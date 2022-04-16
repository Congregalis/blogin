package v1

import (
	"net/http"

	"github.com/Congregalis/gin-demo/pkg/app"
	"github.com/Congregalis/gin-demo/pkg/e"
	"github.com/Congregalis/gin-demo/pkg/logging"
	"github.com/Congregalis/gin-demo/pkg/setting"
	"github.com/Congregalis/gin-demo/pkg/util"
	"github.com/Congregalis/gin-demo/service/article_service"
	"github.com/Congregalis/gin-demo/service/tag_service"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type GetArticlesForm struct {
	TagID int `form:"tag_id" validate:"required,min=1"`
	State int `form:"state" validate:"eq=0|eq=1"`
}

func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	form := GetArticlesForm{}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		TagID:      form.TagID,
		State:      form.State,
		PageOffset: util.GetPageOffset(c),
		PageSize:   setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	appG.Reposonse(http.StatusOK, e.SUCCESS, data)
}

type GetArticleForm struct {
	ID int `form:"id" validate:"required,min=1"`
}

func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	form := GetArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{ID: form.ID}
	exists, err := articleService.ExistById()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exists {
		appG.Reposonse(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusOK, e.ERROR_GET_ARTICLES_FAIL, nil)
	}

	appG.Reposonse(http.StatusOK, e.SUCCESS, article)
}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" validate:"required,min=1"`
	Title         string `form:"title" validate:"required,max=100"`
	Desc          string `form:"desc" validate:"required,max=255"`
	Content       string `form:"content" validate:"required,max=65535"`
	CreatedBy     string `form:"created_by" validate:"required,max=100"`
	CoverImageUrl string `form:"cover_image_url" validate:"max=255"`
	State         int    `form:"state" validate:"eq=0|eq=1"`
}

func AddArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	form := AddArticleForm{}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	// 检查 tag 是否存在
	tagService := tag_service.Tag{
		ID: form.TagID,
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

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		CreatedBy:     form.CreatedBy,
		State:         form.State,
	}
	if err := articleService.Add(); err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Reposonse(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID            int    `form:"id" validate:"required,min=1"`
	TagID         int    `form:"tag_id" validate:"required,min=1"`
	Title         string `form:"title" validate:"max=100"`
	Desc          string `form:"desc" validate:"max=255"`
	Content       string `form:"content" validate:"max=65535"`
	ModifiedBy    string `form:"modified_by" validate:"required,max=100"`
	CoverImageUrl string `form:"cover_image_url" validate:"max=255"`
	State         int    `form:"state" validate:"eq=0|eq=1"`
}

func EditArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	form := EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}
	exists, err := articleService.ExistById()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Reposonse(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	// 检查 tag 是否存在
	tagService := tag_service.Tag{
		ID: form.TagID,
	}
	exists, err = tagService.ExistById()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Reposonse(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Reposonse(http.StatusOK, e.SUCCESS, nil)
}

type DeleteArticleForm struct {
	ID int `form:"id" validate:"required,min=1"`
}

func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	form := DeleteArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Reposonse(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{ID: form.ID}
	exists, err := articleService.ExistById()
	if err != nil {
		appG.Reposonse(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exists {
		appG.Reposonse(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		logging.Error(err)
		appG.Reposonse(http.StatusOK, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Reposonse(http.StatusOK, e.SUCCESS, nil)
}
