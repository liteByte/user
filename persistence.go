package main

import (
	"github.com/jinzhu/gorm"
	"os"
)

type Persistence interface {
	Create(c *Model) error
	Migrate(c *Model) error
	UpdateFields(c *Model, updates map[string]interface{}) error
	Delete(c *Model) error
	Find(filter []map[string]string, order map[string]string, offset, limit int) ([]Model, error)
	UpdateMany(updates map[string]interface{}, filter []map[string]string) error
	DeleteMany(filter []map[string]string) error
}

type PersistenceHandler struct {
	DB *gorm.DB
}

func (h *PersistenceHandler) Create(c *Model) error {
	v, _ := os.LookupEnv("ENV")
	if v == "test" {
		return nil
	}
	if err := h.DB.Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (h *PersistenceHandler) Migrate(c *Model) error {
	v, _ := os.LookupEnv("ENV")
	if v == "test" {
		return nil
	}
	if err := h.DB.AutoMigrate(c).Error; err != nil {
		return err
	}
	return nil
}

func (h *PersistenceHandler) UpdateFields(c *Model, updates map[string]interface{}) error {
	v, _ := os.LookupEnv("ENV")
	if v == "test" {
		return nil
	}
	if err := h.DB.Model(c).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (h *PersistenceHandler) Delete(c *Model) error {
	v, _ := os.LookupEnv("ENV")
	if v == "test" {
		return nil
	}
	if err := h.DB.Delete(c).Error; err != nil {
		return err
	}
	return nil
}

func (h *PersistenceHandler) Find(filter []map[string]string, order map[string]string, offset, limit int) ([]Model, error) {

	var models []Model

	v, _ := os.LookupEnv("ENV")
	if v == "test" {
		return models, nil
	}

	db := h.DB

	db = h.applyFilter(db, filter)
	db = h.applyOrder(db, order)
	db = h.applyPagination(db, offset, limit)

	if err := db.Find(&models).Error; err != nil {
		return models, err
	}

	return models, nil
}

func (h *PersistenceHandler) UpdateMany(updates map[string]interface{}, filter []map[string]string) error {

	v, _ := os.LookupEnv("ENV")
	if v == "test" {
		return nil
	}

	db := h.DB

	db = h.applyFilter(db, filter)

	var model Model
	if err := db.Model(&model).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (h *PersistenceHandler) DeleteMany(filter []map[string]string) error {

	v, _ := os.LookupEnv("ENV")
	if v == "test" {
		return nil
	}

	db := h.DB

	db = h.applyFilter(db, filter)

	var model Model
	if err := db.Delete(&model).Error; err != nil {
		return err
	}

	return nil
}

func (h *PersistenceHandler) applyFilter(db *gorm.DB, filter []map[string]string) *gorm.DB {
	for _, q := range filter {
		for k, v := range q {
			db = db.Where(k, v)
		}
	}
	return db
}

func (h *PersistenceHandler) applyOrder(db *gorm.DB, order map[string]string) *gorm.DB {
	for k, v := range order {
		db = db.Order(k + " " + v)
	}
	return db
}

func (h *PersistenceHandler) applyPagination(db *gorm.DB, offset, limit int) *gorm.DB {
	return db.Offset(offset).Limit(limit)
}
