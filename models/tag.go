package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageOffset int, pageSize int, maps interface{}) ([]*Tag, error) {
	var tags []*Tag
	if err := db.Where(maps).Offset(pageOffset).Limit(pageSize).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func GetTagTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag

	if err := db.Select("id").Where("name = ?", name).First(&tag).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return tag.ID > 0, nil
}

func AddTag(name string, state int, createdBy string) error {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}

	if err := db.Create(&tag).Error; err != nil {
		return err
	}

	return nil
}

func ExistTagByID(id int) (bool, error) {
	var tag Tag

	if err := db.Select("id").Where("id = ?", id).First(&tag).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return tag.ID > 0, nil
}

func EditTag(id int, data interface{}) error {
	if err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func DeleteTag(id int) error {
	if err := db.Model(&Tag{}).Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}

	return nil
}
