package core

import (
	"OperatorAutomation/pkg/core/action"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ActionPlaceholder struct {
	ValueString string
	ValueInt int
	ValueBool bool
	ValueFile string `formWidget:"file"`

	ValueStringTags string `json:"value_string_tags" formType:"string" formOrder:"1" formTitle:"My ValueStringTags"`
	ValueIntTags int `json:"value_int_tags" formType:"number" formOrder:"2" formTitle:"My ValueIntTags"`
	ValueBoolTags bool `json:"value_bool_tags" formType:"boolean" formOrder:"3" formTitle:"My ValueBoolTags"`
}

func Test_ActionGetter(t *testing.T) {
	var action1 action.IAction = action.FormAction{
		Name:          "ActionGroup1Item1Name",
		UniqueCommand: "ActionGroup1Item1Cmd",
		Placeholder:   &ActionPlaceholder{ValueString: "MyString"},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return "", nil
		},
	}

	var action2 action.IAction = action.FormAction{
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
	assert.Equal(t, action1.GetJsonFormResultPlaceholder(), group.GetActions()[0].GetJsonFormResultPlaceholder())
	assert.NotNil(t, action1.GetActionExecuteCallback())

	form := action1.GetJsonForm()
	jsonBytes, err := json.Marshal(form)
	assert.Nil(t, err)
	jsonText := string(jsonBytes)

	assert.Equal(t, remarshal(`
		{
		  "order": [
			"value_string_tags",
			"value_int_tags",
			"value_bool_tags",
			"ValueString",
			"ValueInt",
			"ValueBool",
			"ValueFile"
		  ],
		  "properties": {
			"ValueBool": {
			  "default": false,
			  "title": "ValueBool",
			  "type": "boolean"
			},
			"ValueFile": {
			  "default": "",
			  "title": "ValueFile",
			  "type": "string",
			  "widget": "file"
			},
			"ValueInt": {
			  "default": 0,
			  "title": "ValueInt",
			  "type": "number"
			},
			"ValueString": {
			  "default": "MyString",
			  "title": "ValueString",
			  "type": "string"
			},
			"value_bool_tags": {
			  "default": false,
			  "title": "My ValueBoolTags",
			  "type": "boolean"
			},
			"value_int_tags": {
			  "default": 0,
			  "title": "My ValueIntTags",
			  "type": "number"
			},
			"value_string_tags": {
			  "default": "",
			  "title": "My ValueStringTags",
			  "type": "string"
			}
		  }
		}
	`), jsonText)

	assert.Equal(t, action2.GetName(), group.GetActions()[1].GetName())
	assert.Equal(t, action2.GetUniqueCommand(), group.GetActions()[1].GetUniqueCommand())
	assert.Equal(t, action2.GetJsonFormResultPlaceholder(), group.GetActions()[1].GetJsonFormResultPlaceholder())
	assert.NotNil(t, action2.GetActionExecuteCallback())
}

// Unmarshal object and Marshal it again to remove white spaces and new lines
func remarshal(input string) string {
	var dynamicObj interface{}
	err := json.Unmarshal([]byte(input), &dynamicObj)

	if err != nil  {
		panic(err)
	}

	jsonBytes, err := json.Marshal(dynamicObj)
	if err != nil  {
		panic(err)
	}

	return string(jsonBytes)
}

func Test_ActionExecution(t *testing.T) {

	counter := 0
	var action1 action.IAction = action.FormAction{
		Name:          "Action1Name",
		UniqueCommand: "Action1Cmd",
		Placeholder:   &ActionPlaceholder{ValueString: "OldValue"},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			counter += 1
			return placeholder.(*ActionPlaceholder).ValueString, nil
		},
	}

	assert.Equal(t, "Action1Name", action1.GetName())
	assert.Equal(t, "Action1Cmd", action1.GetUniqueCommand())

	placeholder := action1.GetJsonFormResultPlaceholder().(*ActionPlaceholder)
	assert.Equal(t, "OldValue", placeholder.ValueString)
	placeholder.ValueString = "NewValue"
	value, err := action1.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.Equal(t, "NewValue", value)
	assert.Equal(t, 1, counter)
}
