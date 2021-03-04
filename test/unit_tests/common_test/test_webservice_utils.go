package common_test

import (
	"encoding/json"
	"net/http/httptest"
)

func GetMessageAndCode(responseRecorder *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(responseRecorder.Body.Bytes(), target)

	if err != nil {
		panic(err)
	}

	receivedJson := string(responseRecorder.Body.Bytes())
	serializedTargetBytes, err := json.Marshal(target)
	if err != nil {
		panic(err)
	}

	same := receivedJson == string(serializedTargetBytes)
	if !same {
		panic(err)
	}
}
