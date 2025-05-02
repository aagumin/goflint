package sparkconf

import (
	"github.com/aagumin/goflint/flint/common"
	"strings"
)

type FrozenSparkConf struct {
	props map[string]string
}

// NewSparkConf creates a new SparkConf instance. If props empty, get deafut system propertiers.
func NewFrozenConf(data map[string]string) FrozenSparkConf {
	// TODO load spark-default.conf
	return FrozenSparkConf{
		props: data,
	}
}

// Get returns a configuration value
func (conf FrozenSparkConf) Get(key string) string {
	return conf.props[key]
}

// GetAll returns all configuration properties
func (conf FrozenSparkConf) GetAll() map[string]string {
	// Return a copy of the map to prevent external modifications
	copyd := make(map[string]string, len(conf.props))
	for k, v := range conf.props {
		copyd[k] = v
	}
	return copyd
}

// Contains checks if a configuration exists
func (conf FrozenSparkConf) Contains(key string) bool {
	_, exists := conf.props[key]
	return exists
}

// ToCommandLineArgs converts the configuration to command-line arguments
// for spark-submit in the format: --conf key=value
func (conf FrozenSparkConf) ToCommandLineArgs() []string {
	args := make([]string, 0, len(conf.props))
	for k, v := range conf.props {
		if strings.HasPrefix(k, "spark.") {
			args = append(args, "--conf", k+"="+v)
		}
	}
	return args
}

func (conf FrozenSparkConf) IsEmpty() bool {
	return len(conf.props) == 0
}

// Merge combines another SyncSparkConf into this one and return NEW FrozenSparkConf
func (conf FrozenSparkConf) Merge(other common.SparkConf) common.SparkConf {
	a := conf.props

	for k, v := range other.GetAll() {
		a[k] = v
	}
	return NewFrozenConf(a)
}

func (conf FrozenSparkConf) Repr() string {
	return strings.Join(conf.ToCommandLineArgs(), " ")
}
