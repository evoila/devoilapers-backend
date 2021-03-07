package core

import (
	"OperatorAutomation/pkg/core/action"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ActionPlaceholder struct {
	Value string
}

func Test_ActionGetter(t *testing.T) {
	var action1 action.IAction = action.Action{
		Name:          "ActionGroup1Item1Name",
		UniqueCommand: "ActionGroup1Item1Cmd",
		Placeholder:   &ActionPlaceholder{},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return "", nil
		},
	}

	var action2 action.IAction = action.Action{
		Name:          "ActionGroup1Item2Name",
		UniqueCommand: "ActionGroup1Item2Cmd",
		Placeholder:   &ActionPlaceholder{},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return "", nil
		},
	}

	var group action.IActionGroup = action.ActionGroup{
		Name: "ActionGroup1",
		Actions: []action.IAction{
			action1,
			action2,
		},
	}

	assert.Equal(t, "ActionGroup1", group.GetName())
	assert.Equal(t, 2, len(group.GetActions()))

	assert.Equal(t, action1.GetName(), group.GetActions()[0].GetName())
	assert.Equal(t, action1.GetUniqueCommand(), group.GetActions()[0].GetUniqueCommand())
	assert.Equal(t, action1.GetPlaceholder(), group.GetActions()[0].GetPlaceholder())
	assert.NotNil(t, action1.GetActionExecuteCallback())

	assert.Equal(t, action2.GetName(), group.GetActions()[1].GetName())
	assert.Equal(t, action2.GetUniqueCommand(), group.GetActions()[1].GetUniqueCommand())
	assert.Equal(t, action2.GetPlaceholder(), group.GetActions()[1].GetPlaceholder())
	assert.NotNil(t, action2.GetActionExecuteCallback())
}

func Test_ActionExecution(t *testing.T) {

	counter := 0
	var action1 action.IAction = action.Action{
		Name:          "Action1Name",
		UniqueCommand: "Action1Cmd",
		Placeholder:   &ActionPlaceholder{Value: "OldValue"},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			counter += 1
			return placeholder.(*ActionPlaceholder).Value, nil
		},
	}

	assert.Equal(t, "Action1Name", action1.GetName())
	assert.Equal(t, "Action1Cmd", action1.GetUniqueCommand())

	placeholder := action1.GetPlaceholder().(*ActionPlaceholder)
	assert.Equal(t, "OldValue", placeholder.Value)
	placeholder.Value = "NewValue"
	value, err := action1.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.Equal(t, "NewValue", value)
	assert.Equal(t, 1, counter)
}
