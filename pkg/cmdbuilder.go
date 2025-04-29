package goflint

import (
	"strings"

	"goflint/pkg/common"
)

type SparkSubmit struct {
	conf            common.SparkConf
	args            []string
	application     string
	applicationArgs []string
}

func CrateOrUpdate() *SparkSubmit {
	// если base == nil, начинаем с пустого
	return &SparkSubmit{}
}

// spark://host:port, mesos://host:port, yarn,
// k8s://https://host:port, or local (Default: local[*]).
func (s SparkSubmit) WithMaster(master *string) SparkSubmit {
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
func (s SparkSubmit) WithDeployMode(mode *string) SparkSubmit {
	var deployMode string

	if mode == nil {
		deployMode = common.ClientDeployMode
	} else {
		deployMode = *mode
	}
	s.args = append(s.args, "--deploy-mode", deployMode)
	return s
}

func (s SparkSubmit) WithSparkConf(conf common.SparkConf) SparkSubmit {
	var x common.SparkConf = s.conf
	s.conf = x.Merge(conf)
	return s
}

func (s SparkSubmit) Application(application string) SparkSubmit {
	s.application = application
	return s
}

func (s SparkSubmit) ApplicationArgs(applicationArgs ...string) SparkSubmit {
	s.applicationArgs = append(s.applicationArgs, applicationArgs...)
	return s
}

func (s SparkSubmit) WithName(name string) SparkSubmit {
	s.args = append(s.args, "--name", name)
	return s
}

func (s SparkSubmit) buildArgs() string {
	if len(s.args) == 0 {
		return ""
	}
	return strings.Join(s.args, " ")
}

func (s SparkSubmit) buildAppArgs() string {
	if len(s.applicationArgs) == 0 {
		return ""
	}
	return strings.Join(s.applicationArgs, " ")
}

func (s SparkSubmit) buildConf() string {
	if s.conf.IsEmpty() {
		return ""
	}
	return strings.Join(s.conf.ToCommandLineArgs(), " ")
}

func (s SparkSubmit) Build() SparkApp {
	// args + cfg + app + apparg
	var b []string
	b = []string{s.buildArgs(), s.buildConf(), s.application, s.buildAppArgs()}
	cmd := strings.Join(b, " ")
	return SparkApp{repr: cmd}
}
