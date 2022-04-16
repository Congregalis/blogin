package v1

import (
	"net/http"

	"github.com/Congregalis/gin-demo/pkg/app"
	"github.com/Congregalis/gin-demo/pkg/e"
	"github.com/Congregalis/gin-demo/pkg/logging"
	"github.com/Congregalis/gin-demo/pkg/setting"
	"github.com/Congregalis/gin-demo/pkg/util"
	"github.com/Congregalis/gin-demo/service/tag_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type GetTagsForm struct {
	State int `form:"state" validate:"eq=0|eq=1"`
}

func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	form := GetTagsForm{}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		State:      form.State,
		PageOffset: util.GetPageOffset(c),
		PageSize:   setting.AppSetting.PageSize,
	}

	tags, err := tagService.GetAll()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Reposonse(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" validate:"required,max=100"`
	CreatedBy string `form:"created_by" validate:"required,max=100"`
	State     int    `form:"state" validate:"eq=0|eq=1"`
}

func AddTag(c *gin.Context) {
	appG := app.Gin{C: c}
	form := AddTagForm{}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}

	exists, err := tagService.ExistByName()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Reposonse(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	if err := tagService.Add(); err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Reposonse(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" validate:"required,min=1"`
	Name       string `form:"name" validate:"required,max=100"`
	ModifiedBy string `form:"modified_by" validate:"required,max=100"`
	State      int    `form:"state" validate:"eq=0|eq=1"`
}

func EditTag(c *gin.Context) {
	appG := app.Gin{C: c}
	form := EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistById()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Reposonse(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Reposonse(http.StatusOK, e.SUCCESS, nil)
}

type DeleteTagForm struct {
	ID int `form:"id" validate:"required,min=1"`
}

func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	form := DeleteTagForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID: form.ID,
	}

	exists, err := tagService.ExistById()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Reposonse(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Reposonse(http.StatusOK, e.SUCCESS, nil)
}
