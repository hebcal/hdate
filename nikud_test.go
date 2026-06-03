package hdate_test

import (
	"fmt"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/stretchr/testify/assert"
)

func TestHebrewStripNikkud(t *testing.T) {
	assert := assert.New(t)
	src := "חֲנוּכָּה יוֹם ד׳ (בְּשַׁבָּת)"
	dest := "חנוכה יום ד׳ (בשבת)"
	assert.Equal(dest, hdate.HebrewStripNikkud(src))
	assert.Equal("א־ב", hdate.HebrewStripNikkud("א־ב"))
	assert.Equal("אב", hdate.HebrewStripNikkud("אֿב"))
}

func ExampleHebrewStripNikkud() {
	src := "וְהָאָ֗רֶץ הָיְתָ֥ה תֹ֙הוּ֙ וָבֹ֔הוּ"
	dest := hdate.HebrewStripNikkud(src)
	fmt.Println(dest)
	// Output:
	// והארץ היתה תהו ובהו
}
