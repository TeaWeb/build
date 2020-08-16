package scheduling

import "testing"

func TestRoundRobinScheduling_Next(t *testing.T) {
	s := &RoundRobinScheduling{}
	s.Add(&TestCandidate{
		Name:   "a",
		Weight: 5,
	})
	s.Add(&TestCandidate{
		Name:   "b",
		Weight: 10,
	})
	s.Add(&TestCandidate{
		Name:   "c",
		Weight: 20,
	})
	s.Add(&TestCandidate{
		Name:   "d",
		Weight: 30,
	})
	s.Start()

	for _, c := range s.Candidates {
		t.Log(c.(*TestCandidate).Name, c.CandidateWeight())
	}

	t.Log(s.currentWeights)

	for i := 0; i < 100; i ++ {
		t.Log("===", "round", i, "===")
		t.Log(s.Next(nil))
		t.Log(s.currentWeights)
		t.Log(s.rawWeights)
	}
}

func TestRoundRobinScheduling_Two(t *testing.T) {
	s := &RoundRobinScheduling{}
	s.Add(&TestCandidate{
		Name:   "a",
		Weight: 10,
	})
	s.Add(&TestCandidate{
		Name:   "b",
		Weight: 10,
	})
	s.Start()

	for _, c := range s.Candidates {
		t.Log(c.(*TestCandidate).Name, c.CandidateWeight())
	}

	t.Log(s.currentWeights)

	for i := 0; i < 100; i ++ {
		t.Log("===", "round", i, "===")
		t.Log(s.Next(nil))
		t.Log(s.currentWeights)
		t.Log(s.rawWeights)
	}
}

func TestRoundRobinScheduling_NextPerformance(t *testing.T) {
	s := &RoundRobinScheduling{}
	s.Add(&TestCandidate{
		Name:   "a",
		Weight: 1,
	})
	s.Add(&TestCandidate{
		Name:   "b",
		Weight: 2,
	})
	s.Add(&TestCandidate{
		Name:   "c",
		Weight: 3,
	})
	s.Add(&TestCandidate{
		Name:   "d",
		Weight: 6,
	})
	s.Start()

	for _, c := range s.Candidates {
		t.Log(c.(*TestCandidate).Name, c.CandidateWeight())
	}

	t.Log(s.currentWeights)

	hits := map[string]uint{}
	for _, c := range s.Candidates {
		hits[c.(*TestCandidate).Name] = 0
	}
	for i := 0; i < 100*10000; i ++ {
		c := s.Next(nil)
		hits[c.(*TestCandidate).Name] ++
	}

	t.Log(hits)
}
