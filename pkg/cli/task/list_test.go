package task

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/joao.rufino/pomo/pkg/cli/test"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestTaskListBuildContainerListOptions(t *testing.T) {

	contexts := []struct {
		listOpts         *listOptions
		expectedSort     bool
		expectedAsJSON   bool
		expectedAll      bool
		expectedLimit    int
		expectedDuration string
	}{
		{
			listOpts: &listOptions{
				all:      true,
				asJSON:   true,
				sort:     true,
				limit:    5,
				duration: "24h",
			},
			expectedSort:     true,
			expectedAsJSON:   true,
			expectedAll:      true,
			expectedLimit:    5,
			expectedDuration: "24h",
		},
		{
			listOpts: &listOptions{
				all:      true,
				asJSON:   false,
				sort:     true,
				limit:    -1,
				duration: "0h",
			},
			expectedSort:     true,
			expectedAsJSON:   false,
			expectedAll:      true,
			expectedLimit:    1,
			expectedDuration: "0h",
		},
		{
			listOpts: &listOptions{
				all:      false,
				asJSON:   true,
				sort:     false,
				limit:    1,
				duration: "rttteta12",
			},
			expectedSort:     false,
			expectedAsJSON:   true,
			expectedAll:      false,
			expectedLimit:    1,
			expectedDuration: "24h",
		},
	}

	for _, c := range contexts {
		options, err := validateTaskListOptions(c.listOpts)
		assert.NilError(t, err)

		assert.Check(t, is.Equal(c.expectedAll, options.all))
		assert.Check(t, is.Equal(c.expectedAsJSON, options.asJSON))
		assert.Check(t, is.Equal(c.expectedSort, options.sort))
		assert.Check(t, is.Equal(c.expectedLimit, options.limit))
		assert.Check(t, is.Equal(c.expectedDuration, options.duration))

	}
}
func TestTaskListErrors(t *testing.T) {
	testCases := []struct {
		args          []string
		flags         map[string]string
		taskListFunc  func(listOptions) error
		expectedError string
	}{
		{
			args: []string{"--in"},
			flags: map[string]string{
				"format": "{{invalid}}",
			},
			expectedError: `unknown flag --in`},
		{
			flags: map[string]string{
				"format": "{{list}}",
			},
			expectedError: `wrong number of args for join`,
		},
		{
			taskListFunc: func(listOptions) error {
				return fmt.Errorf("error listing containers")
			},
			expectedError: "error listing containers",
		},
	}

	buf := new(bytes.Buffer)
	mockCli := test.NewMockCli()
	cmd := NewTaskListCommand(mockCli)
	cmd.SetOut(buf)

	for _, tc := range testCases {
		cmd.SetArgs(tc.args)
		for key, value := range tc.flags {
			cmd.Flags().Set(key, value)
		}
		cmd.Execute()
		assert.Check(t, is.Equal(buf.String(), tc.expectedError))
	}
}

func Test_ExecuteCommand(t *testing.T) {
	cmd := NewTaskListCommand(nil)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--in", "testisawesome"})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "testisawesome" {
		t.Fatalf("expected \"%s\" got \"%s\"", "testisawesome", string(out))
	}
}
