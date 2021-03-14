package unit_tests

import (
	action2 "OperatorAutomation/pkg/core/action"
	"encoding/json"
	"fmt"
	"testing"
)

type MyStruct struct {
	Value int `json:"value" formType:"number" formTitle:"My value: "`
	Name string `formType:"string" formTitle:"My Name: "`
	File string `formType:"string" formWidget:"file" formTitle:"My File: "`
}

func Test_Json(t *testing.T) {
	action := action2.Action{
		Name: "MyAct",
		Placeholder: &MyStruct{Value: 10},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {

			fmt.Println(placeholder)
			return nil, nil
		},
	}

	res := action.GetJsonForm()
	jBytes, _ := json.Marshal(res)
	jString := string(jBytes)
	print(jString)

}