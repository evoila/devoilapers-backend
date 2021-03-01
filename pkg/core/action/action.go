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
	GetPlaceholder() interface{}
	// Getaction execute callback function for executing an action.
	// Placeholder value equals the placeholder of action.
	// Return value could be any struct
	GetActionExecuteCallback() func(placeholder interface{}) (interface{},error)
}
// Action
type Action struct{
	Name string
	UniqueCommand string
	Placeholder interface{}
	ActionExecuteCallback func(placeholder interface{}) (interface{},error)
}

func (a Action) GetName() string {
	return a.Name
}

func (a Action) GetUniqueCommand() string {
	return a.UniqueCommand
}

func (a Action) GetPlaceholder() interface{} {
	return a.Placeholder
}

func (a Action) GetActionExecuteCallback() func(placeholder interface{}) (interface{}, error) {
	return a.ActionExecuteCallback
}