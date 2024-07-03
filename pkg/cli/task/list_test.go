package task

import (
	"bytes"
	"testing"

	"github.com/joaorufino/pomo/pkg/cli/test"
	testClient "github.com/joaorufino/pomo/pkg/client/test"
	"github.com/joaorufino/pomo/pkg/core/models"
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
		doc           string
		args          []string
		flags         map[string]string
		taskListFunc  func(listOptions) error
		expectedError string
	}{
		{
			doc:           "unknown flag",
			args:          []string{"--fake"},
			flags:         map[string]string{},
			expectedError: `unknown flag: --fake`,
		},
	}
	buf := new(bytes.Buffer)
	mockCli := test.NewMockCli()
	cmd := NewTaskListCommand(mockCli)
	cmd.SetOut(buf)

	for _, tc := range testCases {
		t.Run(tc.doc, func(t *testing.T) {
			cmd.SetArgs(tc.args)
			for key, value := range tc.flags {
				cmd.Flags().Set(key, value)
			}
			assert.ErrorContains(t, cmd.Execute(), tc.expectedError)
		})
	}
}

func TestRunTaskListWithInvalidArguments(t *testing.T) {
	var testcases = []struct {
		doc         string
		options     listOptions
		expectedErr string
	}{
		{
			doc: "invalid duration",
			options: listOptions{
				duration: "r",
			},
			expectedErr: "time: invalid duration \"r\"",
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.doc, func(t *testing.T) {
			err := list(test.NewMockCli(), &testcase.options)
			assert.Error(t, err, testcase.expectedErr)
		})
	}
}

func TestRunTaskList(t *testing.T) {
	var testcases = []struct {
		doc            string
		clientOptions  testClient.MockClientOptions
		options        listOptions
		expectedResult string
	}{
		{
			doc: "all values",
			options: listOptions{
				duration: "24h",
			},
			clientOptions: testClient.MockClientOptions{
				List: &models.List{
					{
						ID:         1,
						Message:    "ola",
						Pomodoros:  nil,
						Tags:       []string{},
						NPomodoros: 0,
						Duration:   1,
					},
				},
			},
			expectedResult: "",
		},
	}

	for _, testcase := range testcases {
		mockCli := test.NewMockCli()
		client := testClient.NewMockClient(mockCli.Config(), testcase.clientOptions)
		mockCli.SetClient(&client)
		t.Run(testcase.doc, func(t *testing.T) {
			err := list(mockCli, &testcase.options)
			if err != nil {
				assert.Error(t, err, testcase.expectedResult)
			}
		})
	}
}
