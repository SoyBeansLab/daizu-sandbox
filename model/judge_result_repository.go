package model

type JudgeResultRepository interface {
	UpdateJudgeResult(*JudgeResult) error
}
