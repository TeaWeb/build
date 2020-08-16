package teautils

import (
	"fmt"
	"strconv"
	"testing"
)

func TestParseVariables(t *testing.T) {
	v := ParseVariables("hello, ${name}, world", func(s string) string {
		return "Lu"
	})
	t.Log(v)
}

func TestParseNoVariables(t *testing.T) {
	for i := 0; i < 2; i++ {
		v := ParseVariables("hello, world", func(s string) string {
			return "Lu"
		})
		t.Log(v)
	}
}

func BenchmarkParseVariables(b *testing.B) {
	_ = ParseVariables("hello, ${name}, ${age}, ${gender}, ${home}, world", func(s string) string {
		return "Lu"
	})

	for i := 0; i < b.N; i++ {
		_ = ParseVariables("hello, ${name}, ${age}, ${gender}, ${home}, world", func(s string) string {
			return "Lu"
		})
	}
}

func BenchmarkParseVariablesUnique(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ParseVariables("hello, ${name} "+strconv.Itoa(i%1000), func(s string) string {
			return "Lu"
		})
	}
}

func BenchmarkParseNoVariables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ParseVariables("hello, world, "+fmt.Sprintf("%d", i%1000), func(s string) string {
			return "Lu"
		})
	}
}
