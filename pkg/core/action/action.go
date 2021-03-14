package action

import (
	"reflect"
)

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
	// Get toggle group
	GetToggleGroup() string
	// Get placeholder
	GetJsonFormResultPlaceholder() interface{}
	// Get action execute callback function for executing an action.
	// Placeholder value equals the placeholder of action.
	// Return value could be any struct
	GetActionExecuteCallback() func(placeholder interface{}) (interface{},error)
	// Get the json form object
	GetJsonForm() interface{}
}

// Action
type Action struct{
	Name string
	UniqueCommand string
	Placeholder interface{}
	ActionExecuteCallback func(placeholder interface{}) (interface{},error)
}

func (a Action) GetToggleGroup() string {
	panic("implement me")
}

func (a Action) GetName() string {
	return a.Name
}

func (a Action) GetUniqueCommand() string {
	return a.UniqueCommand
}

func (a Action) GetJsonFormResultPlaceholder() interface{} {
	return a.Placeholder
}

func (a Action) GetActionExecuteCallback() func(placeholder interface{}) (interface{}, error) {
	return a.ActionExecuteCallback
}

func (a Action) GetJsonForm() interface{} {
	placeholder := a.GetJsonFormResultPlaceholder()

	elemsType := reflect.ValueOf(placeholder)
	if elemsType.Kind() == reflect.Ptr {
		elemsType = elemsType.Elem()
	}

	resultObj := map[string]interface{}{}
	// Loop all properties of the struct
	for i := 0; i< elemsType.NumField(); i++ {
		property := elemsType.Type().Field(i)
		propertyValue :=  elemsType.Field(i)
		propertyMapResult := map[string]interface{}{}


		propertyFormId := property.Tag.Get("json")
		if propertyFormId == "" {
			propertyFormId = property.Name
		}

		// Form type
		propertyFormType := property.Tag.Get("formType")
		if propertyFormType == "" {
			propertyFormType = property.Type.Name()
		}

		propertyMapResult["type"] = propertyFormType

		// Form title
		propertyFormTitle := property.Tag.Get("formTitle")
		if propertyFormTitle == "" {
			propertyFormTitle = property.Name
		}

		propertyMapResult["title"] = propertyFormTitle

		// Form widget
		propertyWidget := property.Tag.Get("formWidget")
		if propertyWidget != "" {
			propertyMapResult["widget"] = propertyWidget
		}

		// Get value via reflection
		var value interface{} = nil
		switch propertyValue.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value= propertyValue.Int()
				break
			case reflect.Bool:
				value = propertyValue.Bool()
				break
			case reflect.String:
				value = propertyValue.String()
				break
		}

		if value != nil {
			propertyMapResult["default"] = value
		}

		// Store sub map into root map
		resultObj[propertyFormId] = propertyMapResult
	}

	return map[string]interface{}{ "properties" : resultObj}
}
