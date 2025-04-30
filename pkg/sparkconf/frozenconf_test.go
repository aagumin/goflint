package sparkconf

import (
	"goflint/pkg/common"
	"reflect"
	"strings"
	"testing"
)

func TestNewFrozenConf(t *testing.T) {
	// Test with empty data
	emptyData := map[string]string{}
	conf := NewFrozenConf(emptyData)

	if conf.props == nil {
		t.Fatal("Initialized properties map is nil")
	}

	if len(conf.props) != 0 {
		t.Errorf("Initial properties map should be empty, got %d elements", len(conf.props))
	}

	// Test with provided data
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	conf = NewFrozenConf(testData)

	if !reflect.DeepEqual(conf.props, testData) {
		t.Errorf("Properties map should match the provided data, got %v, expected %v", conf.props, testData)
	}
}

func TestFrozenSparkConf_Get(t *testing.T) {
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	conf := NewFrozenConf(testData)

	// Check existing key
	if val := conf.Get("key1"); val != "value1" {
		t.Errorf("Get('key1') returned %q, expected 'value1'", val)
	}

	// Check existing key
	if val := conf.Get("key2"); val != "value2" {
		t.Errorf("Get('key2') returned %q, expected 'value2'", val)
	}

	// Check non-existent key
	if val := conf.Get("nonexistent"); val != "" {
		t.Errorf("Get() for non-existent key should return empty string, got %q", val)
	}
}

func TestFrozenSparkConf_GetAll(t *testing.T) {
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	conf := NewFrozenConf(testData)
	props := conf.GetAll()

	// Check that all elements are returned
	if !reflect.DeepEqual(props, testData) {
		t.Errorf("GetAll() returned %v, expected %v", props, testData)
	}

	// Check that modifying the returned map doesn't affect the original
	props["key3"] = "value3"
	if conf.Contains("key3") {
		t.Error("Modifying the map returned by GetAll() should not affect the original configuration")
	}

	// Check that original data remains unchanged
	originalProps := conf.GetAll()
	if reflect.DeepEqual(originalProps, props) {
		t.Error("GetAll() should return a copy of the properties map")
	}
}

func TestFrozenSparkConf_Contains(t *testing.T) {
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	conf := NewFrozenConf(testData)

	// Check existing keys
	if !conf.Contains("key1") {
		t.Error("Contains('key1') returned false, expected true")
	}

	if !conf.Contains("key2") {
		t.Error("Contains('key2') returned false, expected true")
	}

	// Check non-existent key
	if conf.Contains("nonexistent") {
		t.Error("Contains('nonexistent') returned true, expected false")
	}
}

