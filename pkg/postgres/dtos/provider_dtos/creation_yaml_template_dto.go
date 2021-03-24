package provider_dtos

type ProviderYamlTemplateDto struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Annotations struct {
			CurrentPrimary string `yaml:"current-primary"`
		} `yaml:"annotations"`
		Labels struct {
			CrunchyPghaScope string `yaml:"crunchy-pgha-scope"`
			DeploymentName   string `yaml:"deployment-name"`
			Name             string `yaml:"name"`
			PgCluster        string `yaml:"pg-cluster"`
			PgoVersion       string `yaml:"pgo-version"`
			Pgouser          string `yaml:"pgouser"`
		} `yaml:"labels"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		BackrestStorage struct {
			Accessmode         string `yaml:"accessmode"`
			MatchLabels        string `yaml:"matchLabels"`
			Name               string `yaml:"name"`
			Size               string `yaml:"size"`
			Storageclass       string `yaml:"storageclass"`
			Storagetype        string `yaml:"storagetype"`
			Supplementalgroups string `yaml:"supplementalgroups"`
		} `yaml:"BackrestStorage"`
		PrimaryStorage struct {
			Accessmode         string `yaml:"accessmode"`
			MatchLabels        string `yaml:"matchLabels"`
			Name               string `yaml:"name"`
			Size               string `yaml:"size"`
			Storageclass       string `yaml:"storageclass"`
			Storagetype        string `yaml:"storagetype"`
			Supplementalgroups string `yaml:"supplementalgroups"`
		} `yaml:"PrimaryStorage"`
		ReplicaStorage struct {
			Accessmode         string `yaml:"accessmode"`
			MatchLabels        string `yaml:"matchLabels"`
			Name               string `yaml:"name"`
			Size               string `yaml:"size"`
			Storageclass       string `yaml:"storageclass"`
			Storagetype        string `yaml:"storagetype"`
			Supplementalgroups string `yaml:"supplementalgroups"`
		} `yaml:"ReplicaStorage"`
		Annotations struct {
		} `yaml:"annotations"`
		Ccpimage       string `yaml:"ccpimage"`
		Ccpimageprefix string `yaml:"ccpimageprefix"`
		Ccpimagetag    string `yaml:"ccpimagetag"`
		Clustername    string `yaml:"clustername"`
		Database       string `yaml:"database"`
		Exporterport   string `yaml:"exporterport"`
		Limits         struct {
		} `yaml:"limits"`
		Name         string `yaml:"name"`
		Namespace    string `yaml:"namespace"`
		PgDataSource struct {
			RestoreFrom string `yaml:"restoreFrom"`
			RestoreOpts string `yaml:"restoreOpts"`
		} `yaml:"pgDataSource"`
		Pgbadgerport    string `yaml:"pgbadgerport"`
		Pgoimageprefix  string `yaml:"pgoimageprefix"`
		PodAntiAffinity struct {
			Default    string `yaml:"default"`
			PgBackRest string `yaml:"pgBackRest"`
			PgBouncer  string `yaml:"pgBouncer"`
		} `yaml:"podAntiAffinity"`
		Port        string        `yaml:"port"`
		Tolerations []interface{} `yaml:"tolerations"`
		Tls         struct {
			CaSecret  string `yaml:"caSecret,omitempty"`
			TlsSecret string `yaml:"tlsSecret,omitempty"`
		} `yaml:"tls,omitempty"`
		TlsOnly    bool   `yaml:"tlsOnly,omitempty"`
		User       string `yaml:"user"`
		Userlabels struct {
			PgoVersion string `yaml:"pgo-version"`
		} `yaml:"userlabels"`
	} `yaml:"spec"`
}
