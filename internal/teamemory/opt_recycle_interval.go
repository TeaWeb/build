package teamemory

type RecycleIntervalOpt struct {
	Interval int
}

func NewRecycleIntervalOpt(interval int) *RecycleIntervalOpt {
	return &RecycleIntervalOpt{
		Interval: interval,
	}
}
