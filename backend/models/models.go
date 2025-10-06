package models

import "gorm.io/gorm"

// User represents the admin user model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // Hide password from JSON responses
}

// Provider represents the telecommunication provider model
type Provider struct {
	gorm.Model
	Name    string   `json:"name" gorm:"unique"`
	Address string   `json:"address"`
	Towers  []*Tower `json:"towers,omitempty" gorm:"many2many:provider_towers;"`
}

// Tower represents the telecommunication tower model
type Tower struct {
	gorm.Model
	PhotoURL    string      `json:"photo_url"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
	Kelurahan   string      `json:"kelurahan"`
	Kecamatan   string      `json:"kecamatan"`
	Address     string      `json:"address"`
	Tinggi      float64     `json:"tinggi"`
	Tipe        string      `json:"tipe"`
	Status      string      `json:"status" gorm:"default:'active'"`
	Providers   []*Provider `json:"providers,omitempty" gorm:"many2many:provider_towers;"`
}

// BlankspotArea represents an area on the map with weak signal or no signal
type BlankspotArea struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Kelurahan   string `json:"kelurahan"` // New field
	Coordinates string `json:"coordinates" gorm:"type:text"` // Store as JSON string: [[lat, lon], [lat, lon], ...]
	Type        string `json:"type"`                        // e.g., "Blankspot", "Weak Signal"
	Color       string `json:"color"`                       // e.g., "#FF0000", "#FFFF00"
}