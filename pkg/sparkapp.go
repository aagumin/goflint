package goflint

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"goflint/pkg/common"
)

type SparkApp struct {
	cmd *SparkSubmit
}

func (s *SparkApp) Submit(ctx context.Context) error {
	sparkHome, present := os.LookupEnv(common.EnvSparkHome)
	if !present {
		return fmt.Errorf("env %s is not specified", common.EnvSparkHome)
	}
	command := filepath.Join(sparkHome, "bin", "spark-submit")
	fmt.Println(command, s.cmd.Repr())
	cmd := exec.Command(command, s.cmd.Repr())
	_, err := cmd.Output()
	if err != nil {
		var errorMsg string
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			errorMsg = string(exitErr.Stderr)
		}
		// The driver pod of the application already exists.
		if strings.Contains(errorMsg, common.ErrorCodePodAlreadyExists) {
			return fmt.Errorf("driver pod already exist")
		}
		if errorMsg != "" {
			return fmt.Errorf("failed to run spark-submit: %s", errorMsg)
		}
		return fmt.Errorf("failed to run spark-submit: %v", err)
	}
	return nil
}

func (s *SparkApp) Status(ctx context.Context) (string, error) {
	// TODO: Implement real status logic
	return "", fmt.Errorf("not implemented")
}

func (s *SparkApp) Kill(ctx context.Context) error {
	// TODO: Implement real kill logic
	return fmt.Errorf("not implemented yet")
}
