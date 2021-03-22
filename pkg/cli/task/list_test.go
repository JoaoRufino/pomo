package task

import (
	"testing"

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
