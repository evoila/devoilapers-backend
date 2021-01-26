package kibana

import form "goclient/RequestForms/Kibana"

func getScalePlaceHolder() form.ScaleRequest {
	return form.ScaleRequest{
		Replicas: 1,
	}
}
