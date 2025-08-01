package hdate_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/stretchr/testify/assert"
)

func TestHebrew2RD(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(int64(733359), hdate.ToRD(5769, hdate.Cheshvan, 15))
	assert.Equal(int64(711262), hdate.ToRD(5708, hdate.Iyyar, 6))
	assert.Equal(int64(249), hdate.ToRD(3762, hdate.Tishrei, 1))
	assert.Equal(int64(72), hdate.ToRD(3761, hdate.Nisan, 1))
	assert.Equal(int64(1), hdate.ToRD(3761, hdate.Tevet, 18))
	assert.Equal(int64(0), hdate.ToRD(3761, hdate.Tevet, 17))
	assert.Equal(int64(-1), hdate.ToRD(3761, hdate.Tevet, 16))
	assert.Equal(int64(-16), hdate.ToRD(3761, hdate.Tevet, 1))
	assert.Equal(int64(2278650), hdate.ToRD(9999, hdate.Elul, 29))
	assert.Equal(int64(731840), hdate.ToRD(5765, hdate.Tishrei, 1))
	assert.Equal(int64(731957), hdate.ToRD(5765, hdate.Shvat, 1))
	assert.Equal(int64(731987), hdate.ToRD(5765, hdate.Adar1, 1))
	assert.Equal(int64(732017), hdate.ToRD(5765, hdate.Adar2, 1))
	assert.Equal(int64(732038), hdate.ToRD(5765, hdate.Adar2, 22))
	assert.Equal(int64(732046), hdate.ToRD(5765, hdate.Nisan, 1))
	assert.Equal(int64(710347), hdate.ToRD(5706, hdate.Kislev, 7))
	assert.Equal(int64(336499), hdate.ToRD(4682, hdate.Nisan, 15))
}

func TestRD2Hebrew(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("15 Cheshvan 5769", hdate.FromRD(733359).String())
	assert.Equal("6 Iyyar 5708", hdate.FromRD(711262).String())
	assert.Equal("1 Tishrei 3762", hdate.FromRD(249).String())
	assert.Equal("1 Nisan 3761", hdate.FromRD(72).String())
	assert.Equal("8 Nisan 3761", hdate.FromRD(79).String())
	assert.Equal("18 Tevet 3761", hdate.FromRD(1).String())
	assert.Equal("17 Tevet 3761", hdate.FromRD(0).String())
	assert.Equal("16 Tevet 3761", hdate.FromRD(-1).String())
	assert.Equal("1 Tevet 3761", hdate.FromRD(-16).String())
	assert.Equal("30 Kislev 3761", hdate.FromRD(-17).String())
	assert.Equal("29 Elul 9999", hdate.FromRD(2278650).String())
	assert.Equal("22 Adar II 5765", hdate.FromRD(732038).String())
}

func TestMonthNames(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("Nisan", hdate.Nisan.String())
	assert.Equal("Tishrei", hdate.Tishrei.String())
	assert.Equal("Sh'vat", hdate.Shvat.String())
	assert.Equal("Adar I", hdate.Adar1.String())
	assert.Equal("Adar II", hdate.Adar2.String())
}

func TestMonthNames2(t *testing.T) {
	assert := assert.New(t)
	hd := hdate.New(5782, hdate.Adar1, 15)
	assert.Equal("Adar I", hd.MonthName("en"))
	assert.Equal("אַדָר א׳", hd.MonthName("he"))
	assert.Equal("אדר א׳", hd.MonthName("he-x-NoNikud"))
	hd = hdate.New(5783, hdate.Adar1, 15)
	assert.Equal("Adar", hd.MonthName("en"))
	assert.Equal("אַדָר", hd.MonthName("he"))
	assert.Equal("אדר", hd.MonthName("he-x-NoNikud"))
}

func TestAdar2ResetToAdar1(t *testing.T) {
	assert := assert.New(t)
	hd := hdate.New(5782, hdate.Adar1, 15)
	assert.Equal(hdate.Adar1, hd.Month())
	hd = hdate.New(5782, hdate.Adar2, 15)
	assert.Equal(hdate.Adar2, hd.Month())
	hd = hdate.New(5783, hdate.Adar1, 15)
	assert.Equal(hdate.Adar1, hd.Month())
	hd = hdate.New(5783, hdate.Adar2, 15)
	assert.Equal(hdate.Adar1, hd.Month())
}

func TestMonthFromName(t *testing.T) {
	toTest := []struct {
		s string
		m hdate.HMonth
	}{
		{"adar", hdate.Adar2},
		{"Adar I", hdate.Adar1},
		{"Adar II", hdate.Adar2},
		{"Adar 1", hdate.Adar1},
		{"Adar 2", hdate.Adar2},
		{"Adar1", hdate.Adar1},
		{"Adar2", hdate.Adar2},
		{"אדר א", hdate.Adar1},
		{"אדר ב", hdate.Adar2},
		{"אדר א׳", hdate.Adar1},
		{"אדר ב׳", hdate.Adar2},
		{"אדר", hdate.Adar2},
		{"Iyyar", hdate.Iyyar},
		{"Iyar", hdate.Iyyar},
		{"tammuz", hdate.Tamuz},
		{"Si", hdate.Sivan},
		{"Sh", hdate.Shvat},
	}
	for _, item := range toTest {
		month, err := hdate.MonthFromName(item.s)
		assert.Equal(t, nil, err)
		assert.Equal(t, item.m, month)
	}
}

