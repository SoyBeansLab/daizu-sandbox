package model

type JudgeResult struct {
	Result          string
	AllTestCases    int
	PassedTestCases int
	ExecuteTime     int
	CodeSize        int
	CompileMessage  string
}
