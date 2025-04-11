package sparkconf

import (
	"strings"
	"sync"
)

type SparkConf struct {
	mu    sync.RWMutex
	props map[string]string
}

// NewSparkConf creates a new SparkConf instance
func NewSparkConf() *SparkConf {
	return &SparkConf{
		props: make(map[string]string),
	}
}

// Set sets a configuration property
func (conf *SparkConf) Set(key, value string) *SparkConf {
	conf.mu.Lock()
	defer conf.mu.Unlock()
	conf.props[key] = value
	return conf
}

// Get returns a configuration value
func (conf *SparkConf) Get(key string) string {
	conf.mu.RLock()
	defer conf.mu.RUnlock()
	return conf.props[key]
}

// GetAll returns all configuration properties
func (conf *SparkConf) GetAll() map[string]string {
	conf.mu.RLock()
	defer conf.mu.RUnlock()

	// Return a copy of the map to prevent external modifications
	copyd := make(map[string]string, len(conf.props))
	for k, v := range conf.props {
		copyd[k] = v
	}
	return copyd
}

// Contains checks if a configuration exists
func (conf *SparkConf) Contains(key string) bool {
	conf.mu.RLock()
	defer conf.mu.RUnlock()
	_, exists := conf.props[key]
	return exists
}

// ToCommandLineArgs converts the configuration to command-line arguments
// for spark-submit in the format: --conf key=value
func (conf *SparkConf) ToCommandLineArgs() []string {
	conf.mu.RLock()
	defer conf.mu.RUnlock()

	args := make([]string, 0, len(conf.props))
	for k, v := range conf.props {
		if strings.HasPrefix(k, "spark.") {
			args = append(args, "--conf", k+"="+v)
		}
	}
	return args
}

// Remove deletes a configuration property
func (conf *SparkConf) Remove(key string) *SparkConf {
	conf.mu.Lock()
	defer conf.mu.Unlock()
	delete(conf.props, key)
	return conf
}

// SetIfMissing sets a configuration property if not already set
func (conf *SparkConf) SetIfMissing(key, value string) *SparkConf {
	conf.mu.Lock()
	defer conf.mu.Unlock()

	if _, exists := conf.props[key]; !exists {
		conf.props[key] = value
	}
	return conf
}

// Merge combines another SparkConf into this one
func (conf *SparkConf) Merge(other *SparkConf) *SparkConf {
	other.mu.RLock()
	defer other.mu.RUnlock()
	conf.mu.Lock()
	defer conf.mu.Unlock()

	for k, v := range other.props {
		conf.props[k] = v
	}
	return conf
}

//conf := spark.NewSparkConf().
//SetAppName("My Spark App").
//SetMaster("local[*]").
//Set("spark.executor.memory", "2g").
//Set("spark.driver.memory", "1g")
//
//args := conf.ToCommandLineArgs()
//// Output: ["--conf", "spark.app.name=My Spark App",
////          "--conf", "spark.master=local[*]",
////          "--conf", "spark.executor.memory=2g",
////          "--conf", "spark.driver.memory=1g"]
//
//defaultConf := spark.NewSparkConf().
//Set("spark.default.parallelism", "8").
//Set("spark.serializer", "org.apache.spark.serializer.KryoSerializer")
//
//customConf := spark.NewSparkConf().
//SetAppName("Custom App").
//SetMaster("yarn")
//
//mergedConf := defaultConf.Merge(customConf)

//Add validation for common Spark properties
//
//Implement environment variable support
//
//Add YAML/JSON serialization
//
//Create helper methods for common configurations:
//
//func (conf *SparkConf) EnableDynamicAllocation() *SparkConf {
//	return conf.
//		Set("spark.dynamicAllocation.enabled", "true").
//		Set("spark.shuffle.service.enabled", "true")
//}
