package ratelimit

import (
	"sync"
	"time"
)

type Service struct {
	Id              string
	LimitBpm        int
	LimitRpm        int
	LastRequestTime time.Time
	Bpm             int
	Rpm             int
	m               sync.RWMutex
}

func NewService(
	id string,
	LimitBpm int,
	limitRpm int,
	lastRequestTime time.Time,
	bpm int,
	rpm int) *Service {
	s := Service{
		Id:              id,
		LimitBpm:        limitRpm,
		LimitRpm:        limitRpm,
		Bpm:             bpm,
		Rpm:             rpm,
		LastRequestTime: lastRequestTime,
	}

	return &s
}

func (s *Service) update(bodySize int) {
	defer s.m.Unlock()
	s.m.Lock()
	s.Bpm += bodySize
	s.Rpm += 1
}

func (s *Service) init() {
	defer s.m.Unlock()
	s.m.Lock()
	s.Bpm = 0
	s.Rpm = 0
}

type RateLimitError struct {
	message string
}

func (e *RateLimitError) Error() string {
	return e.message
}

type RateLimit struct {
	services map[string]*Service
}

func NewRateLimit() *RateLimit {
	r := RateLimit{}
	r.services = make(map[string]*Service)
	return &r
}

func Of(service *Service) *RateLimit {
	r := RateLimit{}
	r.services = make(map[string]*Service)
	r.services[service.Id] = service

	return &r
}

func (r *RateLimit) Add(service *Service) {
	r.services[service.Id] = service
}

func (r *RateLimit) isExist(serviceId string, service *Service) bool {
	return r.services[serviceId] == service && r.services[serviceId].Id == serviceId
}

func (s *Service) exceedOneMinute(current time.Time) bool {
	return current.Sub(s.LastRequestTime).Minutes() >= 1
}

func (s *Service) exceedBpm() bool {
	return s.Bpm >= s.LimitBpm
}

func (s *Service) exceedRpm() bool {
	return s.Rpm >= s.LimitRpm
}

func (r *RateLimit) ProcessRateLimit(current time.Time, service *Service) (bool, error) {
	if !r.isExist(service.Id, service) {
		return false, &RateLimitError{
			message: "not exist service",
		}
	}

	if service.exceedOneMinute(current) {
		service.init()
		return true, nil
	}

	if service.exceedBpm() {
		return false, &RateLimitError{
			message: "bpm over",
		}
	}

	if service.exceedRpm() {
		return false, &RateLimitError{
			message: "rpm over",
		}
	}
	return true, nil
}
