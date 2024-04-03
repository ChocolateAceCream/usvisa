package job

type JobGroup struct {
	USVisa USVisaJob
}

var JobGroupInstance = new(JobGroup)
