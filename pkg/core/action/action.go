package action

type IActionGroup interface {
	// Get action group name
	GetName() string
	// Get actions of this action group
	GetActions() []IAction
}

type ActionGroup struct {
	// Name of the action group
	Name string
	// Actions of this group
	Actions []IAction
}

func (ag ActionGroup) GetName() string {
	return ag.Name
}

func (ag ActionGroup) GetActions() []IAction {
	return ag.Actions
}

type IAction interface {
	// Get action name
	GetName() string
	// Get the unique command
	GetUniqueCommand() string
	// Get placeholder
	GetJsonFormResultPlaceholder() interface{}
	// Get action execute callback function for executing an action.
	// Placeholder value equals the placeholder of action.
	// Return value could be any struct
	GetActionExecuteCallback() func(placeholder interface{}) (interface{}, error)
	// Get the json form object
	GetJsonForm() interface{}
	// Determine if this action is a toggle action
	GetIsToggleAction() bool
}
