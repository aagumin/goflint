package sparkconf

import (
	"goflint/pkg/common"
	"strings"
	"sync"
)

type SyncSparkConf struct {
	mu    sync.RWMutex
	props map[string]string
}

// NewSparkConf creates a new SyncSparkConf instance
func NewSparkConf() *SyncSparkConf {
	return &SyncSparkConf{
		props: make(map[string]string),
	}
}

// Set sets a configuration property
func (conf *SyncSparkConf) Set(key, value string) *SyncSparkConf {
	conf.mu.Lock()
	defer conf.mu.Unlock()
	conf.props[key] = value
	return conf
}

// Get returns a configuration value
func (conf *SyncSparkConf) Get(key string) string {
	conf.mu.RLock()
	defer conf.mu.RUnlock()
	return conf.props[key]
}

// GetAll returns all configuration properties
func (conf *SyncSparkConf) GetAll() map[string]string {
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
func (conf *SyncSparkConf) Contains(key string) bool {
	conf.mu.RLock()
	defer conf.mu.RUnlock()
	_, exists := conf.props[key]
	return exists
}

// ToCommandLineArgs converts the configuration to command-line arguments
// for spark-submit in the format: --conf key=value
func (conf *SyncSparkConf) ToCommandLineArgs() []string {
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
func (conf *SyncSparkConf) Remove(key string) *SyncSparkConf {
	conf.mu.Lock()
	defer conf.mu.Unlock()
	delete(conf.props, key)
	return conf
}

// SetIfMissing sets a configuration property if not already set
func (conf *SyncSparkConf) SetIfMissing(key, value string) *SyncSparkConf {
	conf.mu.Lock()
	defer conf.mu.Unlock()

	if _, exists := conf.props[key]; !exists {
		conf.props[key] = value
	}
	return conf
}

// Merge combines another SyncSparkConf into this one
func (conf *SyncSparkConf) Merge(other common.SparkConf) common.SparkConf {
	conf.mu.Lock()
	defer conf.mu.Unlock()

	for k, v := range other.GetAll() {
		conf.props[k] = v
	}
	return conf
}
