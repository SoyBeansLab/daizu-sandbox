package runner

// Job ...
type Job struct {
	Image       string
	Source      string
	TimeLimit   int
	MemoryLimit int
}

// NewJob ...
func NewJob(image, source string, tl, ml int) (j Job) {
	j = Job{
		Image:       image,
		Source:      source,
		TimeLimit:   tl,
		MemoryLimit: ml,
	}
	return
}
