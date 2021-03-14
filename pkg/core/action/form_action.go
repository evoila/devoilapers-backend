package action

import (
	"reflect"
	"sort"
	"strconv"
)

// Action
type FormAction struct {
	Name                  string
	UniqueCommand         string
	Placeholder           interface{}
	ActionExecuteCallback func(placeholder interface{}) (interface{}, error)
}

func (a FormAction) GetIsToggleAction() bool {
	return false
}

func (a FormAction) GetName() string {
	return a.Name
}

func (a FormAction) GetUniqueCommand() string {
	return a.UniqueCommand
}

func (a FormAction) GetJsonFormResultPlaceholder() interface{} {
	return a.Placeholder
}

func (a FormAction) GetActionExecuteCallback() func(placeholder interface{}) (interface{}, error) {
	return a.ActionExecuteCallback
}

// Converts the placeholder into a json form object
// Valid tags:
// json: (string) Property name which will also be used by the response value
// formOrder: (int) Order index. Ensure all properties are using this tag if one uses it
// formType: (string) Type of the form. Valid: string, boolean, number. For files refer to string and use widget
// formTitle: (string) Title displayed in the form
// formWidget: (string) Widget that is used. Valid: file
// Note:
// default value is the assigned property value
func (a FormAction) GetJsonForm() interface{} {
	// Build form accordingly to
	// https://github.com/guillotinaweb/ngx-schema-form

	// Get the placeholder struct pointer
	placeholder := a.GetJsonFormResultPlaceholder()

	// Get type of placeholder
	placeholderType := reflect.ValueOf(placeholder)

	// Ensure we only allow pointers
	if placeholderType.Kind() != reflect.Ptr {
		panic("Placeholder is not of type pointer")
	}

	// Unwrap underlying struct of the pointer
	placeholderType = placeholderType.Elem()

	resultObj := map[string]interface{}{}

	orderIndicies := make([]int, 0)
	propertyOrder := map[int]string{}

	// Loop all properties of the struct
	for i := 0; i < placeholderType.NumField(); i++ {
		property := placeholderType.Type().Field(i)
		propertyValue := placeholderType.Field(i)
		propertyMapResult := map[string]interface{}{}

		// Get id
		propertyFormId := property.Tag.Get("json")
		if propertyFormId == "" {
			propertyFormId = property.Name
		}

		// Remember the property order
		propertyFormOrder := 1000 + i // 1000 Offset for properties without order tags
		propertyFormOrderString := property.Tag.Get("formOrder")
		if propertyFormOrderString != "" {
			numericOrder, err := strconv.Atoi(propertyFormOrderString)
			if err != nil {
				panic("Invalid order symbol")
			}

			propertyFormOrder = numericOrder
		}

		propertyOrder[propertyFormOrder] = propertyFormId
		orderIndicies = append(orderIndicies, propertyFormOrder)

		// Form type
		propertyFormType := property.Tag.Get("formType")
		if propertyFormType == "" {
			propertyFormType = property.Type.Name()

			switch property.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				propertyFormType = "number"
				break
			case reflect.Bool:
				propertyFormType = "boolean"
				break
			case reflect.String:
				propertyFormType = "string"
				break
			default:
				panic("Unsupported auto conversion of type " + property.Type.Name() + " of property " + property.Name)
			}
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
			value = propertyValue.Int()
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

	// Append the sorted the properties to a list
	sort.Ints(orderIndicies)
	order := make([]string, 0)
	for _, indicie := range orderIndicies {
		order = append(order, propertyOrder[indicie])
	}

	// Wrap into root object
	return map[string]interface{}{"properties": resultObj, "order": order}
}
