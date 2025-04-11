package goflint

import (
	"goflint/pkg/common"
	"goflint/pkg/sparkconf"
	"strings"
)

type SparkSubmitCmd struct {
	conf           *sparkconf.SparkConf
	args           []string
	application    *string
	aplicationArgs []string
	command        string
}

func NewSparkSubmitCmd(base *SparkApp) *SparkSubmitCmd {
	// если base == nil, начинаем с пустого
	return &SparkSubmitCmd{}
}

// spark://host:port, mesos://host:port, yarn,
// k8s://https://host:port, or local (Default: local[*]).
func (s *SparkSubmitCmd) WithMaster(master *string) *SparkSubmitCmd {
	if master != nil {
		s.args = append(s.args, "--master", *master)
	} else {
		findedMaster, err := common.GetK8SMasterURL()
		if err != nil {
			s.args = append(s.args, "--master", findedMaster)
		} else {
			s.args = append(s.args, "--master", "local[*]")
		}
	}
	return s
}

// Whether to launch the driver program locally ("client") or
// on one of the worker machines inside the cluster ("cluster")
// (Default: client)
func (s *SparkSubmitCmd) WithDeployMode(mode *string) *SparkSubmitCmd {
	var deployMode string

	if mode == nil {
		deployMode = common.ClientDeployMode
	} else {
		deployMode = *mode
	}
	s.args = append(s.args, "--deploy-mode", deployMode)
	return s
}

func (s *SparkSubmitCmd) WithSparkConf(conf *sparkconf.SparkConf) *SparkSubmitCmd {
	s.conf = s.conf.Merge(conf)
	return s
}

func (s *SparkSubmitCmd) Application(application *string) *SparkSubmitCmd {
	s.application = application
	return s
}

func (s *SparkSubmitCmd) ApplicationArgs(applicationArgs ...string) *SparkSubmitCmd {
	s.aplicationArgs = append(s.aplicationArgs, applicationArgs...)
	return s
}

func (s *SparkSubmitCmd) buildArgs() string {
	return strings.Join(s.args, " ")
}
func (s *SparkSubmitCmd) buildAppArgs() string {
	return strings.Join(s.aplicationArgs, " ")
}
func (s *SparkSubmitCmd) buildConf() string {
	return strings.Join(s.conf.ToCommandLineArgs(), " ")
}

func (s *SparkSubmitCmd) Build() SparkApp {
	// args + cfg + app + apparg
	var b []string
	b = []string{s.buildArgs(), s.buildConf(), *s.application, s.buildAppArgs()}
	cmd := strings.Join(b, " ")
	return SparkApp{command: s, repr: cmd}
}
