package elasticsearch

import (
	form "goclient/RequestForms/Elasticsearch"
	"goclient/RestStruct/Elasticsearch/structs"
)

func getBackupPlaceHolder() form.CreateSnapshotRequest {
	return form.CreateSnapshotRequest{
		My_Repo:        "my repository",
		My_Snapshot:    "my backup",
		CreateSnapshot: structs.CreateSnapshot{},
	}
}

func getScalePlaceHolder() form.ScaleRequest {
	return form.ScaleRequest{
		Replicas: 1,
	}
}
