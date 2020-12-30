package model

type TaskRepository interface {
	FindBySubmissionID(id string) *Task
}
