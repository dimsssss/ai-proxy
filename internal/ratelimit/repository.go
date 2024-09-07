package ratelimit

import (
	"gorm.io/gorm"
)

func GetServices(db *gorm.DB) []*Service {
	services := []*Service{}
	db.Table("rate_limits").Find(&services)
	return services
}
