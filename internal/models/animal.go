package models

type Animal struct {
	Base
	Name    string    `gorm:"uniqueIndex;not null" json:"name"`
	Rfid    string    `gorm:"uniqueIndex" json:"rfid"`
	Devices []*Device `gorm:"many2many:devices_animals;" json:"devices"`
}
