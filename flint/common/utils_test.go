package common

import (
	"os"
	"testing"
)

func unSetFix(originalHostEnv string, originalPortEnv string) error {
	err := os.Setenv(EnvKubernetesServiceHost, originalHostEnv)
	if err != nil {
		return err
	}
	err = os.Setenv(EnvKubernetesServicePort, originalPortEnv)
	if err != nil {
		return err
	}
	return nil
}

func TestGetK8SMasterURL(t *testing.T) {
	originalHostEnv := os.Getenv(EnvKubernetesServiceHost)
	originalPortEnv := os.Getenv(EnvKubernetesServicePort)

	defer unSetFix(originalHostEnv, originalPortEnv)

	tests := []struct {
		name         string
		hostEnvValue string
		portEnvValue string
		expectedURL  string
		expectError  bool
	}{
		{
			name:         "IPv4 address",
			hostEnvValue: "192.168.1.1",
			portEnvValue: "443",
			expectedURL:  "k8s://https://192.168.1.1:443",
			expectError:  false,
		},
		{
			name:         "IPv6 address",
			hostEnvValue: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			portEnvValue: "6443",
			expectedURL:  "k8s://https://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:6443",
			expectError:  false,
		},
		{
			name:         "Domain name",
			hostEnvValue: "kubernetes.default.svc",
			portEnvValue: "443",
			expectedURL:  "k8s://https://kubernetes.default.svc:443",
			expectError:  false,
		},
		{
			name:         "IPv6 in brackets",
			hostEnvValue: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
			portEnvValue: "6443",
			expectedURL:  "k8s://https://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:6443",
			expectError:  false,
		},
		{
			name:         "Missing host env",
			hostEnvValue: "",
			portEnvValue: "443",
			expectedURL:  "",
			expectError:  true,
		},
		{
			name:         "Missing port env",
			hostEnvValue: "192.168.1.1",
			portEnvValue: "",
			expectedURL:  "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv(EnvKubernetesServiceHost, tt.hostEnvValue)
			if err != nil {
				return
			}
			err = os.Setenv(EnvKubernetesServicePort, tt.portEnvValue)
			if err != nil {
				return
			}

			url, err := GetK8SMasterURL()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if url != tt.expectedURL {
				t.Errorf("Expected URL %q, got %q", tt.expectedURL, url)
			}
		})
	}
}

// TestGetK8SMasterURLWithUnsetEnv проверяет случай, когда переменные среды не установлены вообще
func TestGetK8SMasterURLWithUnsetEnv(t *testing.T) {
	originalHostEnv := os.Getenv(EnvKubernetesServiceHost)
	originalPortEnv := os.Getenv(EnvKubernetesServicePort)

	defer unSetFix(originalHostEnv, originalPortEnv)

	err := os.Unsetenv(EnvKubernetesServiceHost)
	if err != nil {
		return
	}
	err = os.Unsetenv(EnvKubernetesServicePort)
	if err != nil {
		return
	}

	url, err := GetK8SMasterURL()

	if err == nil {
		t.Error("Expected error when environment variables are not set, but got none")
	}

	if url != "" {
		t.Errorf("Expected empty URL, got %q", url)
	}
}
