package teamemory

type LimitCountOpt struct {
	Count int
}

func NewLimitCountOpt(count int) *LimitCountOpt {
	return &LimitCountOpt{
		Count: count,
	}
}
