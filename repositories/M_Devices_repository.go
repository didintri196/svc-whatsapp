package repositories

import (
	"svc-whatsapp/domain/models"
	"time"

	"gorm.io/gorm"
)

type (
	IMDevicesRepository interface {
		Create(tx *gorm.DB, data *models.Devices) (err error)
		Filter(offset, limit int, orderBy, sort, search string) (model []models.Devices, count int64, err error)
		FilterByUserID(id string, offset, limit int, orderBy, sort, search string) (model []models.Devices, count int64, err error)
		Read(id string, model *models.Devices) (err error)
		ReadByApiKey(apiKey string, model *models.Devices) (err error)
		ReadNotConnectWorker(model *models.Devices) (err error)
		Update(tx *gorm.DB, id string, model models.Devices) (err error)
		UpdateByJID(tx *gorm.DB, jid string, model models.Devices) (err error)
		Delete(id string, tx *gorm.DB) (err error)
	}

	MDevicesRepository struct {
		Postgres *gorm.DB
	}
)

func NewMDevicesRepository(postgres *gorm.DB) IMDevicesRepository {
	return &MDevicesRepository{Postgres: postgres}
}

func (repo MDevicesRepository) Create(tx *gorm.DB, data *models.Devices) (err error) {
	return tx.Omit("id", "deleted_at").Create(data).Error
}

func (repo MDevicesRepository) Filter(offset, limit int, orderBy, sort, search string) (model []models.Devices, count int64, err error) {
	var modelMDevices = models.NewMDevices()

	queryBuilder := repo.Postgres.Model(&modelMDevices)
	queryBuilder.Select("devices.id, devices.jid, devices.server, devices.phone, devices.worker_id, devices.api_key, devices.created_at, devices.updated_at, devices.deleted_at")
	queryBuilder.Where("devices.deleted_at IS NULL")

	if search != "" {
		queryBuilder.Where("devices.phone LIKE '%" + search + "%'")
	}
	countQuery := queryBuilder

	queryBuilder.Order(orderBy + ` ` + sort)
	queryBuilder.Offset(offset).Limit(limit)
	err = queryBuilder.Scan(&model).Error

	if err != nil {
		return model, count, err
	}

	// hitung total data
	countQuery.Offset(-1).Limit(-1).Count(&count)
	return model, count, err
}

func (repo MDevicesRepository) FilterByUserID(id string, offset, limit int, orderBy, sort, search string) (model []models.Devices, count int64, err error) {
	var modelMDevices = models.NewMDevices()

	queryBuilder := repo.Postgres.Model(&modelMDevices)
	queryBuilder.Select("devices.id, devices.jid, devices.server, devices.phone, devices.worker_id, devices.api_key, devices.created_at, devices.updated_at, devices.deleted_at")
	queryBuilder.Where("devices.deleted_at IS NULL")
	queryBuilder.Where("devices.m_user_id = ?", id)

	if search != "" {
		queryBuilder.Where("devices.phone LIKE '%" + search + "%'")
	}

	countQuery := queryBuilder

	queryBuilder.Order(orderBy + ` ` + sort)
	queryBuilder.Offset(offset).Limit(limit)
	err = queryBuilder.Scan(&model).Error

	if err != nil {
		return model, count, err
	}

	// hitung total data
	countQuery.Offset(-1).Limit(-1).Count(&count)
	return model, count, err
}

func (repo MDevicesRepository) Read(id string, model *models.Devices) (err error) {
	return repo.Postgres.Where("id = ?", id).First(model).Error
}

func (repo MDevicesRepository) ReadByApiKey(apiKey string, model *models.Devices) (err error) {
	return repo.Postgres.Where("api_key = ?", apiKey).First(model).Error
}

func (repo MDevicesRepository) ReadNotConnectWorker(model *models.Devices) (err error) {
	return repo.Postgres.Where("worker_id = ?", "").Where("devices.deleted_at IS NULL").First(model).Error
}

func (repo MDevicesRepository) Update(tx *gorm.DB, id string, model models.Devices) (err error) {
	return tx.Model(models.NewMDevices()).Where("id = ?", id).Updates(&model).Error
}

func (repo MDevicesRepository) UpdateByJID(tx *gorm.DB, jid string, model models.Devices) (err error) {
	return tx.Model(models.NewMDevices()).Where("jid = ?", jid).Updates(&model).Error
}

func (repo MDevicesRepository) Delete(id string, tx *gorm.DB) (err error) {
	return tx.Model(models.Devices{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
