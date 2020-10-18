package rest

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

// PostHandler represent the rest handler for post
type PostHandler struct {
	PUsecase domain.PostUsecase
}

// NewPostHandler will initialize the post resource endpoint
func NewPostHandler(app *fiber.App, p domain.PostUsecase) {
	handler := &PostHandler{
		PUsecase: p,
	}

	app.Get("/posts", handler.FetchPost)
	app.Post("/posts", handler.Store)
	app.Get("/posts/:id", handler.GetByID)
	app.Delete("/posts/:id", handler.Delete)
}

// Store will store the new Post base on given data
func (ph *PostHandler) Store(c *fiber.Ctx) (err error) {
	var post domain.Post
	err = c.BodyParser(&post)
	if err != nil {
		c.Response().SetStatusCode(http.StatusUnprocessableEntity)
		return c.JSON(err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&post); !ok {
		return c.JSON(ResponseError{Error: http.StatusBadRequest, Message: err.Error()})
	}

	ctx := c.Context()
	err = ph.PUsecase.Store(ctx, &post)
	if err != nil {
		return c.JSON(ResponseError{Error: getStatusCode(err), Message: err.Error()})
	}

	c.Response().SetStatusCode(http.StatusCreated)
	return c.JSON(post)
}

// FetchPost will fetch the Post based on given params
func (ph *PostHandler) FetchPost(c *fiber.Ctx) error {
	numS := c.Query("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.Query("cursor")
	ctx := c.Context()

	listAr, nextCursor, err := ph.PUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(ResponseError{Error: getStatusCode(err), Message: err.Error()})
	}

	c.Response().SetStatusCode(http.StatusOK)
	c.Set(`X-Cursor`, nextCursor)
	return c.JSON(listAr)
}

// GetByID will get post by given id
func (ph *PostHandler) GetByID(c *fiber.Ctx) error {
	idP, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(ResponseError{Error: http.StatusNotFound, Message: domain.ErrNotFound.Error()})
	}

	id := int64(idP)
	ctx := c.Context()

	post, err := ph.PUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(ResponseError{Error: getStatusCode(err), Message: err.Error()})
	}

	c.Response().SetStatusCode(http.StatusOK)
	return c.JSON(post)
}

// Delete will delete post by given param
func (ph *PostHandler) Delete(c *fiber.Ctx) error {
	idP, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(ResponseError{http.StatusNotFound, domain.ErrNotFound.Error()})
	}

	id := int64(idP)
	ctx := c.Context()

	err = ph.PUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(ResponseError{Error: getStatusCode(err), Message: err.Error()})
	}

	return c.SendStatus(http.StatusNoContent)
}

func isRequestValid(p *domain.Post) (bool, error) {
	validate := validator.New()
	err := validate.Struct(p)

	if err != nil {
		return false, err
	}

	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
