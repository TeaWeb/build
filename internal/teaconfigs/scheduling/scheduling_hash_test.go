package scheduling

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"math/rand"
	"testing"
	"time"
)

func TestHashScheduling_Next(t *testing.T) {
	s := &HashScheduling{}
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

	hits := map[string]uint{}
	for _, c := range s.Candidates {
		hits[c.(*TestCandidate).Name] = 0
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000000; i ++ {
		call := shared.NewRequestCall()
		call.Options["key"] = "192.168.1." + fmt.Sprintf("%d", rand.Int())

		c := s.Next(call)
		hits[c.(*TestCandidate).Name] ++
	}
	t.Log(hits)
}
