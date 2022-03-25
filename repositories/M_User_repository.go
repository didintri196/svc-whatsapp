package repositories

import (
	"svc-whatsapp/domain/models"

	"gorm.io/gorm"
)

type (
	IMUserRepository interface {
		Create(tx *gorm.DB, data *models.MUser) (err error)
		ReadByEmail(email string, model *models.MUser) (err error)
	}

	MUserRepository struct {
		Postgres *gorm.DB
	}
)

func NewMUserRepository(postgres *gorm.DB) IMUserRepository {
	return &MUserRepository{Postgres: postgres}
}

func (repo MUserRepository) Create(tx *gorm.DB, data *models.MUser) (err error) {
	return tx.Omit("id", "deleted_at").Create(data).Error
}

func (repo MUserRepository) ReadByEmail(email string, model *models.MUser) (err error) {
	return repo.Postgres.Where("email = ?", email).First(model).Error
}
