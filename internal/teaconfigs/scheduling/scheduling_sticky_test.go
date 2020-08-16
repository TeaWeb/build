package scheduling

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"testing"
)

func TestStickyScheduling_NextArgument(t *testing.T) {
	s := &StickyScheduling{}
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

	t.Log(s.mapping)

	req, err := http.NewRequest("GET", "http://www.example.com/?backend=c", nil)
	if err != nil {
		t.Fatal(err)
	}

	options := maps.Map{
		"type":  "argument",
		"param": "backend",
	}
	call := shared.NewRequestCall()
	call.Request = req
	call.Options = options
	t.Log(s.Next(call))
	t.Log(options)
}

func TestStickyScheduling_NextCookie(t *testing.T) {
	s := &StickyScheduling{}
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

	t.Log(s.mapping)

	req, err := http.NewRequest("GET", "http://www.example.com/?backend=c", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "backend",
		Value: "c",
	})

	options := maps.Map{
		"type":  "cookie",
		"param": "backend",
	}
	call := shared.NewRequestCall()
	call.Request = req
	call.Options = options
	t.Log(s.Next(call))
	t.Log(options)
}

func TestStickyScheduling_NextHeader(t *testing.T) {
	s := &StickyScheduling{}
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

	t.Log(s.mapping)

	req, err := http.NewRequest("GET", "http://www.example.com/?backend=c", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("backend", "c")

	options := maps.Map{
		"type":  "header",
		"param": "backend",
	}
	call := shared.NewRequestCall()
	call.Request = req
	call.Options = options
	t.Log(s.Next(call))
	t.Log(options)
}
