package model

import "os"

type Task struct {
	Language    string
	Image       string
	Source      string
	TimeLimit   int
	MemoryLimit int
	Cmd         []string
	Stdout      *os.File
	Stderr      *os.File
}
