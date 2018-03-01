package core
  
import (
        "testing"
        "github.com/stretchr/testify/assert"
)

func Test_getUsage(t *testing.T) {
        task_command := TaskCommand{}
        assert.Equal(t, task_command.GetUsage(), "listtasks <username>", "they should be equal")
}

func Test_getDescriptionNil(t *testing.T) {
        task_command := TaskCommand{}
        assert.NotEmpty(t, task_command.GetDescription())
}

func Test_getMatchersNil(t *testing.T) {
        task_command := TaskCommand{}
        assert.NotEmpty(t, task_command.GetMatchers())
        assert.Contains(t, task_command.GetMatchers(), "^listtasks\\s*(.*)$")
}

func Test_getIMMatchersNil(t *testing.T) {
        task_command := TaskCommand{}
        assert.NotEmpty(t, task_command.GetIMMatchers())
}

func Test_getMentionMatchersNil(t *testing.T) {
        task_command := TaskCommand{}
        assert.NotEmpty(t, task_command.GetMentionMatchers())
}

func Test_getHandlerNil(t *testing.T) {
        task_command := TaskCommand{}
        assert.NotEmpty(t, task_command.GetHandler())
}
