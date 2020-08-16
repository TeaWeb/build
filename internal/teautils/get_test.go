package teautils

import "testing"

func TestGetStruct(t *testing.T) {
	object := struct {
		Name  string
		Age   int
		Books []string
		Extend struct {
			Location struct {
				City string
			}
		}
	}{
		Name:  "lu",
		Age:   20,
		Books: []string{"Golang"},
		Extend: struct {
			Location struct {
				City string
			}
		}{
			Location: struct {
				City string
			}{
				City: "Beijing",
			},
		},
	}

	if Get(object, []string{"Name"}) != "lu" {
		t.Fatal("[ERROR]Name != lu")
	}

	if Get(object, []string{"Age"}) != 20 {
		t.Fatal("[ERROR]Age != 20")
	}

	if Get(object, []string{"Books", "0"}) != "Golang" {
		t.Fatal("[ERROR]books.0 != Golang")
	}

	t.Log("Extend.Location:", Get(object, []string{"Extend", "Location"}))

	if Get(object, []string{"Extend", "Location", "City"}) != "Beijing" {
		t.Fatal("[ERROR]Extend.Location.City != Beijing")
	}
}

func TestGetMap(t *testing.T) {
	object := map[string]interface{}{
		"Name": "lu",
		"Age":  20,
		"Extend": map[string]interface{}{
			"Location": map[string]interface{}{
				"City": "Beijing",
			},
		},
	}

	if Get(object, []string{"Name"}) != "lu" {
		t.Fatal("[ERROR]Name != lu")
	}

	if Get(object, []string{"Age"}) != 20 {
		t.Fatal("[ERROR]Age != 20")
	}

	if Get(object, []string{"Books", "0"}) != nil {
		t.Fatal("[ERROR]books.0 != nil")
	}

	t.Log(Get(object, []string{"Extend", "Location"}))

	if Get(object, []string{"Extend", "Location", "City"}) != "Beijing" {
		t.Fatal("[ERROR]Extend.Location.City != Beijing")
	}
}
