package runner

import "os"

// Job ...
type Job struct {
	Language    string
	Image       string
	Source      string
	TimeLimit   int64
	MemoryLimit int64
	Cmd         []string
	Stdout      *os.File
	Stderr      *os.File
}

// Task ...
type Task struct {
	CompileCmd []string `json:"compile"`
	ExecuteCmd []string `json:"execute"`
}
