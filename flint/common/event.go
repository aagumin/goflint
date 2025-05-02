package common

// SparkApplication events
const (
	EventSparkApplicationAdded = "SparkApplicationAdded"

	EventSparkApplicationSubmitted = "SparkApplicationSubmitted"

	EventSparkApplicationSubmissionFailed = "SparkApplicationSubmissionFailed"

	EventSparkApplicationCompleted = "SparkApplicationCompleted"

	EventSparkApplicationFailed = "SparkApplicationFailed"

	EventSparkApplicationPendingRerun = "SparkApplicationPendingRerun"
)

// Spark driver events
const (
	EventSparkDriverPending = "SparkDriverPending"

	EventSparkDriverRunning = "SparkDriverRunning"

	EventSparkDriverCompleted = "SparkDriverCompleted"

	EventSparkDriverFailed = "SparkDriverFailed"

	EventSparkDriverUnknown = "SparkDriverUnknown"
)

// Spark executor events
const (
	EventSparkExecutorPending = "SparkExecutorPending"

	EventSparkExecutorRunning = "SparkExecutorRunning"

	EventSparkExecutorCompleted = "SparkExecutorCompleted"

	EventSparkExecutorFailed = "SparkExecutorFailed"

	EventSparkExecutorUnknown = "SparkExecutorUnknown"
)

const (
	ErrorCodePodAlreadyExists = "code=409"
)

const (
	SparkApplicationFinalizerName          = "sparkoperator.k8s.io/finalizer"
	ScheduledSparkApplicationFinalizerName = "sparkoperator.k8s.io/finalizer"
)

const (
	RSAKeySize = 2048
)

const (
	CAKeyPem      = "ca-key.pem"
	CACertPem     = "ca-cert.pem"
	ServerKeyPem  = "server-key.pem"
	ServerCertPem = "server-cert.pem"
)

// Kubernetes volume types.
const (
	VolumeTypeEmptyDir              = "emptyDir"
	VolumeTypeHostPath              = "hostPath"
	VolumeTypeNFS                   = "nfs"
	VolumeTypePersistentVolumeClaim = "persistentVolumeClaim"
)

const (
	// Epsilon is a small number used to compare 64 bit floating point numbers.
	Epsilon = 1e-9
)
