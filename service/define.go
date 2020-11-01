package service

type AetherService interface {

	// Service Init entry
	Regist(args ...interface{}) error
}

type AetherJobService interface {
	AetherService

	OnStart()

	OnFinish()
}

type AetherJobType int

const (
	AetherJobTypeOnce AetherJobType = 1 //Job only run one time
	AetherJobTypeLoop AetherJobType = 2 //Tasks with timed cycles
)

type AetherTaskStatus int

const (
	AetherTaskStatusRunning AetherTaskStatus = 1 // Running
	AetherTaskStatusFailed  AetherTaskStatus = 2 // Failed
	AetherTaskStatusSucceed AetherTaskStatus = 3 // Succeed
)

type AetherJobStatus int

const (
	AetherJobStatusUnknown  AetherJobStatus = 0 // Unknown
	AetherJobStatusRegisted AetherJobStatus = 1 // Registed
	AetherJobStatusStarted  AetherJobStatus = 2 // Started
	AetherJobStatusCfgError AetherJobStatus = 3 // Bad Config
	AetherJobStatusRunning  AetherJobStatus = 4 // Running
	AetherJobStatusFree     AetherJobStatus = 5 // Free
	AetherJobStatusStopped  AetherJobStatus = 6 // Stopped
)

type AetherJobFunc func() error