func TestFrozenSparkConf_ToCommandLineArgs(t *testing.T) {
	tests := []struct {
		name     string
		setup    map[string]string
		expected []string
	}{
		{
			name:     "empty configuration",
			setup:    map[string]string{},
			expected: []string{},
		},
		{
			name: "only spark.* keys",
			setup: map[string]string{
				"spark.app.name":        "MyApp",
				"spark.executor.memory": "2g",
			},
			expected: []string{
				"--conf", "spark.app.name=MyApp",
				"--conf", "spark.executor.memory=2g",
			},
		},
		{
			name: "mixed keys",
			setup: map[string]string{
				"spark.app.name": "MyApp",
				"non.spark.key":  "value",
			},
			expected: []string{
				"--conf", "spark.app.name=MyApp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := NewFrozenConf(tt.setup)
			args := conf.ToCommandLineArgs()

			// Check number of arguments
			if len(args) != len(tt.expected) {
				t.Fatalf("ToCommandLineArgs() returned %d arguments, expected %d: %v",
					len(args), len(tt.expected), args)
			}

			// Order of arguments is not important as they can be in different order due to map iteration
			// Just check that all pairs "--conf key=value" are present
			for i := 0; i < len(args); i += 2 {
				if i+1 >= len(args) {
					t.Fatalf("Odd number of arguments in result: %v", args)
				}

				if args[i] != "--conf" {
					t.Errorf("Expected argument '--conf', got %q", args[i])
				}

				found := false
				for j := 0; j < len(tt.expected); j += 2 {
					if j+1 < len(tt.expected) && args[i+1] == tt.expected[j+1] {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Argument %q not found in expected list", args[i+1])
				}
			}
		})
	}
}

func TestFrozenSparkConf_IsEmpty(t *testing.T) {
	// Test with empty configuration
	emptyConf := NewFrozenConf(map[string]string{})
	if !emptyConf.IsEmpty() {
		t.Error("IsEmpty() for empty configuration returned false, expected true")
	}

	// Test with non-empty configuration
	nonEmptyConf := NewFrozenConf(map[string]string{"key": "value"})
	if nonEmptyConf.IsEmpty() {
		t.Error("IsEmpty() for non-empty configuration returned true, expected false")
	}
}

func TestFrozenSparkConf_Merge(t *testing.T) {
	// Create source configuration
	srcData := map[string]string{
		"key1":   "value1",
		"common": "srcValue",
	}
	srcConf := NewFrozenConf(srcData)

	// Create configuration to merge
	otherData := map[string]string{
		"key2":   "value2",
		"common": "otherValue",
	}
	otherConf := NewFrozenConf(otherData)

	// Merge configurations
	result := srcConf.Merge(otherConf)

	// Check that result implements SparkConf interface
	_, ok := result.(common.SparkConf)
	if !ok {
		t.Error("Merge() should return an object implementing common.SparkConf interface")
	}

	// Cast result to FrozenSparkConf for testing
	resultConf, ok := result.(FrozenSparkConf)
	if !ok {
		t.Error("Result of Merge() should be of type FrozenSparkConf")
	}

	// Check contents of merged configuration
	expected := map[string]string{
		"key1":   "value1",
		"key2":   "value2",
		"common": "otherValue", // Should be overwritten by otherConf value
	}

	if !reflect.DeepEqual(resultConf.GetAll(), expected) {
		t.Errorf("After Merge() got %v, expected %v", resultConf.GetAll(), expected)
	}

	// Check that source configuration is not modified
	if !reflect.DeepEqual(srcConf.GetAll(), srcData) {
		t.Errorf("Source configuration should not be modified after Merge()")
	}

	// Check that other configuration is not modified
	if !reflect.DeepEqual(otherConf.GetAll(), otherData) {
		t.Errorf("Other configuration should not be modified after Merge()")
	}
}

// Test interface implementation
func TestFrozenSparkConf_Interface(t *testing.T) {
	// Check that FrozenSparkConf implements SparkConf interface
	var conf common.SparkConf = NewFrozenConf(map[string]string{})

	// If code compiles, the interface is implemented correctly
	// But also check a few methods

	testData := map[string]string{
		"key1":           "value1",
		"spark.app.name": "TestApp",
	}

	conf = NewFrozenConf(testData)

	// Check Get
	if val := conf.Get("key1"); val != "value1" {
		t.Errorf("Through interface Get('key1') returned %q, expected 'value1'", val)
	}

	// Check GetAll
	allProps := conf.GetAll()
	if !reflect.DeepEqual(allProps, testData) {
		t.Errorf("Through interface GetAll() returned %v, expected %v", allProps, testData)
	}

	// Check ToCommandLineArgs
	args := conf.ToCommandLineArgs()
	if len(args) != 2 {
		t.Errorf("Through interface ToCommandLineArgs() returned wrong number of arguments: %v", args)
	}
}

func TestFrozenSparkConf_Repr(t *testing.T) {
	// Test with empty configuration
	emptyConf := NewFrozenConf(map[string]string{})
	emptyRepr := emptyConf.Repr()
	if emptyRepr != "" {
		t.Errorf("Repr() for empty configuration returned %q, expected \"{}\"", emptyRepr)
	}

	// Test with single key-value pair
	singleConf := NewFrozenConf(map[string]string{"spark.key1": "value1"})
	singleRepr := singleConf.Repr()
	if singleRepr != "--conf spark.key1=value1" {
		t.Errorf("Repr() for single key configuration returned %q, expected \"--conf spark.key1=value1\"", singleRepr)
	}

	// Test with multiple key-value pairs
	// Since map iteration order is not guaranteed, we need to check if the representation contains all entries
	multiConf := NewFrozenConf(map[string]string{
		"spark.key1": "value1",
		"spark.key2": "value2",
		"spark.key3": "value3",
	})
	multiRepr := multiConf.Repr()

	// Check that the representation starts with { and ends with }
	if !strings.HasPrefix(multiRepr, "--conf") || !strings.HasSuffix(multiRepr, "value3") {
		t.Errorf("Repr() should start with `--conf` and end with `value3`, got %q", multiRepr)
	}

}
