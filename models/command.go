package models

import (
	"time"
)

type Command struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AppID     uint      `json:"app_id"`
	Value     string    `json:"value" gorm:"unique"`
}

type CommandData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Command     string `json:"command"`
}

type CommandRepository struct{}

func (r *CommandRepository) GetCommandData(cmd string) (*CommandData, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	var data CommandData
	err = db.Select("a.name, a.description, a.subject, c.value as command").
		Table("commands c").
		Joins("LEFT JOIN apps a ON a.id = c.app_id").
		Where("c.value = ?", cmd).Scan(&data).Error

	return &data, err
}
