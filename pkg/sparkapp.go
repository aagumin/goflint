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

func (s *SparkApp) Submit(ctx context.Context) (*exec.Cmd, error) {
	sparkHome, present := os.LookupEnv(common.EnvSparkHome)
	if !present {
		return nil, fmt.Errorf("env %s is not specified", common.EnvSparkHome)
	}

	command := filepath.Join(sparkHome, "bin", "spark-submit")

	args := strings.Fields(s.cmd.Repr())

	cmd := exec.CommandContext(ctx, command, args...)

	fmt.Printf("Executing: %s %s", command, strings.Join(args, " "))

	output, err := cmd.CombinedOutput()

	if err != nil {
		var errorMsg string
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			errorMsg = string(exitErr.Stderr)
		}

		if strings.Contains(errorMsg, common.ErrorCodePodAlreadyExists) {
			return nil, fmt.Errorf("driver pod already exists")
		}

		if errorMsg != "" {
			return nil, fmt.Errorf("failed to run spark-submit: %s", errorMsg)
		}
		return nil, fmt.Errorf("failed to run spark-submit: %v", err)
	}
	fmt.Printf("Spark submit succeeded: %s", string(output))
	fmt.Println(output)
	return cmd, nil
}

func (s *SparkApp) Status(ctx context.Context) (string, error) {
	// TODO: Implement real status logic
	return "", fmt.Errorf("not implemented")
}

func (s *SparkApp) Kill(ctx context.Context) error {
	// TODO: Implement real kill logic
	return fmt.Errorf("not implemented yet")
}
