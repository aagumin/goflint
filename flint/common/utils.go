package common

import (
	"fmt"
	"os"
	"strings"
)

func GetK8SMasterURL() (string, error) {
	kubernetesServiceHost := os.Getenv(EnvKubernetesServiceHost)
	if kubernetesServiceHost == "" {
		return "", fmt.Errorf("environment variable %s is not found", EnvKubernetesServiceHost)
	}

	kubernetesServicePort := os.Getenv(EnvKubernetesServicePort)
	if kubernetesServicePort == "" {
		return "", fmt.Errorf("environment variable %s is not found", EnvKubernetesServicePort)
	}
	// check if the host is IPv6 address
	if strings.Contains(kubernetesServiceHost, ":") && !strings.HasPrefix(kubernetesServiceHost, "[") {
		return fmt.Sprintf("k8s://https://[%s]:%s", kubernetesServiceHost, kubernetesServicePort), nil
	}
	return fmt.Sprintf("k8s://https://%s:%s", kubernetesServiceHost, kubernetesServicePort), nil
}
