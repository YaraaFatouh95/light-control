package models

import (
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Zones     []Zone    `gorm:"foreignKey:CityID" json:"zones,omitempty"`
}
type Zone struct {
	ID         uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name       string      `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	CityID     uuid.UUID   `gorm:"type:uuid" json:"city_id"`
	Luminaires []Luminaire `gorm:"foreignKey:ZoneID" json:"luminaires,omitempty"`
}
type Luminaire struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name   string    `gorm:"type:varchar(100);not null" json:"name"`
	ZoneID uuid.UUID `gorm:"type:uuid" json:"zone_id"`
	Status bool      `json:"status"`
	Dim    string    `gorm:"type:varchar(100);not null" json:"dim"`
}

type Command struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Type          string      `gorm:"type:varchar(100);not null" json:"type"`
	Payload       string      `gorm:"type:varchar(100);not null" json:"payload"`
	ScheduledTime time.Time   `json:"scheduled_time"`
	CreatedAt     time.Time   `json:"created_at"`
	Status        string      `gorm:"type:varchar(100);not null" json:"status"`
	Cities        []City      `gorm:"many2many:command_cities;" json:"cities,omitempty"`
	Zones         []Zone      `gorm:"many2many:command_zones;" json:"zones,omitempty"`
	Luminaires    []Luminaire `gorm:"many2many:command_luminaires;" json:"luminaires,omitempty"`
	EntityType    string      `gorm:"type:varchar(100);not null" json:"entity_type"`
	Entities      []uuid.UUID `gorm:"-" json:"entities,omitempty"`
}
