package postgres

import (
	form "goclient/RequestForms/Postgresql"
	structs "goclient/RestStruct/Postgres/structs"
)

func getScalePlaceHolder() form.ScaleRequest {
	return form.ScaleRequest{
		Name:                "cluster name",
		ClusterScaleRequest: structs.ClusterScaleRequest{},
	}
}

func getBackupPlaceHolder() form.CreateBackupRequest {
	return form.CreateBackupRequest{}
}
