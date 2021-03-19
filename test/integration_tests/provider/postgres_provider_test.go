package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres"
	dtos "OperatorAutomation/pkg/postgres/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func CreatePostgresTestProvider(t *testing.T) (*provider.IServiceProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	var pgProvider provider.IServiceProvider = postgres.CreatePostgresProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.ResourcesTemplatesPath,
		kubernetes.NginxInformation(config.Kubernetes.Nginx),
	)

	return &pgProvider, config
}

func Test_Postgres_Provider_Create_Panic_Template_Not_Found(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Postgres provider did not panic if the template could not be found")
		}
	}()

	pgProvider := postgres.CreatePostgresProvider(
		"Server",
		"CaPath",
		"NotExistingPath",
		kubernetes.NginxInformation{})

	_ = pgProvider
}

func Test_Postgres_Provider_GetAttributes(t *testing.T) {
	pgProviderPtr, _ := CreatePostgresTestProvider(t)
	pgProvider := *pgProviderPtr

	assert.NotEqual(t, "", pgProvider.GetServiceImage())
	assert.NotEqual(t, "", pgProvider.GetServiceDescription())
	assert.Equal(t, "Postgres", pgProvider.GetServiceType())

	testUser := unit_test.TestUser{
		KubernetesNamespace: "MyNamespace",
	}

	// Get json form data
	formDataObj1, err := pgProvider.GetJsonForm(testUser)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj1)
	formDataObj2, err := pgProvider.GetJsonForm(testUser)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj2)

	// Ensure they are not the same (because of the random name)
	formData1 := formDataObj1.(dtos.FormQueryDto)
	formData2 := formDataObj2.(dtos.FormQueryDto)

	assert.NotEqual(t,
		formData1.Properties.Common.Properties.ClusterName.Default,
		formData2.Properties.Common.Properties.ClusterName.Default)

	// Generate yaml from form values and ensure it sets the values from the form
	filledForm := dtos.FormResponseDto{}
	filledForm.Common.ClusterName = "MyCluster"

	filledFormData, err := json.Marshal(filledForm)
	assert.Nil(t, err)
	yamlTemplate, err := pgProvider.GetYamlTemplate(testUser, filledFormData)
	assert.Nil(t, err)
	assert.NotNil(t, yamlTemplate)

	yamlObject := yamlTemplate.(dtos.ProviderYamlTemplateDto)

	// Ensure values are set by form as expected
	expectedClusterName := "MyCluster"
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Annotations.CurrentPrimary)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Labels.CrunchyPghaScope)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Labels.DeploymentName)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Labels.Name)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Labels.PgCluster)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Name)
	assert.Equal(t, expectedClusterName, yamlObject.Spec.PrimaryStorage.Name)
	assert.Equal(t, expectedClusterName, yamlObject.Spec.Clustername)
	assert.Equal(t, expectedClusterName, yamlObject.Spec.Database)
	assert.Equal(t, expectedClusterName, yamlObject.Spec.Name)

	expectedNamespace := "MyNamespace"
	assert.Equal(t, expectedNamespace, yamlObject.Metadata.Namespace)
	assert.Equal(t, expectedNamespace, yamlObject.Spec.Namespace)
}

func get_action(service *service.IService, groupname string, actioncommand string) (*action.IAction, error) {
	actionGroups := (*service).GetActionGroups()

	for _, actionGroup := range actionGroups {
		if actionGroup.GetName() != groupname {
			continue
		}

		actions := actionGroup.GetActions()
		for actionIdx, action := range actions {
			if action.GetUniqueCommand() == actioncommand {
				return &actions[actionIdx], nil
			}
		}
	}

	return nil, errors.New("Action not found")
}

func Test_Postgres_Provider_End2End(t *testing.T) {
	pgProviderPtr, config := CreatePostgresTestProvider(t)
	pgProvider := *pgProviderPtr

	user := config.Users[0]

	// Prepare form
	filledForm := dtos.FormResponseDto{}
	filledForm.Common.ClusterName = "pg-test-cluster"

	service1Ptr := common_test.CommonProviderStart(t, pgProviderPtr, user, filledForm, 3)
	service1 := *service1Ptr

	// Testing actions
	// Get database credentials
	actionPtr, err := get_action(service1Ptr, "Informations", "cmd_pg_get_credentials")
	assert.Nil(t, err)
	action := *actionPtr
	assert.Nil(t, action.GetJsonFormResultPlaceholder())
	result, err := action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterCredentials := *result.(*dtos.ClusterCredentialsDto)
	assert.True(t, len(clusterCredentials.Username) > 3)
	assert.True(t, len(clusterCredentials.Password) > 3)
	assert.True(t, clusterCredentials.InternalPort > 1)

	// Exposure
	actionPtr, err = get_action(service1Ptr, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.NotNil(t, err) // Not exposed

	// Expose it
	actionPtr, err = get_action(service1Ptr, "Security", "cmd_pg_expose")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeResult := result.(*dtos.ClusterExposeResponseDto)
	assert.True(t, clusterExposeResult.Port > 1)

	// Check again if it is exposed
	serviceTemp, err := pgProvider.GetService(user, service1.GetName())
	assert.Nil(t, err)
	actionPtr, err = get_action(serviceTemp, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation := result.(*dtos.ClusterExposeResponseDto)
	assert.Equal(t, clusterExposeResult.Port, clusterExposeInformation.Port)

	// Hide it again
	actionPtr, err = get_action(service1Ptr, "Security", "cmd_pg_hide")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)

	// Check again if it is hidden
	serviceTemp, err = pgProvider.GetService(user, service1.GetName())
	assert.Nil(t, err)
	actionPtr, err = get_action(serviceTemp, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.NotNil(t, err) // Not exposed

	// Scale the cluster
	actionPtr, err = get_action(service1Ptr, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale := *(action.GetJsonFormResultPlaceholder().(*dtos.ClusterScaleDto))
	assert.Equal(t, 0, clusterScale.NumberOfReplicas)
	// Try setting the same number of replicas as we have
	clusterScale.NumberOfReplicas = 0
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.NotNil(t, err) // Should create an error
	assert.Nil(t, result)
	// Try setting a negative number of replicas
	clusterScale.NumberOfReplicas = -1
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.NotNil(t, err) // Should create an error
	assert.Nil(t, result)
	// Increment the number of replicas
	clusterScale.NumberOfReplicas = 2
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)
	// Ensure we have 2 replicas now
	serviceTemp, err = pgProvider.GetService(user, service1.GetName())
	assert.Nil(t, err)
	actionPtr, err = get_action(serviceTemp, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale = *(action.GetJsonFormResultPlaceholder().(*dtos.ClusterScaleDto))
	assert.Equal(t, 2, clusterScale.NumberOfReplicas)
	// Decrement the number of replicas
	clusterScale.NumberOfReplicas = 1
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)
	// Ensure we have only 1 replica now
	serviceTemp, err = pgProvider.GetService(user, service1.GetName())
	assert.Nil(t, err)
	actionPtr, err = get_action(serviceTemp, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale = *(action.GetJsonFormResultPlaceholder().(*dtos.ClusterScaleDto))
	assert.Equal(t, 1, clusterScale.NumberOfReplicas)

	common_test.CommonProviderStop(t, pgProviderPtr, user)
}
