package model

import "os"

type Task struct {
	SubmissionID string
	ProblemID    string
	Username     string
	Language     string
	Image        string
	Source       string
	TimeLimit    int
	MemoryLimit  int
	Cmd          []string
	Stdout       *os.File
	Stderr       *os.File
}
