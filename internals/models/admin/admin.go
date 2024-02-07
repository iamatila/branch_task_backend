package model

// import (
// 	"gorm.io/gorm"
// )
import "gorm.io/gorm"

type Admin struct {
	gorm.Model        // Adds some metadata fields to the table
	Adminid    string `json:"adminid" gorm:"unique;not null"`
	Firstname  string `json:"firstname" gorm:"not null"`
	Lastname   string `json:"lastname" gorm:"not null"`
	Email      string `json:"email" gorm:"unique;not null"`
	Username   string `json:"username" gorm:"unique;not null"`
	Admintype  string `json:"admintype" gorm:"not null"`
	Password   string `json:"password" gorm:"not null"`
	Phone      string `json:"phone" gorm:"unique;not null"`
	Gender     string `json:"gender"`
	DOB        string `json:"dob"`
	State      string `json:"state"`
	City       string `json:"city"`
	Country    string `json:"country"`
}

// type UserLogin struct {
// 	Email    string `json:"email"`
// 	Password string
// }