func TestMonthFromNameEmpty(t *testing.T) {
	month, err := hdate.MonthFromName("")
	assert.Equal(t, hdate.HMonth(0), month)
	assert.Equal(t, errors.New("unable to parse month name"), err)
}

func TestMonthFromName1Char(t *testing.T) {
	month, err := hdate.MonthFromName("i")
	assert.Equal(t, hdate.Iyyar, month)
	assert.Equal(t, nil, err)
	month, err = hdate.MonthFromName("S")
	assert.Equal(t, hdate.HMonth(0), month)
	assert.Equal(t, errors.New("unable to parse month name"), err)
}

func TestDaysInHebYear(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(385, hdate.DaysInYear(5779))
	assert.Equal(355, hdate.DaysInYear(5780))
	assert.Equal(353, hdate.DaysInYear(5781))
	assert.Equal(384, hdate.DaysInYear(5782))
	assert.Equal(355, hdate.DaysInYear(5783))
	assert.Equal(383, hdate.DaysInYear(5784))
	assert.Equal(355, hdate.DaysInYear(5785))
	assert.Equal(354, hdate.DaysInYear(5786))
	assert.Equal(385, hdate.DaysInYear(5787))
	assert.Equal(355, hdate.DaysInYear(5788))
	assert.Equal(354, hdate.DaysInYear(5789))
	assert.Equal(383, hdate.DaysInYear(3762))
	assert.Equal(354, hdate.DaysInYear(3671))
	assert.Equal(353, hdate.DaysInYear(1234))
	assert.Equal(355, hdate.DaysInYear(123))
	assert.Equal(355, hdate.DaysInYear(2))
	assert.Equal(355, hdate.DaysInYear(1))

	assert.Equal(353, hdate.DaysInYear(5761))
	assert.Equal(354, hdate.DaysInYear(5762))
	assert.Equal(385, hdate.DaysInYear(5763))
	assert.Equal(355, hdate.DaysInYear(5764))
	assert.Equal(383, hdate.DaysInYear(5765))
	assert.Equal(354, hdate.DaysInYear(5766))
}

func TestDaysInMonth(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(29, hdate.DaysInMonth(hdate.Iyyar, 5780))
	assert.Equal(30, hdate.DaysInMonth(hdate.Sivan, 5780))
	assert.Equal(29, hdate.DaysInMonth(hdate.Cheshvan, 5782))
	assert.Equal(30, hdate.DaysInMonth(hdate.Cheshvan, 5783))
	assert.Equal(30, hdate.DaysInMonth(hdate.Kislev, 5783))
	assert.Equal(29, hdate.DaysInMonth(hdate.Kislev, 5784))

	assert.Equal(30, hdate.DaysInMonth(hdate.Tishrei, 5765))
	assert.Equal(29, hdate.DaysInMonth(hdate.Cheshvan, 5765))
	assert.Equal(29, hdate.DaysInMonth(hdate.Kislev, 5765))
	assert.Equal(29, hdate.DaysInMonth(hdate.Tevet, 5765))
}

func TestWeekday(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(time.Thursday, hdate.New(5769, hdate.Cheshvan, 15).Weekday())
	assert.Equal(time.Saturday, hdate.New(5708, hdate.Iyyar, 6).Weekday())
	assert.Equal(time.Sunday, hdate.New(5708, hdate.Iyyar, 7).Weekday())
	assert.Equal(time.Thursday, hdate.New(3762, hdate.Tishrei, 1).Weekday())
	assert.Equal(time.Tuesday, hdate.New(3761, hdate.Nisan, 1).Weekday())
	assert.Equal(time.Monday, hdate.New(3761, hdate.Tevet, 18).Weekday())
	assert.Equal(time.Sunday, hdate.New(3761, hdate.Tevet, 17).Weekday())
	assert.Equal(time.Saturday, hdate.New(3761, hdate.Tevet, 16).Weekday())
	assert.Equal(time.Friday, hdate.New(3761, hdate.Tevet, 1).Weekday())
	assert.Equal(time.Tuesday, hdate.New(3333, hdate.Sivan, 29).Weekday())
	assert.Equal(time.Monday, hdate.New(3333, hdate.Sivan, 28).Weekday())
	assert.Equal(time.Sunday, hdate.New(3333, hdate.Sivan, 27).Weekday())
	assert.Equal(time.Saturday, hdate.New(3333, hdate.Sivan, 26).Weekday())
	assert.Equal(time.Friday, hdate.New(3333, hdate.Sivan, 25).Weekday())
	assert.Equal(time.Thursday, hdate.New(3333, hdate.Sivan, 24).Weekday())
	assert.Equal(time.Wednesday, hdate.New(3333, hdate.Sivan, 23).Weekday())
}

