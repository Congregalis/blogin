package app

import (
	"net/http"

	"github.com/Congregalis/gin-demo/pkg/e"
	"github.com/Congregalis/gin-demo/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	validate := validator.New()
	err = validate.Struct(form)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logging.Info(err.Error())
		}

		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
