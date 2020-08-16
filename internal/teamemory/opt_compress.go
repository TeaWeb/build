package teamemory

type CompressOpt struct {
	Level int
}

func NewCompressOpt(level int) *CompressOpt {
	return &CompressOpt{
		Level: level,
	}
}
