package http

import (
	"demo-echo/model"
	models "demo-echo/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError struct
type ResponseError struct {
	Message string `json:"message"`
}

// Articlehandle struct
type Articlehandle struct {
	ArUsecase models.ArticleUsecase
}

// NewArticleHandler func API
func NewArticleHandler(e *echo.Echo, articleUsecase model.ArticleUsecase) {
	handler := &Articlehandle{
		ArUsecase: articleUsecase,
	}
	e.GET("/articles", handler.FetchArticle)
	e.GET("/articles/:id", handler.GetByID)
	e.POST("/articles", handler.Store)
	e.DELETE("/articles/:id", handler.Delete)
}

// FetchArticle method GET
func (ar *Articlehandle) FetchArticle(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := ar.ArUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

// GetByID method
func (ar *Articlehandle) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
	}
	id := int64(idP)
	ctx := c.Request().Context()
	article, err := ar.ArUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, article)
}

// Store method
func (ar *Articlehandle) Store(c echo.Context) error {
	var article model.Article
	err := c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	err = ar.ArUsecase.Store(ctx, &article)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, article)
}

// Delete method
func (ar *Articlehandle) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound)
	}
	id := int64(idP)
	ctx := c.Request().Context()

	err = ar.ArUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func isRequestValid(m *model.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
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
	case model.ErrInternalServerError:
		return http.StatusInternalServerError
	case model.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
