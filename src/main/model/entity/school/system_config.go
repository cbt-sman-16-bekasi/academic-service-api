package school

import "gorm.io/gorm"

type SystemConfig struct {
	gorm.Model
	SchoolCode   string `gorm:"not null" json:"SchoolCode"`
	Origin       string `json:"origin"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
