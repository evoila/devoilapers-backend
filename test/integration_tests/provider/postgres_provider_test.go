package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres"
	"OperatorAutomation/pkg/postgres/actions/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func CreatePostgresTestProvider(t *testing.T) (*postgres.PostgresProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	pgProvider := postgres.CreatePostgresProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.YamlTemplatePath,
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
	pgProvider, _ := CreatePostgresTestProvider(t)

	assert.NotEqual(t, "", pgProvider.GetServiceImage())
	assert.NotEqual(t, "", pgProvider.GetServiceDescription())
	assert.Equal(t, "Postgres", pgProvider.GetServiceType())

	testUser := unit_test.TestUser{
		KubernetesNamespace: "A_LONG_NAMESPACE",
	}

	template := *pgProvider.GetTemplate(testUser)
	assert.True(t, strings.Contains(template.GetYAML(), "namespace: "+testUser.KubernetesNamespace))
	assert.Equal(t, 1, len(template.GetImportantSections()))
	assert.Equal(t, "metadata.name", template.GetImportantSections()[0])

	template2 := *pgProvider.GetTemplate(testUser)
	assert.True(t, strings.Contains(template2.GetYAML(), "namespace: "+testUser.KubernetesNamespace))
	assert.Equal(t, 1, len(template2.GetImportantSections()))
	assert.Equal(t, "metadata.name", template2.GetImportantSections()[0])
	assert.NotEqual(t, template2.GetYAML(), template.GetYAML())
}

func get_action(service *service.IService, groupname string, actioncommand string) (*action.IAction, error) {
	actionGroups := (*service).GetActions()

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
	pgProvider, config := CreatePostgresTestProvider(t)

	user := config.Users[0]
	invalidUser := unit_test.TestUser{KubernetesNamespace: "namespace", KubernetesAccessToken: "InvalidToken"}

	// Continue with actual pg provider
	yaml := (*pgProvider.GetTemplate(user)).GetYAML()

	// Check if there is no other service
	services, err := pgProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(services))

	// Try create a service with invalid yaml
	err = pgProvider.CreateService(user, "something")
	assert.NotNil(t, err)

	// Try create a service with invalid user data
	err = pgProvider.CreateService(invalidUser, yaml)
	assert.NotNil(t, err)

	// Create a service
	err = pgProvider.CreateService(user, yaml)
	assert.Nil(t, err)

	// Try check if created with invalid user data
	services, err = pgProvider.GetServices(invalidUser)
	assert.NotNil(t, err)

	// Check if created
	services, err = pgProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	service0 := *services[0]
	assert.NotEqual(t, "", service0.GetName())
	assert.Equal(t, pgProvider.GetServiceType(), service0.GetType())
	assert.Equal(t, 3, len(service0.GetActions()))

	// Try get service with invalid user data
	_, err = pgProvider.GetService(invalidUser, service0.GetName())
	assert.NotNil(t, err)

	// Wait for service to become ok. Postgres needs some extra time.
	var service1 service.IService
	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

		// Try get service with invalid user data
		service1Ptr, err := pgProvider.GetService(user, service0.GetName())
		assert.Nil(t, err)
		service1 = *service1Ptr

		if service1.GetStatus() == service.ServiceStatusOk {
			break
		}
	}

	// Ensure service is ok
	assert.Equal(t, service.ServiceStatusOk, service1.GetStatus())

	// Ensure they have the same attributes
	assert.Equal(t, service0.GetName(), service1.GetName())
	assert.Equal(t, service0.GetType(), service1.GetType())
	assert.Equal(t, service0.GetTemplate().GetImportantSections(), service1.GetTemplate().GetImportantSections())


	// Testing actions
	// Get database credentials
	actionPtr, err := get_action(&service0, "Informations", "cmd_pg_get_credentials")
	assert.Nil(t, err)
	action := *actionPtr
	assert.Nil(t, action.GetPlaceholder())
	result, err := action.GetActionExecuteCallback()(action.GetPlaceholder())
	assert.Nil(t, err)
	clusterCredentials := *result.(*dtos.ClusterCredentialsDto)
	assert.True(t, len(clusterCredentials.Username) > 3)
	assert.True(t, len(clusterCredentials.Password) > 3)
	assert.True(t, clusterCredentials.InternalPort > 1)

	// Exposure
	actionPtr, err = get_action(&service0, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetPlaceholder())
	assert.NotNil(t, err) // Not exposed

	// Expose it
	actionPtr, err = get_action(&service0, "Security", "cmd_pg_expose")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetPlaceholder())
	assert.Nil(t, err)
	clusterExposeResult := result.(*dtos.ClusterExposeResponseDto)
	assert.True(t, clusterExposeResult.Port > 1)

	// Check again if it is exposed
	serviceTemp, err := pgProvider.GetService(user, service0.GetName())
	assert.Nil(t, err)
	actionPtr, err = get_action(serviceTemp, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation := result.(*dtos.ClusterExposeResponseDto)
	assert.Equal(t, clusterExposeResult.Port, clusterExposeInformation.Port)

	// Hide it again
	actionPtr, err = get_action(&service0, "Security", "cmd_pg_hide")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetPlaceholder())
	assert.Nil(t, err)

	// Check again if it is hidden
	serviceTemp, err = pgProvider.GetService(user, service0.GetName())
	assert.Nil(t, err)
	actionPtr, err = get_action(serviceTemp, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetPlaceholder())
	assert.NotNil(t, err) // Not exposed

	// Scale the cluster
	actionPtr, err = get_action(&service0, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale := *(action.GetPlaceholder().(*dtos.ClusterScaleDto))
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
	serviceTemp, err = pgProvider.GetService(user, service0.GetName())
	assert.Nil(t, err)
	actionPtr, err = get_action(serviceTemp, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale = *(action.GetPlaceholder().(*dtos.ClusterScaleDto))
	assert.Equal(t, 2, clusterScale.NumberOfReplicas)
	// Decrement the number of replicas
	clusterScale.NumberOfReplicas = 1
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)
	// Ensure we have only 1 replica now
	serviceTemp, err = pgProvider.GetService(user, service0.GetName())
	assert.Nil(t, err)
	actionPtr, err = get_action(serviceTemp, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale = *(action.GetPlaceholder().(*dtos.ClusterScaleDto))
	assert.Equal(t, 1, clusterScale.NumberOfReplicas)


	// Try delete service with invalid id
	err = pgProvider.DeleteService(user, "some-not-existing-id")
	assert.NotNil(t, err)

	// Try delete service with invalid user
	err = pgProvider.DeleteService(invalidUser, (*services[0]).GetName())
	assert.NotNil(t, err)

	// Delete service
	err = pgProvider.DeleteService(user, (*services[0]).GetName())
	assert.Nil(t, err)
}
