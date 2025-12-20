package employee

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/zhedevops/idm/inner/common"
	"github.com/zhedevops/idm/inner/web"
	"strconv"
	"strings"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
}

// интерфейс сервиса employee.Service
type Svc interface {
	FindById(request ParamIdRequest) (Response, error)
	CreateEmployee(request CreateRequest) (int64, error)
	FindAll() ([]Response, error)
	FilterByIDs(request ParamIdsRequest) ([]Response, error)
	DeleteById(request ParamIdRequest) (int64, error)
	DeleteByIds(request ParamIdsRequest) (int64, error)
}

func NewController(server *web.Server, employeeService Svc) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
	}
}

// функция для регистрации маршрутов
func (c *Controller) RegisterRoutes() {
	// полный маршрут получится "/api/v1/employees"
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
	c.server.GroupApiV1.Get("/employees/:id", c.FindById)
	c.server.GroupApiV1.Get("/employees", c.FindAll)
	c.server.GroupApiV1.Get("/employees/list/:ids", c.FilterByIDs)
	c.server.GroupApiV1.Delete("/employees/:id", c.DeleteById)
	c.server.GroupApiV1.Post("/employees/delete", c.DeleteByIds)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees"
func (c *Controller) CreateEmployee(ctx *fiber.Ctx) error {
	// анмаршалим JSON body запроса в структуру CreateRequest
	var request CreateRequest
	if err := ctx.BodyParser(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод CreateEmployee сервиса employee.Service
	var newEmployeeId, err = c.employeeService.CreateEmployee(request)
	if err != nil {
		switch {
		// если сервис возвращает ошибку RequestValidationError или AlreadyExistsError,
		// то мы возвращаем ответ с кодом 400 (BadRequest)
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())

		// если сервис возвращает другую ошибку, то мы возвращаем ответ с кодом 500 (InternalServerError)
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}

	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	if err = common.OkResponse(ctx, newEmployeeId); err != nil {
		// функция ErrorResponse() формирует и направляет ответ в случае ошибки
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
	}

	return nil
}

func (c *Controller) FindById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	req := ParamIdRequest{Id: id}

	entity, err := c.employeeService.FindById(req)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}

	}

	if err = common.OkResponse(ctx, entity); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error get employee by id")
	}

	return nil
}

func (c *Controller) FindAll(ctx *fiber.Ctx) error {
	resp, err := c.employeeService.FindAll()
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, resp); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error get employee by id")
	}

	return nil
}

func (c *Controller) DeleteById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	req := ParamIdRequest{Id: id}

	entity, err := c.employeeService.DeleteById(req)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}

	}

	if err = common.OkResponse(ctx, entity); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error delete employee by i")
	}

	return nil
}

func (c *Controller) FilterByIDs(ctx *fiber.Ctx) error {
	idsParam := ctx.Params("ids")
	if idsParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ids is required")
	}
	rawIds := strings.Split(idsParam, ",")

	var ids []int64
	for _, v := range rawIds {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid id: "+v)
		}
		ids = append(ids, id)
	}

	req := ParamIdsRequest{Ids: ids}

	entity, err := c.employeeService.FilterByIDs(req)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}

	}

	if err = common.OkResponse(ctx, entity); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error get employee by id")
	}

	return nil
}

func (c *Controller) DeleteByIds(ctx *fiber.Ctx) error {
	idsStr := ctx.Query("ids")
	if idsStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ids is required")
	}
	rawIds := strings.Split(idsStr, ",")

	var ids []int64
	for _, v := range rawIds {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid id: "+v)
		}
		ids = append(ids, id)
	}

	req := ParamIdsRequest{Ids: ids}

	entity, err := c.employeeService.FilterByIDs(req)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}

	}

	if err = common.OkResponse(ctx, entity); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error delete employee by ids")
	}

	return nil
}
