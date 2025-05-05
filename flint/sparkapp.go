package flint

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aagumin/goflint/flint/common"
)

type SparkApp struct {
	sparkSubmit *SparkSubmit
	cmd         *exec.Cmd
	isRunning   bool
}

func (s *SparkApp) Submit(ctx context.Context) ([]byte, error) {
	if s.isRunning {
		return nil, errors.New("already running")
	}
	sparkHome, present := os.LookupEnv(common.EnvSparkHome)
	if !present {
		return nil, fmt.Errorf("env %s is not specified", common.EnvSparkHome)
	}

	command := filepath.Join(sparkHome, "bin", "spark-submit")

	args := strings.Fields(s.sparkSubmit.Repr())

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
		return nil, fmt.Errorf("failed to run spark-submit: %v, \n %s", err, string(output))
	}
	fmt.Printf("Spark submit succeeded: %s", string(output))
	s.isRunning = true
	return output, nil
}

func (s *SparkApp) Status(ctx context.Context) (string, error) {
	if s.cmd != nil {
		return "alreay done or cancel", nil
	}

	return "done", nil
}

func (s *SparkApp) Kill(ctx context.Context) error {
	err := s.cmd.Cancel()
	return err
}
