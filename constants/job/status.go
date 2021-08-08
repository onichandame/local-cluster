package job

type JobStatus string

const (
	PENDING  JobStatus = "PENDING"
	FINISHED JobStatus = "FINISHED"
	FAILED   JobStatus = "FAILED"
)
