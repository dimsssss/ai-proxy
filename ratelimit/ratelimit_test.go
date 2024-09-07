package ratelimit_test

import (
	"testing"
	"time"

	"github.com/dimsssss/ai-proxy/ratelimit"
	"github.com/stretchr/testify/assert"
)

func TestProcessRateLimit_Not_Exist_Service(t *testing.T) {
	now := time.Now()
	service := ratelimit.NewService(
		"1",
		10,
		10,
		now.Add(time.Second),
		10,
		10,
	)

	r := ratelimit.Of(service)

	notRegisteredService := ratelimit.NewService(
		"2",
		10,
		10,
		now.Add(time.Second),
		10,
		10,
	)

	success, err := r.ProcessRateLimit(now, notRegisteredService)
	assert.Equal(t, success, false, "if a service not registered, should be return false")
	assert.Equal(t, err.Error(), "not exist service", "if a service not registered, should be return error")
}

func TestProcessRateLimit_ExceedTime(t *testing.T) {
	now := time.Now()
	service := ratelimit.NewService(
		"1",
		10,
		10,
		now.Add(-2*time.Minute),
		10,
		10,
	)

	r := ratelimit.Of(service)

	success, _ := r.ProcessRateLimit(now, service)
	assert.Equal(t, 0, service.Bpm, "After 1 minute, bpm should be zero")
	assert.Equal(t, 0, service.Rpm, "After 1 minute, rpm should be zero")
	assert.Equal(t, success, true, "After 1 minute, return true")
}

func TestProcessRateLimit_ExceedBpm(t *testing.T) {
	now := time.Now()
	service := ratelimit.NewService(
		"1",
		10,
		10,
		now.Add(1*time.Second),
		10,
		10,
	)

	r := ratelimit.Of(service)

	success, err := r.ProcessRateLimit(now, service)

	assert.Equal(t, success, false, "Within 1 minute, if bpm is greater than the limit, should return fail")
	assert.Equal(t, err.Error(), "bpm over", "Within 1 minute, if bpm is greater than the limit, should return error")
}

func TestProcessRateLimit_ExceedRpm(t *testing.T) {
	now := time.Now()
	service := ratelimit.NewService(
		"1",
		10,
		10,
		now.Add(1*time.Second),
		9,
		11,
	)

	r := ratelimit.Of(service)

	success, err := r.ProcessRateLimit(now, service)

	assert.Equal(t, success, false, "Within 1 minute, if rpm is greater than the limit, should return fail")
	assert.Equal(t, err.Error(), "rpm over", "Within 1 minute, if rpm is greater than the limit, should return error")
}

func TestProcessRateLimit_Success(t *testing.T) {
	now := time.Now()
	service := ratelimit.NewService(
		"1",
		10,
		10,
		now.Add(1*time.Second),
		9,
		9,
	)

	r := ratelimit.Of(service)

	success, err := r.ProcessRateLimit(now, service)

	assert.Equal(t, success, true, "Within 1 minute, if rpm and bpm is less than the limit, should return true")
	assert.Equal(t, err, nil, "Within 1 minute, if rpm and bpm is less than the limit, should return nil")
}
