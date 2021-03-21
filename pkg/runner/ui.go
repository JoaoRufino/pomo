package runner

import (
	"fmt"
	"time"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/joao.rufino/pomo/pkg/server"
)

func render(wheel *server.Wheel, status *server.Status) *widgets.Paragraph {
	var text string
	switch status.State {
	case server.RUNNING:
		text = fmt.Sprintf(
			`[%d/%d] Pomodoros completed

			%s %s remaining


			[q] - quit [p] - pause
			`,
			status.Count,
			status.NPomodoros,
			wheel,
			status.Remaining,
		)
	case server.BREAKING:
		text = `It is time to take a break!

		Once you are ready, press [enter] 
		to begin the next Pomodoro.

		[q] - quit [p] - pause
		`
	case server.PAUSED:
		text = `Pomo is suspended.

		Press [p] to continue.


		[q] - quit [p] - unpause
		`
	case server.COMPLETE:
		text = `This session has concluded. 
		
		Press [q] to exit.


		


		[q] - quit
		`
	}
	par := widgets.NewParagraph()
	par.Text = text
	par.SetRect(0, 10, 20, 20)
	par.Title = fmt.Sprintf("Pomo - %s", status.State)
	par.BorderStyle.Fg = termui.ColorRed
	if status.State == server.RUNNING {
		par.BorderStyle.Fg = termui.ColorGreen
	}
	return par
}

func newBlk() *termui.Block {
	blk := termui.NewBlock()
	blk.Border = false
	return blk
}

func centered(part *widgets.Paragraph) *termui.Grid {
	grid := termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		termui.NewRow(1.0/2,
			termui.NewCol(1.0/2, newBlk()),
		),
		termui.NewRow(1.0/4,
			termui.NewCol(1.0/3, newBlk()),
			termui.NewCol(1.0/3, part),
			termui.NewCol(1.0/3, newBlk()),
		),
		termui.NewRow(1.0/4,
			termui.NewCol(1.0/4, newBlk()),
		),
	)
	return grid
}

func StartUI(runner *TaskRunner) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	wheel := server.Wheel(0)

	defer termui.Close()

	termui.Render(centered(render(&wheel, runner.Status())))
	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<Enter>":
				runner.Toggle()
				termui.Render(centered(render(&wheel, runner.Status())))
			case "q", "<C-c>":
				return
			case "p":
				runner.Pause()
				termui.Render(centered(render(&wheel, runner.Status())))
			case "<Resize>":
				termui.Render(centered(render(&wheel, runner.Status())))
			}
		case <-ticker:
			termui.Render(centered(render(&wheel, runner.Status())))
		}
	}
}
