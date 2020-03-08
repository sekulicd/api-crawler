package sqliteinfra

import (
	"api-crawler/core/collegescorecard/collegedomain"
	"github.com/jinzhu/gorm"
)

type collegeRepository struct {
	db *gorm.DB
}

func NewCollegeRepository(db *gorm.DB) collegedomain.CollegeRepository {
	return collegeRepository{
		db: db,
	}
}

func (c collegeRepository) GetAll() ([]collegedomain.School, error) {
	var schools []collegedomain.School
	err := c.db.Find(&schools).Error
	if err != nil {
		return nil, err
	}
	return schools, nil
}

func (c collegeRepository) Create(s collegedomain.School) error {
	db := c.db.Create(&s)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
