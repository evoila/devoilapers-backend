package common_test

import (
	opaConfig "OperatorAutomation/cmd/service/config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func tryPathOrJoinWithWorkingDirectory(targetPath string, workingDirectory string, isDirectory bool) string {
	targetPath = path.Clean(targetPath)

	absoluteTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		absoluteTargetPath = targetPath
	}

	stats, err := os.Stat(absoluteTargetPath)
	if !os.IsNotExist(err) && stats.IsDir() == isDirectory {
		return absoluteTargetPath
	}

	newPath := path.Join(workingDirectory, targetPath)
	stats, err = os.Stat(newPath)
	if !os.IsNotExist(err) && stats.IsDir() == isDirectory {
		return newPath
	}

	panic("Path: \"" + targetPath + "\" not found during absolute path conversion")
}

func loadConfigAndResolveToAbsolutePaths(t *testing.T, pathFromRoot string) opaConfig.RawConfig {
	var config opaConfig.RawConfig
	var err error

	// Get path of this file
	_, filename, _, _ := runtime.Caller(0)
	// Navigate back to repositiory root
	fmt.Println("Loader file in: " + path.Dir(filename))

	rootDirectoryPath := path.Join(path.Dir(filename), "../../..")
	fmt.Println("Root directory at: " + rootDirectoryPath)

	// Load configuration file
	configPath := tryPathOrJoinWithWorkingDirectory(pathFromRoot, rootDirectoryPath, false)
	fmt.Println("Use config at: " + configPath)

	config, err = opaConfig.LoadConfigurationFromFile(configPath)
	assert.Nil(t, err)

	// Convert relative paths to absolute
	fmt.Println("Try resolve config.WebserverSllCertificate.PrivateKeyFilePath")
	config.WebserverSllCertificate.PrivateKeyFilePath = tryPathOrJoinWithWorkingDirectory(
		config.WebserverSllCertificate.PrivateKeyFilePath,
		rootDirectoryPath,
		false,
	)

	fmt.Println("Try resolve config.WebserverSllCertificate.PublicKeyFilePath")
	config.WebserverSllCertificate.PublicKeyFilePath = tryPathOrJoinWithWorkingDirectory(
		config.WebserverSllCertificate.PublicKeyFilePath,
		rootDirectoryPath,
		false,
	)

	fmt.Println("Try resolve config.Kubernetes.CertificateAuthority")
	config.Kubernetes.CertificateAuthority = tryPathOrJoinWithWorkingDirectory(
		config.Kubernetes.CertificateAuthority,
		rootDirectoryPath,
		false,
	)

	fmt.Println("Try resolve config.ResourcesTemplatesPath")
	config.ResourcesTemplatesPath = tryPathOrJoinWithWorkingDirectory(
		config.ResourcesTemplatesPath,
		rootDirectoryPath,
		true,
	)

	fmt.Println("Try resolveconfig.Kubernetes.Operators.Postgres.PgoCaPath")
	config.Kubernetes.Operators.Postgres.PgoCaPath = tryPathOrJoinWithWorkingDirectory(
		config.Kubernetes.Operators.Postgres.PgoCaPath,
		rootDirectoryPath,
		false,
	)

	return config
}

func GetConfig(t *testing.T) opaConfig.RawConfig {
	var config opaConfig.RawConfig
	configType := os.Getenv("ENV_GITHUB_ACTION")


	if configType == "" {
		// Local
		fmt.Println("Using local appconfig")
		config = loadConfigAndResolveToAbsolutePaths(t, "configs/appconfig.json")
	} else if configType == "remote" {
		// Remote
		fmt.Println("Using remote appconfig")
		config = loadConfigAndResolveToAbsolutePaths(t, "configs/appconfig_remote.json")
	} else {
		// Github action
		fmt.Println("Using github appconfig")
		config = loadConfigAndResolveToAbsolutePaths(t, "configs/appconfig_github_actions.json")
	}

	return config
}
