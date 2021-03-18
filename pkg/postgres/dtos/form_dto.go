package dtos

// Autogenerated from templates/postgres/create_form.json
type FormResponseDto struct {
	Common struct {
		ClusterName        string `json:"clusterName"`
		Username           string `json:"username"`
		InClusterPort      int    `json:"inClusterPort"`
		ClusterCtorageSize int    `json:"clusterCtorageSize"`
	} `json:"common"`
	TLS struct {
		UseTLS            bool   `json:"useTls"`
		TLSMode           string `json:"tlsMode"`
		TLSModeFromSecret struct {
			CaSecret  string `json:"caSecret"`
			TLSSecret string `json:"tlsSecret"`
		} `json:"tlsModeFromSecret"`
	} `json:"tls"`
	Backup struct {
		PerformBackup bool `json:"performBackup"`
		CommonS3Data  struct {
			S3BucketName string `json:"s3bucketName"`
			S3Endpoint   string `json:"s3endpoint"`
			S3Region     string `json:"s3region"`
		} `json:"commonS3data"`
		BackupMode                   string `json:"backupMode"`
		BackupModeFromNewCredentials struct {
			S3Key    string `json:"s3key"`
			S3Secret string `json:"s3secret"`
		} `json:"backupModeFromNewCredentials"`
		BackupModeFromSecret struct {
			S3Key    string `json:"s3key"`
			S3Secret string `json:"s3secret"`
		} `json:"backupModeFromSecret"`
	} `json:"backup"`
}

// Autogenerated from templates/postgres/create_form.json
type FormQueryDto struct {
	Properties struct {
		Common struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Properties  struct {
				ClusterName struct {
					Type    string `json:"type"`
					Title   string `json:"title"`
					Default string `json:"default"`
					Widget  struct {
						ID string `json:"id"`
					} `json:"widget"`
				} `json:"clusterName"`
				Username struct {
					Type   string `json:"type"`
					Title  string `json:"title"`
					Widget struct {
						ID string `json:"id"`
					} `json:"widget"`
					Default string `json:"default"`
				} `json:"username"`
				InClusterPort struct {
					Type   string `json:"type"`
					Title  string `json:"title"`
					Widget struct {
						ID string `json:"id"`
					} `json:"widget"`
					Default int `json:"default"`
				} `json:"inClusterPort"`
				ClusterCtorageSize struct {
					Type   string `json:"type"`
					Title  string `json:"title"`
					Widget struct {
						ID string `json:"id"`
					} `json:"widget"`
					Default int `json:"default"`
				} `json:"clusterCtorageSize"`
			} `json:"properties"`
			Order []string `json:"order"`
		} `json:"common"`
		TLS struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Properties  struct {
				UseTLS struct {
					Type   string `json:"type"`
					Widget struct {
						ID string `json:"id"`
					} `json:"widget"`
					Description string `json:"description"`
					Default     bool   `json:"default"`
				} `json:"useTls"`
				TLSMode struct {
					Type    string `json:"type"`
					Default string `json:"default"`
					Widget  struct {
						ID string `json:"id"`
					} `json:"widget"`
					OneOf []struct {
						Enum        []string `json:"enum"`
						Description string   `json:"description"`
					} `json:"oneOf"`
					VisibleIf struct {
						UseTLS []bool `json:"useTls"`
					} `json:"visibleIf"`
				} `json:"tlsMode"`
				TLSModeFromSecret struct {
					Type       string `json:"type"`
					Properties struct {
						CaSecret struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"caSecret"`
						TLSSecret struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"tlsSecret"`
					} `json:"properties"`
					VisibleIf struct {
						TLSMode []string `json:"tlsMode"`
					} `json:"visibleIf"`
				} `json:"tlsModeFromSecret"`
				TLSModeFromFile struct {
					Type       string `json:"type"`
					Name       string `json:"name"`
					Properties struct {
						CaCert struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"caCert"`
						TLSPrivateKey struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"tlsPrivateKey"`
						TLSCertificate struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"tlsCertificate"`
					} `json:"properties"`
					VisibleIf struct {
						TLSMode []string `json:"tlsMode"`
					} `json:"visibleIf"`
				} `json:"tlsModeFromFile"`
			} `json:"properties"`
			Order []string `json:"order"`
		} `json:"tls"`
		Backup struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Properties  struct {
				PerformBackup struct {
					Type   string `json:"type"`
					Widget struct {
						ID string `json:"id"`
					} `json:"widget"`
					Description string `json:"description"`
					Default     bool   `json:"default"`
				} `json:"performBackup"`
				CommonS3Data struct {
					Type        string `json:"type"`
					Description string `json:"description"`
					Properties  struct {
						S3BucketName struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"s3bucketName"`
						S3Endpoint struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"s3endpoint"`
						S3Region struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"s3region"`
					} `json:"properties"`
					VisibleIf struct {
						PerformBackup []bool `json:"performBackup"`
					} `json:"visibleIf"`
				} `json:"commonS3data"`
				BackupMode struct {
					Type    string `json:"type"`
					Default string `json:"default"`
					Widget  struct {
						ID string `json:"id"`
					} `json:"widget"`
					OneOf []struct {
						Enum        []string `json:"enum"`
						Description string   `json:"description"`
					} `json:"oneOf"`
					VisibleIf struct {
						PerformBackup []bool `json:"performBackup"`
					} `json:"visibleIf"`
				} `json:"backupMode"`
				BackupModeFromSecret struct {
					Type       string `json:"type"`
					Properties struct {
						S3Key struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"s3key"`
						S3Secret struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"s3secret"`
					} `json:"properties"`
					VisibleIf struct {
						BackupMode []string `json:"backupMode"`
					} `json:"visibleIf"`
				} `json:"backupModeFromSecret"`
				BackupModeFromNewCredentials struct {
					Type       string `json:"type"`
					Properties struct {
						S3Key struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"s3key"`
						S3Secret struct {
							Type   string `json:"type"`
							Title  string `json:"title"`
							Widget struct {
								ID string `json:"id"`
							} `json:"widget"`
						} `json:"s3secret"`
					} `json:"properties"`
					VisibleIf struct {
						BackupMode []string `json:"backupMode"`
					} `json:"visibleIf"`
				} `json:"backupModeFromNewCredentials"`
			} `json:"properties"`
			Order []string `json:"order"`
		} `json:"backup"`
	} `json:"properties"`
}