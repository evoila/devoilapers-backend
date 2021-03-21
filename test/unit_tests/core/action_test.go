package core

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/test/unit_tests/common_test"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ActionPlaceholder struct {
	ValueString string
	ValueInt    int
	ValueBool   bool
	ValueFile   string `formWidget:"file"`

	ValueStringTags string `json:"value_string_tags" formType:"string" formOrder:"1" formTitle:"My ValueStringTags"`
	ValueIntTags    int    `json:"value_int_tags" formType:"number" formOrder:"2" formTitle:"My ValueIntTags"`
	ValueBoolTags   bool   `json:"value_bool_tags" formType:"boolean" formOrder:"3" formTitle:"My ValueBoolTags"`
}

func Test_FormActionGetter(t *testing.T) {
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
	assert.False(t, action1.GetIsToggleAction())

	form := action1.GetJsonForm()
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `
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
	`, form))

	assert.Equal(t, action2.GetName(), group.GetActions()[1].GetName())
	assert.Equal(t, action2.GetUniqueCommand(), group.GetActions()[1].GetUniqueCommand())
	assert.Equal(t, action2.GetJsonFormResultPlaceholder(), group.GetActions()[1].GetJsonFormResultPlaceholder())
	assert.NotNil(t, action2.GetActionExecuteCallback())
	assert.False(t, action2.GetIsToggleAction())

	// Check nil placeholders
	var actionNil action.IAction = action.FormAction{
		Name:          "NilAction",
		UniqueCommand: "CmdNilAction",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return "", nil
		},
	}

	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `{"properties": {}}`, actionNil.GetJsonForm()))

}

func Test_FormAction_Invalid_NonPointerPlaceholder(t *testing.T) {
	// Gets called on destruction. Ensures there was a panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Invalid form action does not panic")
		}
	}()

	var actionUnsupportedType action.IAction = action.FormAction{
		Name:          "NilAction",
		UniqueCommand: "CmdNilAction",
		Placeholder:   ActionPlaceholder{ValueString: "Some value"}, // Not a pointer
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return "", nil
		},
	}

	// Should panic
	actionUnsupportedType.GetJsonForm()
}

func Test_FormAction_Invalid_PlaceholderField(t *testing.T) {
	unsupportedTypeStruct := struct {
		InvalidField []string // Not supported field
	}{}

	// Gets called on destruction. Ensures there was a panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Invalid formaction does not panic")
		}
	}()

	var actionUnsupportedType action.IAction = action.FormAction{
		Name:          "NilAction",
		UniqueCommand: "CmdNilAction",
		Placeholder:   &unsupportedTypeStruct,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return "", nil
		},
	}

	// Should panic
	actionUnsupportedType.GetJsonForm()
}

func Test_FormAction_Invalid_FormOrderField(t *testing.T) {
	unsupportedFormOrderStruct := struct {
		InvalidField string `formOrder:"NotANumber"` // Not supported field
	}{}

	// Gets called on destruction. Ensures there was a panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Invalid formaction does not panic")
		}
	}()

	var actionUnsupportedType action.IAction = action.FormAction{
		Name:          "NilAction",
		UniqueCommand: "CmdNilAction",
		Placeholder:   &unsupportedFormOrderStruct,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return "", nil
		},
	}

	// Should panic
	actionUnsupportedType.GetJsonForm()
}

func Test_ToggleActionGetter(t *testing.T) {
	state := false
	var action1 action.IAction = action.CreateToggleAction(
		"ToggleAction",
		"cmd_toggle_action",
		func() (bool, error) {
			// Get
			return state, nil
		},
		func() (interface{}, error) {
			// Set
			state = true
			return map[string]interface{}{"now": state}, nil
		},
		func() (interface{}, error) {
			// Unset
			state = false
			return map[string]interface{}{"now": state}, nil
		},
	)

	assert.Equal(t, "ToggleAction", action1.GetName())
	assert.Equal(t, "cmd_toggle_action", action1.GetUniqueCommand())
	assert.True(t, action1.GetIsToggleAction())
	assert.NotNil(t, action1.GetJsonFormResultPlaceholder())
	assert.NotNil(t, action1.GetActionExecuteCallback())

	// Check ngx schema form
	form := action1.GetJsonForm()
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `
			{
			  "properties": {
				"toggle": {
				  "oneOf": [
					{
					  "description": "get",
					  "enum": [
						"get"
					  ]
					},
					{
					  "description": "set",
					  "enum": [
						"set"
					  ]
					},
					{
					  "description": "unset",
					  "enum": [
						"unset"
					  ]
					}
				  ],
				  "title": "toggle",
				  "type": "string",
				  "widget": "select"
				}
			  }
			}
	`, form))

	// Get
	placeholder := action1.GetJsonFormResultPlaceholder()
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"get"}`), placeholder))
	response, err := action1.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `false`, response))

	// Set
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"set"}`), placeholder))
	response, err = action1.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `{"now":true}`, response))

	// Get again
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"get"}`), placeholder))
	response, err = action1.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `true`, response))

	// Unset
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"unset"}`), placeholder))
	response, err = action1.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `{"now":false}`, response))

	// Get again
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"get"}`), placeholder))
	response, err = action1.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `false`, response))

	// Invalid state
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"somestate"}`), placeholder))
	response, err = action1.GetActionExecuteCallback()(placeholder)
	assert.NotNil(t, err)
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
