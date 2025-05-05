package flint

import (
	"strings"

	"github.com/aagumin/goflint/flint/common"
)

type SparkSubmit struct {
	conf            common.SparkConf
	args            []string
	application     string
	applicationArgs []string
}
type SparkSubmitOption func(*SparkSubmit)

func NewSparkApp(options ...SparkSubmitOption) *SparkSubmit {
	s := &SparkSubmit{}

	for _, option := range options {
		option(s)
	}

	return s
}

func ExtendSparkApp(app *SparkApp, options ...SparkSubmitOption) *SparkSubmit {
	s := &SparkSubmit{conf: app.sparkSubmit.conf, application: app.sparkSubmit.application, applicationArgs: app.sparkSubmit.applicationArgs, args: app.sparkSubmit.args}

	for _, option := range options {
		option(s)
	}

	return s
}

func WithExistApp(app *SparkApp) SparkSubmitOption {
	return func(s *SparkSubmit) {
		s.args = app.sparkSubmit.applicationArgs
		s.application = app.sparkSubmit.application
		s.applicationArgs = app.sparkSubmit.applicationArgs
		s.conf = app.sparkSubmit.conf
	}
}

// spark://host:port, mesos://host:port, yarn,
// k8s://https://host:port, or local (Default: local[*]).
func WithMaster(master string) SparkSubmitOption {
	return func(s *SparkSubmit) {
		if master != "" {
			s.args = append(s.args, "--master", master)
		} else {
			findedMaster, err := common.GetK8SMasterURL()
			if err == nil && findedMaster != "" {
				s.args = append(s.args, "--master", findedMaster)
			} else {
				s.args = append(s.args, "--master", "local")
			}
		}
	}
}

func WithMainClass(mainClass string) SparkSubmitOption {
	return func(s *SparkSubmit) {
		if mainClass != "" {
			s.args = append(s.args, "--class", mainClass)
		}
	}
}

// Whether to launch the driver program locally ("client") or
// on one of the worker machines inside the cluster ("cluster")
// (Default: client)
func WithDeployMode(mode *string) SparkSubmitOption {
	return func(s *SparkSubmit) {
		var deployMode string
		if mode == nil {
			deployMode = common.ClientDeployMode
		} else {
			deployMode = *mode
		}
		s.args = append(s.args, "--deploy-mode", deployMode)
	}
}

func WithSparkConf(conf common.SparkConf) SparkSubmitOption {
	return func(s *SparkSubmit) {
		if s.conf == nil {
			s.conf = conf
		} else {
			s.conf = s.conf.Merge(conf)
		}
	}
}

func WithApplication(application string) SparkSubmitOption {
	return func(s *SparkSubmit) {
		s.application = application
	}
}

func WithApplicationArgs(applicationArgs ...string) SparkSubmitOption {
	return func(s *SparkSubmit) {
		s.applicationArgs = append(s.applicationArgs, applicationArgs...)
	}
}

func WithName(name string) SparkSubmitOption {
	return func(s *SparkSubmit) {
		s.args = append(s.args, "--name", name)
	}
}

func (s *SparkSubmit) buildArgs() string {
	if len(s.args) == 0 {
		return ""
	}
	return strings.Join(s.args, " ")
}

func (s *SparkSubmit) buildAppArgs() string {
	if len(s.applicationArgs) == 0 {
		return ""
	}
	return strings.Join(s.applicationArgs, " ")
}

func (s *SparkSubmit) buildConf() string {
	return strings.Join(s.conf.ToCommandLineArgs(), " ")
}

func (s *SparkSubmit) Repr() string {
	b := []string{s.buildArgs(), s.buildConf(), s.application, s.buildAppArgs()}
	cmd := strings.Join(b, " ")
	return cmd
}

func (s *SparkSubmit) Build() SparkApp {
	// args + cfg + app + apparg
	return SparkApp{sparkSubmit: s}
}
