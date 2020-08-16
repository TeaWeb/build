package teageo

import (
	"runtime"
	"strconv"
	"testing"
)

func TestIP2City(t *testing.T) {
	SetupDB()
	t.Log(IP2City("114.240.210.253", true))
	t.Log(IP2City("114.240.210.253", true))
	t.Log(IP2City("114.240.210.253", true))
	t.Log(IP2City("114.240.210.253", false))
	t.Log(IP2City("114.240.210.253", false))
}

func BenchmarkIP2CityNoCache(b *testing.B) {
	runtime.GOMAXPROCS(1)

	SetupDB()

	for i := 0; i < b.N; i ++ {
		IP2City("114.240."+strconv.Itoa(i%240)+"."+strconv.Itoa(i%253), false)
	}
}

func BenchmarkIP2CityCache(b *testing.B) {
	runtime.GOMAXPROCS(1)

	SetupDB()

	for i := 0; i < b.N; i ++ {
		IP2City("114.240."+strconv.Itoa(i%100)+"."+strconv.Itoa(i%253), true)
	}

	b.Log(ip2cityGrid.Stat())
}
