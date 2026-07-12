package hdate_test

import (
	"testing"
	"time"

	"github.com/hebcal/hdate"
)

var (
	sinkRD    int64
	sinkHD    hdate.HDate
	sinkInt   int
	sinkBool  bool
	sinkStr   string
	sinkMonth hdate.HMonth
	sinkWday  time.Weekday
)

func BenchmarkToRD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkRD = hdate.ToRD(5765+i%50, hdate.Sivan, 15)
	}
}

func BenchmarkFromRD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkHD = hdate.FromRD(733000 + int64(i%10000))
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkHD = hdate.New(5765+i%50, hdate.Sivan, 15)
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hd := hdate.FromRD(733000 + int64(i%10000))
		sinkRD = hd.Abs()
	}
}

func BenchmarkFromGregorian(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkHD = hdate.FromGregorian(2008, time.November, 1+i%28)
	}
}

func BenchmarkGreg(b *testing.B) {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkInt, _, _ = hd.Greg()
	}
}

func BenchmarkDaysInYear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkInt = hdate.DaysInYear(5765 + i%50)
	}
}

func BenchmarkDaysInMonth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkInt = hdate.DaysInMonth(hdate.Kislev, 5765+i%50)
	}
}

func BenchmarkIsLeapYear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkBool = hdate.IsLeapYear(5765 + i%50)
	}
}

func BenchmarkWeekday(b *testing.B) {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkWday = hd.Weekday()
	}
}

func BenchmarkNext(b *testing.B) {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkHD = hd.Next()
	}
}

func BenchmarkMonthNameEn(b *testing.B) {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkStr = hd.MonthName("en")
	}
}

func BenchmarkMonthNameHe(b *testing.B) {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkStr = hd.MonthName("he")
	}
}

func BenchmarkMonthFromName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkMonth, _ = hdate.MonthFromName("Cheshvan")
	}
}

func BenchmarkString(b *testing.B) {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkStr = hd.String()
	}
}

// 5767 has a long Cheshvan and 5765 is a leap year, so both benchmarks below
// exercise the month-shifting branches rather than the trivial path.
func BenchmarkGetYahrzeit(b *testing.B) {
	hd := hdate.New(5767, hdate.Cheshvan, 30)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkHD, _ = hdate.GetYahrzeit(5768+i%50, hd)
	}
}

func BenchmarkGetBirthdayOrAnniversary(b *testing.B) {
	hd := hdate.New(5765, hdate.Adar1, 30)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkHD, _ = hdate.GetBirthdayOrAnniversary(5766+i%50, hd)
	}
}

func BenchmarkHebrewStripNikkud(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkStr = hdate.HebrewStripNikkud("בְּרֵאשִׁית")
	}
}

func BenchmarkHebrewStripNikkudClean(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkStr = hdate.HebrewStripNikkud("בראשית")
	}
}
