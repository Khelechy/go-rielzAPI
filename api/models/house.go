//house.go
package models

import (
    "errors"
    "github.com/jinzhu/gorm"
    "strings"
)

type House struct {
    gorm.Model
    HouseType        string `gorm:"size:100;not null;" json:"house_type"`
    State string `gorm:"not null"                 json:"state"`
	Description 	string `gorm:"not null"         json:"description"`
    Location    string `gorm:"size:100;not null"        json:"location"`
    Rooms    int    `gorm:"not null"                 json:"rooms"`
	AvailableRooms    int    `gorm:"not null"                 json:"available_rooms"`
	BathRooms    int    `gorm:"not null"                 json:"bathrooms"`
    Price    int `gorm:"not null"        json:"price"`
	LongLat    string    `gorm:"not null"                 json:"long_lat"`
    CreatedBy   User   `gorm:"foreignKey:UserID;"       json:"-"`
    UserID      uint   `gorm:"not null"                 json:"user_id"`
}

func (v *House) Prepare() {
    v.HouseType = strings.TrimSpace(v.HouseType)
    v.Description = strings.TrimSpace(v.Description)
    v.Location = strings.TrimSpace(v.Location)
    v.CreatedBy = User{}
}

func (v *House) Validate() error {
    if v.HouseType == "" {
        return errors.New("HouseType is required")
    }
    if v.Description == "" {
        return errors.New("Description about house is required")
    }
    if v.Location == "" {
        return errors.New("Location of house is required")
    }
    if v.Price < 0 {
        return errors.New("Price of house is invalid")
    }
    if v.Rooms < 0 {
        return errors.New("Number of Rooms of house is invalid")
    }
	if v.BathRooms < 0 {
        return errors.New("Number of BathRooms of house is invalid")
    }
    return nil
}

func (v *House) Save(db *gorm.DB) (*House, error) {
    var err error

    // Debug a single operation, show detailed log for this operation
    err = db.Debug().Create(&v).Error
    if err != nil {
        return &House{}, err
    }
    return v, nil
}

func GetHouses(db *gorm.DB) (*[]House, error) {
    houses := []House{}
    if err := db.Debug().Table("houses").Find(&houses).Error; err != nil {
        return &[]House{}, err
    }
    return &houses, nil
}

func GetHousesByLandLord(id uint, db *gorm.DB) (*[]House, error){
	houses := []House{}
	if err := db.Debug().Table("houses").Where("user_id = ?", id).Find(&houses).Error; err != nil {
		return &[]House{}, err
	}

	return &houses, nil
}

func GetHousesByState(name string, db *gorm.DB) (*[]House, error){
	houses := []House{}
	if err := db.Debug().Table("houses").Where("state = ?", name).Find(&houses).Error; err != nil {
		return &[]House{}, err
	}

	return &houses, nil
}

func GetHouseById(id int, db *gorm.DB) (*House, error) {
    house := &House{}
    if err := db.Debug().Table("houses").Where("id = ?", id).First(house).Error; err != nil {
        return nil, err
    }
    return house, nil
}

func (v *House) UpdateHouse(id int, db *gorm.DB) (*House, error) {
    if err := db.Debug().Table("houses").Where("id = ?", id).Updates(House{
        HouseType:        v.HouseType,
        Description: v.Description,
        Location:    v.Location,
        Rooms:    v.Rooms,
        AvailableRooms: v.AvailableRooms,
		BathRooms:	v.BathRooms,
        Price:    v.Price}).Error; err != nil {
        return &House{}, err
    }
    return v, nil
}

func DeleteHouse(id int, db *gorm.DB) error {
    if err := db.Debug().Table("houses").Where("id = ?", id).Delete(&House{}).Error; err != nil {
        return err
    }
    return nil
}