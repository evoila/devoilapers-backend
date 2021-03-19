package dtos

type CreateRepoDto struct {
	Master_timeout string      `json:"master_timeout"`
	Timeout        string      `json:"timeout"`
	Type           string      `json:"type"`
	Settings       RepoSetting `json:"settings"`
	Verify         bool        `json:"verify"`
}

type RepoSetting struct {
	Chunk_size                 string `json:"chunk_size"`
	Compress                   bool   `json:"compress"`
	Max_number_of_snapshots    int32  `json:"max_number_of_snapshots"`
	Max_restore_bytes_per_sec  string `json:"max_restore_bytes_per_sec"`
	Max_snapshot_bytes_per_sec string `json:"max_snapshot_bytes_per_sec"`
	Readonly                   bool   `json:"readonly"`
	Location                   string `json:"location"`
	Delegate_type              string `json:"delegate_type"`
	Url                        string `json:"url"`
}

type RepoDto struct {
	Repository string
	Body       CreateRepoDto
}

type CreateSnapshotDto struct {
	Ignore_Unavailable   bool     `json:"ignore_unavailable"`
	Indices              string   `json:"indices"`
	Include_Global_State bool     `json:"include_global_state"`
	Master_Timeout       string   `json:"master_timeout"`
	Metadata             Metadata `json:"metadata"`
	Partial              bool     `json:"partial"`
	Wait_For_Completion  bool     `json:"wait_for_completion"`
}

type Metadata struct {
	Taken_By      string `json:"taken_by"`
	Taken_Because string `json:"taken_because"`
}

type SnapshotDto struct {
	Repository string
	Snapshot   string
	Body       CreateSnapshotDto
}
