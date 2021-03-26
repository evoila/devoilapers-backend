package action

import (
	"errors"
)

type ToggleActionPlaceholder struct {
	Toggle string `json:"toggle"`
}

func CreateToggleAction(
	name string,
	uniqueCommand string,
	queryExecuteCallback func() (bool, error),
	setExecuteCallback func() (interface{}, error),
	unsetExecuteCallback func() (interface{}, error)) ToggleAction {

	return ToggleAction{
		Name:                 name,
		UniqueCommand:        uniqueCommand,
		QueryExecuteCallback: queryExecuteCallback,
		SetExecuteCallback:   setExecuteCallback,
		UnsetExecuteCallback: unsetExecuteCallback,
		placeholder:          &ToggleActionPlaceholder{},
	}
}

// Toggle action. Use factory method CreateToggleAction to construct.
type ToggleAction struct {
	Name                 string
	UniqueCommand        string
	QueryExecuteCallback func() (bool, error)
	SetExecuteCallback   func() (interface{}, error)
	UnsetExecuteCallback func() (interface{}, error)

	placeholder *ToggleActionPlaceholder
}

func (a ToggleAction) GetIsToggleAction() bool {
	return true
}

func (a ToggleAction) GetName() string {
	return a.Name
}

func (a ToggleAction) GetUniqueCommand() string {
	return a.UniqueCommand
}

func (a ToggleAction) GetJsonFormResultPlaceholder() interface{} {
	return a.placeholder
}

func (a ToggleAction) GetActionExecuteCallback() func(placeholder interface{}) (interface{}, error) {
	return a.executionCallback
}

func (a ToggleAction) GetJsonForm() interface{} {
	// Create ngx json form with 3 options to get,set,unset this toggle action
	return map[string]interface{}{
		"properties": map[string]interface{}{
			"toggle": map[string]interface{}{
				"type":   "string",
				"title":  "toggle",
				"widget": "select",
				"oneOf": []map[string]interface{}{
					{
						"enum":        []string{"get"},
						"description": "get",
					},
					{
						"enum":        []string{"set"},
						"description": "set",
					},
					{
						"enum":        []string{"unset"},
						"description": "unset",
					},
				},
			},
		},
	}
}

// Gets filled placeholder that is defined above
func (a ToggleAction) executionCallback(placeholder interface{}) (interface{}, error) {
	toggleDto := placeholder.(*ToggleActionPlaceholder)

	switch toggleDto.Toggle {
	case "get":
		return a.QueryExecuteCallback()
	case "set":
		return a.SetExecuteCallback()
	case "unset":
		return a.UnsetExecuteCallback()
	}

	return nil, errors.New("invalid toggle state")
}
