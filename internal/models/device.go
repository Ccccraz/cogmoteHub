package models

type Device struct {
	Base
	Type     string    `gorm:"not null" json:"type"`
	Nickname string    `json:"nickname"`
	Hostname string    `gorm:"uniqueIndex;not null" json:"hostname"`
	Os       string    `gorm:"not null" json:"os"`
	Animals  []*Animal `gorm:"many2many:devices_animals;" json:"animals"`
}
