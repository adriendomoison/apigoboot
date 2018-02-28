package dbmodel

import "github.com/jinzhu/gorm"

// Access database object
type Access struct {
	gorm.Model
	Client       string `gorm:"NOT NULL"`
	UserId       uint   `gorm:"NOT NULL"`
	Authorize    string `gorm:"NOT NULL"`
	Previous     string `gorm:"NOT NULL"`
	AccessToken  string `gorm:"NOT NULL;PRIMARY KEY"`
	RefreshToken string `gorm:"NOT NULL"`
	ExpiresIn    int32  `gorm:"NOT NULL"`
	Scope        string `gorm:"NOT NULL"`
	RedirectUri  string `gorm:"NOT NULL"`
}

// Authorize database object
type Authorize struct {
	gorm.Model
	Client      string `gorm:"NOT NULL"`
	UserId      uint   `gorm:"NOT NULL"`
	Code        string `gorm:"NOT NULL;PRIMARY KEY"`
	ExpiresIn   int32  `gorm:"NOT NULL"`
	Scope       string `gorm:"NOT NULL"`
	RedirectUri string `gorm:"NOT NULL"`
	State       string `gorm:"NOT NULL"`
}

// Client database object
type Client struct {
	UserId      uint   `gorm:"NOT NULL"`
	Id          string `gorm:"NOT NULL;PRIMARY KEY"`
	Secret      string `gorm:"NOT NULL"`
	RedirectUri string `gorm:"NOT NULL"`
}

// Refresh database object
type Refresh struct {
	gorm.Model
	Token  string `gorm:"NOT NULL;PRIMARY KEY"`
	Access string `gorm:"NOT NULL"`
}
