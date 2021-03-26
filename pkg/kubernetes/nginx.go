package kubernetes

import (
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"math/rand"
	"strconv"
	"sync"
)

var nginxMutex sync.Mutex

type NginxInformation struct {
	Namespace        string
	TcpConfigMapName string
	DeploymentName   string
	ContainerName    string
}

func NginxGetExposedTcpPort(api *K8sApi,
	nginxInformation NginxInformation,
	serviceNamespace string,
	serviceName string,
	servicePort int) (int, error) {

	servicePortStr := strconv.Itoa(servicePort)
	nginxIdentifier := serviceNamespace + "/" + serviceName + ":" + servicePortStr

	tcpConfigMap, err := api.ClientSet.CoreV1().ConfigMaps(nginxInformation.Namespace).Get(
		context.TODO(),
		nginxInformation.TcpConfigMapName,
		metav1.GetOptions{},
	)

	if err != nil {
		return -1, err
	}

	// Find exposed port
	for port, existingIdentifier := range tcpConfigMap.Data {
		if nginxIdentifier == existingIdentifier {
			numericPort, err := strconv.Atoi(port)
			return numericPort, err
		}
	}

	return -1, errors.New("no port found. The service is not exposed")
}

func NginxCloseTcpPort(api *K8sApi,
	nginxInformation NginxInformation,
	serviceNamespace string,
	serviceName string,
	servicePort int) error {

	// Ensure only one can open a port since we are accessing two files. There is a race condition.
	nginxMutex.Lock()
	defer nginxMutex.Unlock()

	servicePortStr := strconv.Itoa(servicePort)
	nginxIdentifier := serviceNamespace + "/" + serviceName + ":" + servicePortStr

	// Get nginx tcp map
	tcpConfigMap, err := api.ClientSet.CoreV1().ConfigMaps(nginxInformation.Namespace).Get(
		context.TODO(),
		nginxInformation.TcpConfigMapName,
		metav1.GetOptions{},
	)

	if err != nil {
		return err
	}

	// Find exposed port
	exposedPortStr := ""
	for port, existingIdentifier := range tcpConfigMap.Data {
		if nginxIdentifier == existingIdentifier {
			exposedPortStr = port
			break
		}
	}

	if exposedPortStr == "" {
		return errors.New("service is not exposed")
	}

	exposedPort, err := strconv.Atoi(exposedPortStr)
	if err != nil {
		return err
	}

	// Otherwise delete the port definition from the tcp-list
	delete(tcpConfigMap.Data, exposedPortStr)

	_, err = api.ClientSet.CoreV1().ConfigMaps(nginxInformation.Namespace).Update(
		context.TODO(),
		tcpConfigMap,
		metav1.UpdateOptions{})

	if err != nil {
		return err
	}

	// Additionally we also have to update the deployment
	deployment, err := api.ClientSet.AppsV1().Deployments(nginxInformation.Namespace).Get(
		context.TODO(),
		nginxInformation.DeploymentName,
		metav1.GetOptions{},
	)

	if err != nil {
		return err
	}

	indexToRemove := -1
	for containerIdx, container := range deployment.Spec.Template.Spec.Containers {
		// Find index of port definition
		for portIdx, portDefinition := range container.Ports {
			if int(portDefinition.HostPort) == exposedPort {
				indexToRemove = portIdx
				break
			}
		}

		// If the port definition could be found
		if indexToRemove >= 0 {
			// Remove port definition
			deployment.Spec.Template.Spec.Containers[containerIdx].Ports = append(
				deployment.Spec.Template.Spec.Containers[containerIdx].Ports[:indexToRemove],
				deployment.Spec.Template.Spec.Containers[containerIdx].Ports[indexToRemove+1:]...)

			break
		}
	}

	// Nothing to remove
	if indexToRemove < 0 {
		return nil
	}

	// Otherwise update the deployment without the port
	_, err = api.ClientSet.AppsV1().Deployments(nginxInformation.Namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{})

	return err
}

func NginxOpenRandomTcpPort(
	api *K8sApi,
	nginxInformation NginxInformation,
	serviceNamespace string,
	serviceName string,
	servicePort int) (int, error) {

	// Ensure only one can open a port since we are accessing two files. There is a race condition.
	nginxMutex.Lock()
	defer nginxMutex.Unlock()

	//TODO: Also check if udp ports are reserved
	tcpConfigMap, err := api.ClientSet.CoreV1().ConfigMaps(nginxInformation.Namespace).Get(context.TODO(), nginxInformation.TcpConfigMapName, metav1.GetOptions{})
	if err != nil {
		return -1, err
	}

	servicePortStr := strconv.Itoa(servicePort)
	nginxIdentifier := serviceNamespace + "/" + serviceName + ":" + servicePortStr

	// Check if the identifier is already exposed on a port.
	for port, existingIdentifier := range tcpConfigMap.Data {
		if nginxIdentifier == existingIdentifier {
			return -1, errors.New("service already exposed on port: " + port)
		}
	}

	// Generate a new unused port number
	minPort := 10000
	maxPort := 50000
	newPortNum := minPort
	newPort := strconv.Itoa(newPortNum)

	portFound := false
	for i := minPort; i < maxPort; i++ {
		newPortNum = rand.Intn(maxPort-minPort) + minPort
		newPort = strconv.Itoa(newPortNum)
		_, exists := tcpConfigMap.Data[newPort]

		if !exists {
			portFound = true
			break
		}
	}

	// Ensure that we have found a free port
	if !portFound {
		return -1, errors.New("could not find an unused port")
	}

	// Follow nginx convention to patch tcp-services config map
	patch := fmt.Sprintf(`
		{
		  "data": {
			"%s": "%s"
		  }
		}
	`, newPort, nginxIdentifier)

	// Patch the config map with the new port
	_, err = api.ClientSet.CoreV1().ConfigMaps(nginxInformation.Namespace).Patch(
		context.TODO(),
		nginxInformation.TcpConfigMapName,
		types.StrategicMergePatchType,
		[]byte(patch),
		metav1.PatchOptions{})

	if err != nil {
		return -1, err
	}

	// Create deployment patch as described by nginx
	deploymentPatch := fmt.Sprintf(`
			{
			  "spec": {
				"template": {
				  "spec": {
					"containers": [
					  {
						"name": "%s",
						"ports": [
						  {
							"containerPort": %s,
							"hostPort": %s
						  }
						]
					  }
					]
				  }
				}
			  }
			}
   `, nginxInformation.ContainerName, newPort, newPort)

	// Patch the new port into the deployment
	_, err = api.ClientSet.AppsV1().Deployments(nginxInformation.Namespace).Patch(
		context.TODO(),
		nginxInformation.DeploymentName,
		types.StrategicMergePatchType,
		[]byte(deploymentPatch),
		metav1.PatchOptions{})

	return newPortNum, err
}
