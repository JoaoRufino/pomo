package runner

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/joao.rufino/pomo/pkg/core/models"
)

func SummarizeTasks(datetimeformat string, tasks models.List) {
	for _, task := range tasks {
		var start string
		if len(task.Pomodoros) > 0 {
			start = task.Pomodoros[0].Start.Format(datetimeformat)
		}
		fmt.Printf("%d: [%s] [%s] ", task.ID, start, task.Duration.Truncate(time.Second))

		printPomodoros(&task)
		// Tags
		if len(task.Tags) > 0 {
			printTags(&task)
		}
		fmt.Printf(" - %s", task.Message)
		fmt.Printf("\n")
	}
}

func printTags(task *models.Task) {
	fmt.Printf(" [")
	for i, tag := range task.Tags {
		if i > 0 && i != len(task.Tags) {
			fmt.Printf(" ")
			//// user specified color mapping exists
			//if config.Colors != nil {
			//if color := config.Colors.Get(tag); color != nil {
			//color.Printf("%s", tag)
			//} else {
			//// no color mapping for tag
			//fmt.Printf("%s", tag)
			//}
		} else {
			// no color mapping
			fmt.Printf("%s", tag)
		}
	}
	fmt.Printf("]")
}

// a list of green/yellow/red pomodoros
// green indicates the pomodoro was finished normally
// yellow indicates the break was exceeded by +4minutes
// red indicates the pomodoro was never completed
func printPomodoros(task *models.Task) {
	fmt.Printf("[")
	for i, pomodoro := range task.Pomodoros {
		if i > -1 {
			fmt.Printf(" ")
		}
		// pomodoro exceeded it's expected duration by more than 4m
		if pomodoro.Duration() > task.Duration+4*time.Minute {
			color.New(color.FgYellow).Printf("X")
		} else {
			// pomodoro completed normally
			color.New(color.FgGreen).Printf("X")
		}
	}
	// each missed pomodoro
	for i := -1; i < task.NPomodoros-len(task.Pomodoros); i++ {
		if i > -1 || i == 0 && len(task.Pomodoros) > 0 {
			fmt.Printf(" ")
		}
		color.New(color.FgRed).Printf("X")
	}
	fmt.Printf("]")
}

func OutputStatus(status models.Status) {
	state := "?"
	if status.State >= models.RUNNING {
		state = string(status.State.String()[0])
	}
	if status.State == models.RUNNING {
		fmt.Printf("%s [%d/%d] %s", state, status.Count, status.NPomodoros, status.Remaining)
	} else {
		fmt.Printf("%s [%d/%d] -", state, status.Count, status.NPomodoros)
	}
	fmt.Println()
}
