package task

import (
	"io/ioutil"
	"sort"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRemoveForce(t *testing.T) {
	var removed []string

	cmd := NewTaskDeleteCommand(nil)
	cmd.SetOut(ioutil.Discard)

	t.Run("without force", func(t *testing.T) {
		cmd.SetArgs([]string{"nosuchcontainer", "mycontainer"})
		removed = []string{}
		assert.ErrorContains(t, cmd.Execute(), "No such container")
		sort.Strings(removed)
		assert.DeepEqual(t, removed, []string{"mycontainer", "nosuchcontainer"})
	})
	t.Run("with force", func(t *testing.T) {
		cmd.SetArgs([]string{"--force", "nosuchcontainer", "mycontainer"})
		removed = []string{}
		assert.NilError(t, cmd.Execute())
		sort.Strings(removed)
		assert.DeepEqual(t, removed, []string{"mycontainer", "nosuchcontainer"})
	})
}
