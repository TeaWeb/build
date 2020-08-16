package teamemory

type LimitSizeOpt struct {
	Size int64
}

func NewLimitSizeOpt(size int64) *LimitSizeOpt {
	return &LimitSizeOpt{
		Size: size,
	}
}
