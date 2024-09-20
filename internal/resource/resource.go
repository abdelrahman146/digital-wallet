package resource

import "gorm.io/gorm"

type Resources struct {
	DB     *gorm.DB
	Broker Broker
}
