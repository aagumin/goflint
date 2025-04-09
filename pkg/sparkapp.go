package goflint

import (
	"context"
	"os/exec"
)

type SparkApp struct {
	appPath string
	args    []string
	opts    map[string]string
	flags   map[string]bool
}

// Compile-time check
var _ Submitter = (*SparkApp)(nil)
var _ Monitor = (*SparkApp)(nil)

func (s *SparkApp) Submit(ctx context.Context) error {
	cmdArgs := s.buildArgs()
	cmd := exec.CommandContext(ctx, "spark-submit", cmdArgs...)
	return cmd.Run()
}

func (s *SparkApp) Status(ctx context.Context) (string, error) {
	// TODO: Implement real status logic
	return "UNKNOWN", nil
}

func (s *SparkApp) Kill(ctx context.Context) error {
	// TODO: Implement real kill logic
	return nil
}

func (s *SparkApp) buildArgs() []string {
	args := []string{}
	for k, v := range s.opts {
		args = append(args, k, v)
	}
	for flag, enabled := range s.flags {
		if enabled {
			args = append(args, flag)
		}
	}
	args = append(args, s.appPath)
	args = append(args, s.args...)
	return args
}
