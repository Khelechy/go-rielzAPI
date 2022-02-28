package models


import (
    "errors"
    "strings"

    "github.com/jinzhu/gorm"
)

type Tenant struct {
    gorm.Model
    Email        string `gorm:"not null" json:"email"`
    FirstName    string `gorm:"size:100;not null"              json:"firstname"`
    LastName     string `gorm:"size:100;not null"              json:"lastname"`
    PhoneNumber  string `gorm:"size:100;not null"              json:"phonenumber"`
    Apartment   House   `gorm:"foreignKey:HouseId;"       json:"-"`
	HouseId  int `gorm:"not null"              json:"house_id"`
}


func (v *Tenant) Prepare() {
    v.FirstName = strings.TrimSpace(v.FirstName)
    v.LastName = strings.TrimSpace(v.LastName)
    v.PhoneNumber = strings.TrimSpace(v.PhoneNumber)
    v.Email = strings.TrimSpace(v.Email)
}

func (v *Tenant) Validate() error {
    if v.FirstName == "" {
        return errors.New("Firstname is required")
    }
    if v.LastName == "" {
        return errors.New("Lastname is required")
    }
    if v.PhoneNumber == "" {
        return errors.New("PhoneNumber is required")
    }
    if v.Email == "" {
        return errors.New("Email is required")
    }
    if v.HouseId < 0 {
        return errors.New("HouseId of house is invalid")
    }
    return nil
}

func (u *Tenant) SaveTenant(db *gorm.DB) (*Tenant, error) {
    var err error

    // Debug a single operation, show detailed log for this operation
    err = db.Debug().Create(&u).Error
    if err != nil {
        return &Tenant{}, err
    }
    return u, nil
}