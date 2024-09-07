package ratelimit

import (
	"sync"
	"time"
)

type Service struct {
	Id              string `gorm:"column:service_id"`
	LimitBpm        int    `gorm:"column:bpm"`
	LimitRpm        int    `gorm:"column:rpm"`
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

func (s *Service) Update(bodySize int) {
	defer s.m.Unlock()
	s.m.Lock()
	s.Bpm += bodySize
	s.Rpm += 1
	s.LastRequestTime = time.Now()
}

func (s *Service) init() {
	defer s.m.Unlock()
	s.m.Lock()
	s.Bpm = 0
	s.Rpm = 0
	s.LastRequestTime = time.Now()
}

type RateLimitError struct {
	message string
}

func (e *RateLimitError) Error() string {
	return e.message
}

type InvalidError struct {
	message string
}

func (e *InvalidError) Error() string {
	return e.message
}

type RateLimit struct {
	Services map[string]*Service
}

func NewRateLimit() *RateLimit {
	r := RateLimit{}
	r.Services = make(map[string]*Service)
	return &r
}

func NewRateLimitWith(services []*Service) *RateLimit {
	r := RateLimit{}
	r.Services = make(map[string]*Service)
	r.convert(services)
	return &r
}

func Of(service *Service) *RateLimit {
	r := RateLimit{}
	r.Services = make(map[string]*Service)
	r.Services[service.Id] = service

	return &r
}

func (r *RateLimit) convert(services []*Service) {
	for _, service := range services {
		r.Add(service)
	}
}

func (r *RateLimit) Add(service *Service) {
	r.Services[service.Id] = service
}

func (r *RateLimit) isExist(serviceId string, service *Service) bool {
	if s, ok := r.Services[serviceId]; ok {
		return s == service && s.Id == serviceId
	}
	return false
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
		return false, &InvalidError{
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
