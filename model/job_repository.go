package model

type JobRepository interface {
	GetJob() (*Job, error)
}