func hd2iso(hd hdate.HDate) string {
	year, month, day := hd.Greg()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return d.Format(time.RFC3339)[:10]
}

func TestBefore(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	hd := hdate.FromTime(d)
	assert.Equal("2014-02-15", hd2iso(hd.Before(time.Saturday)))
}

func TestOnOrBefore(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-15", hd2iso(hdate.FromTime(d).OnOrBefore(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(hdate.FromTime(d).OnOrBefore(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(hdate.FromTime(d).OnOrBefore(time.Saturday)))

}

func TestNearest(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(hdate.FromTime(d).Nearest(time.Saturday)))
	d = time.Date(2014, time.February, 18, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-15", hd2iso(hdate.FromTime(d).Nearest(time.Saturday)))
}

func TestOnOrAfter(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(hdate.FromTime(d).OnOrAfter(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(hdate.FromTime(d).OnOrAfter(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(hdate.FromTime(d).OnOrAfter(time.Saturday)))
}

func TestAfter(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(hdate.FromTime(d).After(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(hdate.FromTime(d).After(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(hdate.FromTime(d).After(time.Saturday)))
}

func TestToString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("4 Tevet 5511", hdate.New(5511, hdate.Tevet, 4).String())
	assert.Equal("4 Elul 5782", hdate.New(5782, hdate.Elul, 4).String())
	assert.Equal("29 Adar II 5749", hdate.New(5749, hdate.Adar2, 29).String())
}

func TestGreg(t *testing.T) {
	assert := assert.New(t)
	hd := hdate.New(5765, hdate.Adar2, 22)
	gy, gm, gd := hd.Greg()
	assert.Equal(2005, gy)
	assert.Equal(time.April, gm)
	assert.Equal(2, gd)

	hd = hdate.New(5513, hdate.Tishrei, 6)
	gy, gm, gd = hd.Greg()
	assert.Equal(1752, gy)
	assert.Equal(time.September, gm)
	assert.Equal(14, gd)

	hd = hdate.New(5513, hdate.Tishrei, 5)
	gy, gm, gd = hd.Greg()
	assert.Equal(1752, gy)
	assert.Equal(time.September, gm)
	assert.Equal(2, gd)
}

func TestProlepticGreg(t *testing.T) {
	assert := assert.New(t)
	hd := hdate.New(5765, hdate.Adar2, 22)
	gy, gm, gd := hd.ProlepticGreg()
	assert.Equal(2005, gy)
	assert.Equal(time.April, gm)
	assert.Equal(2, gd)

	hd = hdate.New(5513, hdate.Tishrei, 6)
	gy, gm, gd = hd.ProlepticGreg()
	assert.Equal(1752, gy)
	assert.Equal(time.September, gm)
	assert.Equal(14, gd)

	hd = hdate.New(5513, hdate.Tishrei, 5)
	gy, gm, gd = hd.ProlepticGreg()
	assert.Equal(1752, gy)
	assert.Equal(time.September, gm)
	assert.Equal(13, gd)
}

func TestHDateJsonMarshal(t *testing.T) {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	b, err := json.Marshal(hd)
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`{"hy":5769,"hm":"Cheshvan","hd":15}`), b)
}

func TestHDateJsonUnMarshal(t *testing.T) {
	hdJson := `{"hy":5783,"hm":"Kislev","hd":18}`
	var hd hdate.HDate
	json.Unmarshal([]byte(hdJson), &hd)
	assert.Equal(t, 5783, hd.Year())
	assert.Equal(t, hdate.Kislev, hd.Month())
	assert.Equal(t, 18, hd.Day())
}

func TestFromGregorian(t *testing.T) {
	hd := hdate.FromGregorian(2008, time.November, 13)
	assert.Equal(t, "15 Cheshvan 5769", hd.String())
	hd = hdate.FromGregorian(1752, time.September, 14)
	assert.Equal(t, "6 Tishrei 5513", hd.String())
	hd = hdate.FromGregorian(1752, time.September, 2)
	assert.Equal(t, "5 Tishrei 5513", hd.String())
}

func TestFromProlepticGregorian(t *testing.T) {
	hd := hdate.FromProlepticGregorian(2008, time.November, 13)
	assert.Equal(t, "15 Cheshvan 5769", hd.String())
	hd = hdate.FromProlepticGregorian(1752, time.September, 14)
	assert.Equal(t, "6 Tishrei 5513", hd.String())
	hd = hdate.FromProlepticGregorian(1752, time.September, 13)
	assert.Equal(t, "5 Tishrei 5513", hd.String())
	hd = hdate.FromProlepticGregorian(1752, time.September, 9)
	assert.Equal(t, "1 Tishrei 5513", hd.String())
	hd = hdate.FromProlepticGregorian(1752, time.September, 2)
	assert.Equal(t, "23 Elul 5512", hd.String())
}
