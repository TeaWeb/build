package scheduling

import (
	"sync"
	"testing"
)

type TestCandidate struct {
	Name   string
	Weight uint
}

func (this *TestCandidate) CandidateWeight() uint {
	return this.Weight
}

func (this *TestCandidate) CandidateCodes() []string {
	return []string{this.Name}
}

func TestRandomScheduling_Next(t *testing.T) {
	s := &RandomScheduling{}
	s.Add(&TestCandidate{
		Name:   "a",
		Weight: 10,
	})
	s.Add(&TestCandidate{
		Name:   "b",
		Weight: 10,
	})
	s.Add(&TestCandidate{
		Name:   "c",
		Weight: 10,
	})
	s.Add(&TestCandidate{
		Name:   "d",
		Weight: 30,
	})
	s.Start()

	/**for _, c := range s.array {
		t.Log(c.(*TestCandidate).Name, ":", c.CandidateWeight())
	}**/

	hits := map[string]uint{}
	for _, c := range s.array {
		hits[c.(*TestCandidate).Name] = 0
	}

	t.Log("count:", s.count, "array length:", len(s.array))

	var locker sync.Mutex
	var wg = sync.WaitGroup{}
	wg.Add(100 * 10000)
	for i := 0; i < 100*10000; i ++ {
		go func() {
			defer wg.Done()

			c := s.Next(nil)

			locker.Lock()
			defer locker.Unlock()
			hits[c.(*TestCandidate).Name] ++
		}()
	}
	wg.Wait()

	t.Log(hits)
}

func TestRandomScheduling_NextZero(t *testing.T) {
	s := &RandomScheduling{}
	s.Add(&TestCandidate{
		Name:   "a",
		Weight: 0,
	})
	s.Start()
	t.Log(s.Next(nil))
}
