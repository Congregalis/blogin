package v1

import (
	"fmt"
	"net/http"

	"github.com/Congregalis/gin-demo/models"
	"github.com/Congregalis/gin-demo/pkg/e"
	"github.com/Congregalis/gin-demo/pkg/logging"
	"github.com/Congregalis/gin-demo/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{username, password}
	ok, _ := valid.Valid(a)

	var errMsg string
	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		if models.CheckAuth(username, password) {
			token, err := util.GenerateToken(username)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		errMsg = "("
		for _, err := range valid.Errors {
			errMsg += fmt.Sprintf("err.Key: %s, err.Message: %s ", err.Key, err.Message)
		}
		errMsg += ")"

		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code) + errMsg,
		"data": data,
	})
}
