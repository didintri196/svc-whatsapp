package usecase

import (
	"errors"
	"svc-whatsapp/domain/constants"
	"svc-whatsapp/domain/constants/messages"
	"svc-whatsapp/domain/models"
	"svc-whatsapp/domain/presenters"
	"svc-whatsapp/domain/requests"
	"svc-whatsapp/repositories"
	"svc-whatsapp/utils"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

type (
	IMDevicesUsecase interface {
		AddMDevices(req *requests.MDevicesRequest) (err error)
		FilterMDevices(ctx *gin.Context, filter *requests.FilterRequest) (presenter presenters.ArrayFilterDevicesPresenter, meta presenters.MetaResponsePresenter, err error)
		//ReadByID(id string) (presenter presenters.DetailMDevicesPresenter, err error)
		//UpdateByID(id string, req *requests.MDevicesRequest) (err error)
		//Delete(id string) (err error)
	}

	MDevicesUsecase struct {
		*Contract
	}
)

func NewMDevicesUsecase(ucContract *Contract) IMDevicesUsecase {
	return &MDevicesUsecase{ucContract}
}

func (uc MDevicesUsecase) AddMDevices(req *requests.MDevicesRequest) (err error) {
	// init
	id := uuid.New().String()
	apikey, _ := utils.NewHashHelper().HashAndSalt(id)
	model := &models.Devices{
		ID:       id,
		MUserId:  req.MUserId,
		Jid:      req.Jid,
		Server:   req.Server,
		Phone:    req.Phone,
		WorkerID: "",
		ApiKey:   apikey,
	}

	// save new candidate
	repo := repositories.NewMDevicesRepository(uc.Postgres)
	if err = repo.Create(uc.Postgres, model); err != nil {
		NewErrorLog("MDevicesUsecase.Add", "repo.Create", err.Error())
		return err
	}

	return err
}

func (uc MDevicesUsecase) FilterMDevices(ctx *gin.Context, filter *requests.FilterRequest) (presenter presenters.ArrayFilterDevicesPresenter, meta presenters.MetaResponsePresenter, err error) {
	//DWYOR
	id, ok := ctx.Get(constants.JWTPayloadUUID)
	if !ok {
		return presenter, meta, errors.New(messages.InterfaceConversionErrorMessage)
	}

	//init repo
	repoMDevices := repositories.NewMDevicesRepository(uc.Postgres)

	//set pagination
	offset, limit, page, orderBy, sort := presenters.SetPaginationParameter(filter.Page, filter.PerPage, filter.Order, filter.Sort)

	//get data filter
	modelMDevices, total, err := repoMDevices.FilterByUserID(id.(string), offset, limit, orderBy, sort, filter.Search)
	if err != nil {
		NewErrorLog("MDevicesUsecase.Filter", "repoMDevices.Filter", err.Error())
		return presenter, meta, err
	}

	//build presenter
	presenter = presenters.NewArrayFilterDevicesPresenter().Build(modelMDevices)

	//set pagination
	meta = presenters.SetPaginationResponse(page, limit, int(total))

	return presenter, meta, err
}

//func (uc MDevicesUsecase) ReadByID(id string) (presenter presenters.DetailMDevicesPresenter, err error) {
//
//	// read email
//	model := models.NewMDevices()
//	repo := repositories.NewMDevicesRepository(uc.Postgres)
//	if err = repo.Read(id, model); err != nil {
//		api.NewErrorLog("MDevicesUsecase.ReadByID", "repo.Read", err.Error())
//		return presenter, err
//	}
//
//	presenter = presenters.NewDetailMDevicesPresenter().Build(model)
//
//	return presenter, err
//}
//
//func (uc MDevicesUsecase) UpdateByID(id string, req *requests.MDevicesRequest) (err error) {
//
//	// read email
//	repo := repositories.NewMDevicesRepository(uc.Postgres)
//	modelMDevices := models.Devices{
//		Name: req.Name,
//	}
//	if err = repo.Update(uc.PostgresTX, id, modelMDevices); err != nil {
//		api.NewErrorLog("MDevicesUsecase.UpdateByID", "repo.Update", err.Error())
//		return err
//	}
//
//	return err
//}
//
//func (uc MDevicesUsecase) Delete(id string) (err error) {
//
//	// cek file is exist
//	model := models.NewMDevices()
//	repo := repositories.NewMDevicesRepository(uc.Postgres)
//	if err = repo.Read(id, model); err != nil {
//		api.NewErrorLog("MDevicesUsecase.Delete", "repo.Read", err.Error())
//		return err
//	}
//
//	//delete file
//	if err = repo.Delete(id, uc.PostgresTX); err != nil {
//		api.NewErrorLog("MDevicesUsecase.Delete", "repo.Delete", err.Error())
//		return err
//	}
//
//	return err
//}
