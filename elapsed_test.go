package hdate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElapsedDays(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(int64(2110760), elapsedDays(5780))
	assert.Equal(int64(2084447), elapsedDays(5708))
	assert.Equal(int64(1373677), elapsedDays(3762))
	assert.Equal(int64(1340455), elapsedDays(3671))
	assert.Equal(int64(450344), elapsedDays(1234))
	assert.Equal(int64(44563), elapsedDays(123))
	assert.Equal(int64(356), elapsedDays(2))
	assert.Equal(int64(1), elapsedDays(1))
	assert.Equal(int64(2104174), elapsedDays(5762))
	assert.Equal(int64(2104528), elapsedDays(5763))
	assert.Equal(int64(2104913), elapsedDays(5764))
	assert.Equal(int64(2105268), elapsedDays(5765))
	assert.Equal(int64(2105651), elapsedDays(5766))
}
