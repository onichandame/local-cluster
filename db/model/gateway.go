package model

import (
	"github.com/onichandame/local-cluster/constants/gateway"
	"gorm.io/gorm"
)

type Gateway struct {
	gorm.Model
	Selectable
	Status        gateway.GatewayStatus
	Port          uint   `gorm:"not null"`
	ServiceName   string `gorm:"not null"`
	InterfaceName string `gorm:"not null"`
	External      bool
}
