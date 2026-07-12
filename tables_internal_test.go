package hdate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// The conversion routines were rewritten to look up precomputed tables instead
// of summing month lengths in a loop. The functions below are the original
// implementations, kept verbatim as a reference so the tables can be proven
// equivalent to them over the whole practical range of the calendar, including
// the out-of-range months the old loops happened to accept.

func refDaysInMonth(month HMonth, year int) int {
	switch month {
	case Iyyar, Tamuz, Elul, Tevet, Adar2:
		return 29
	}
	if (month == Adar1 && !IsLeapYear(year)) ||
		(month == Cheshvan && !LongCheshvan(year)) ||
		(month == Kislev && ShortKislev(year)) {
		return 29
	}
	return 30
}

func refToRD(year int, month HMonth, day int) int64 {
	tempabs := int64(day)
	if month < Tishrei {
		monthsInYear := HMonth(MonthsInYear(year))
		for m := Tishrei; m <= monthsInYear; m++ {
			tempabs += int64(refDaysInMonth(m, year))
		}
		for m := Nisan; m < month; m++ {
			tempabs += int64(refDaysInMonth(m, year))
		}
	} else {
		for m := Tishrei; m < month; m++ {
			tempabs += int64(refDaysInMonth(m, year))
		}
	}
	return Epoch + elapsedDays(year) + tempabs - 1
}

func refFromRD(rataDie int64) (int, HMonth, int) {
	approx := float64(rataDie-Epoch) / avgHebrewYearDays
	year := int(approx)
	for newYear(year) <= rataDie {
		year++
	}
	year--
	var month HMonth
	if rataDie < refToRD(year, Nisan, 1) {
		month = Tishrei
	} else {
		month = Nisan
	}
	for rataDie > refToRD(year, month, refDaysInMonth(month, year)) {
		month++
	}
	day := 1 + rataDie - refToRD(year, month, 1)
	return year, month, int(day)
}

const (
	diffMinYear = 1
	diffMaxYear = 8000
)

func TestDaysInMonthMatchesReference(t *testing.T) {
	for year := diffMinYear; year <= diffMaxYear; year++ {
		// Includes the months outside [Nisan, Adar2] that the old switch
		// silently reported as 30 days.
		for month := HMonth(0); month <= Adar2+1; month++ {
			if got, want := DaysInMonth(month, year), refDaysInMonth(month, year); got != want {
				t.Fatalf("DaysInMonth(%d, %d) = %d, reference says %d", month, year, got, want)
			}
		}
	}
}

func TestToRDMatchesReference(t *testing.T) {
	for year := diffMinYear; year <= diffMaxYear; year++ {
		for month := HMonth(0); month <= Adar2+1; month++ {
			for day := 1; day <= 30; day++ {
				if got, want := ToRD(year, month, day), refToRD(year, month, day); got != want {
					t.Fatalf("ToRD(%d, %d, %d) = %d, reference says %d", year, month, day, got, want)
				}
			}
		}
	}
}

func TestFromRDMatchesReference(t *testing.T) {
	if testing.Short() {
		t.Skip("exhaustive; skipped under -short")
	}
	// The reference year search assumes its estimate never overshoots, which
	// only holds once the estimate has room to be corrected upward, so start a
	// year in rather than at the epoch. FromRD is exercised right down to the
	// epoch by the round-trip test below.
	first := newYear(diffMinYear + 1)
	last := newYear(diffMaxYear) - 1
	for rataDie := first; rataDie <= last; rataDie++ {
		hd := FromRD(rataDie)
		year, month, day := refFromRD(rataDie)
		if hd.Year() != year || hd.Month() != month || hd.Day() != day {
			t.Fatalf("FromRD(%d) = %d %d %d, reference says %d %d %d",
				rataDie, hd.Year(), hd.Month(), hd.Day(), year, month, day)
		}
	}
}

// TestRoundTripExhaustive walks every day of the calendar from the epoch
// forward and proves the two directions are exact inverses.
func TestRoundTripExhaustive(t *testing.T) {
	if testing.Short() {
		t.Skip("exhaustive; skipped under -short")
	}
	first := newYear(diffMinYear)
	last := newYear(diffMaxYear) - 1
	for rataDie := first; rataDie <= last; rataDie++ {
		hd := FromRD(rataDie)
		if got := hd.Abs(); got != rataDie {
			t.Fatalf("FromRD(%d).Abs() = %d", rataDie, got)
		}
		if got := ToRD(hd.Year(), hd.Month(), hd.Day()); got != rataDie {
			t.Fatalf("ToRD(%s) = %d, want %d", hd, got, rataDie)
		}
		// New must agree with FromRD, including the abs it caches.
		if got := New(hd.Year(), hd.Month(), hd.Day()); got != hd {
			t.Fatalf("New(%s) = %+v, want %+v", hd, got, hd)
		}
		if day := hd.Day(); day < 1 || day > hd.DaysInMonth() {
			t.Fatalf("FromRD(%d) produced out-of-range day %d", rataDie, day)
		}
	}
}

// TestDayOfYearCoverage proves every table row is fully populated: each year
// type must account for exactly its own length in days, with no gap that would
// leave a zero (invalid) month in monthByDayOfYear.
func TestDayOfYearCoverage(t *testing.T) {
	for typ := 0; typ < numYearTypes; typ++ {
		yearLen := yearLenOfType(typ)
		total := 0
		for m := Nisan; m <= Adar2; m++ {
			if m == Adar2 && yearLen <= 355 {
				continue // no Adar II in a common year
			}
			total += daysInMonth(typ, m)
		}
		assert.Equal(t, yearLen, total, "month lengths must sum to the year length for type %d", typ)
		for d := 0; d < yearLen; d++ {
			m := HMonth(monthByDayOfYear[typ][d])
			assert.True(t, Nisan <= m && m <= Adar2,
				"year type %d day %d maps to invalid month %d", typ, d, m)
		}
	}
}

func TestYearTypeRoundTrip(t *testing.T) {
	for typ := 0; typ < numYearTypes; typ++ {
		assert.Equal(t, typ, yearType(yearLenOfType(typ)))
	}
	// And every year the calendar can actually produce lands in the table.
	for year := diffMinYear; year <= diffMaxYear; year++ {
		typ := yearType(DaysInYear(year))
		assert.True(t, typ >= 0 && typ < numYearTypes, "year %d has type %d", year, typ)
		assert.Equal(t, DaysInYear(year), yearLenOfType(typ), "year %d", year)
	}
}

// TestDayOnOrBeforeNegativeRD guards the floored modulus. R.D. numbers are
// negative for Hebrew years before ~3762, and a truncated remainder there
// returned a day *after* rataDie, which is the opposite of what the function
// promises.
func TestDayOnOrBeforeNegativeRD(t *testing.T) {
	for rataDie := int64(-1373000); rataDie < -1372900; rataDie++ {
		for wday := time.Sunday; wday <= time.Saturday; wday++ {
			got := DayOnOrBefore(wday, rataDie)
			assert.LessOrEqual(t, got, rataDie, "must not return a day after rataDie")
			assert.Greater(t, got, rataDie-7, "must be within the preceding week")
			assert.Equal(t, wday, FromRD(got).Weekday(), "must land on the requested weekday")
		}
	}
}
