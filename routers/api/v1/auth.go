package v1

import (
	"net/http"

	"github.com/Congregalis/gin-demo/models"
	"github.com/Congregalis/gin-demo/pkg/app"
	"github.com/Congregalis/gin-demo/pkg/e"
	"github.com/Congregalis/gin-demo/pkg/util"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `form:"username" validate:"required,max=50"`
	Password string `form:"password" validate:"required,max=50"`
}

func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	a := auth{}

	httpCode, errCode := app.BindAndValid(c, &a)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if models.CheckAuth(a.Username, a.Password) {
		token, err := util.GenerateToken(a.Username)
		if err != nil {
			code = e.ERROR_AUTH_TOKEN
		} else {
			data["token"] = token
			code = e.SUCCESS
		}
	} else {
		code = e.ERROR_AUTH
	}

	appG.Reposonse(http.StatusOK, code, data)
}
