package ratelimit_test

import (
	"sync"
	"testing"
	"time"

	"github.com/dimsssss/ai-proxy/internal/ratelimit"
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

func TestProcessRateLimit_MultiThread(t *testing.T) {
	rl := ratelimit.NewRateLimit()
	service := ratelimit.NewService("testService", 1000, 10, time.Now(), 0, 0)
	rl.Add(service)

	startTime := time.Now()

	var wg sync.WaitGroup

	goroutines := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for time.Since(startTime) < time.Second {
				result, err := rl.ProcessRateLimit(time.Now(), service)

				if err != nil {
					if _, ok := err.(*ratelimit.RateLimitError); ok {
						continue
					}
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if result {
					service.Update(1)
				}
			}
		}()
	}

	wg.Wait()

	if service.Bpm > service.LimitBpm {
		t.Errorf("Bpm exceeded limit: %d", service.Bpm)
	}
	if service.Rpm > service.LimitRpm {
		t.Errorf("Rpm exceeded limit: %d", service.Rpm)
	}
}
