package runner

// Job ...
type Job struct {
	Language    string
	Image       string
	Source      string
	TimeLimit   int64
	MemoryLimit int64
	Cmd         []string
}

// Task ...
type Task struct {
	CompileCmd []string `json:"compile"`
	ExecuteCmd []string `json:"execute"`
}
