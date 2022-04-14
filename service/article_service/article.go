package article_service

import (
	"encoding/json"

	"github.com/Congregalis/gin-demo/models"
	"github.com/Congregalis/gin-demo/pkg/gredis"
	"github.com/Congregalis/gin-demo/pkg/logging"
	"github.com/Congregalis/gin-demo/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageOffset int
	PageSize   int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}

	// 删掉所有 article list （因为目前无法判断所编辑的 article 属于哪个 list）
	gredis.LikeDeletes("ARTICLE_LIST")
	return models.AddArticle(article)
}

func (a *Article) Delete() error {
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		gredis.Delete(key)
	}

	// 删掉所有 tag list （因为目前无法判断所新增的 tag 属于哪个 list）
	gredis.LikeDeletes("ARTICLE_LIST")

	err := models.DeleteArticle(a.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *Article) Edit() error {
	maps := make(map[string]interface{})
	maps["tag_id"] = a.TagID
	maps["modified_by"] = a.ModifiedBy
	maps["state"] = a.State

	if a.Title != "" {
		maps["title"] = a.Title
	}
	if a.Desc != "" {
		maps["desc"] = a.Desc
	}
	if a.Content != "" {
		maps["content"] = a.Content
	}
	if a.CoverImageUrl != "" {
		maps["cover_image_url"] = a.CoverImageUrl
	}

	// 删掉所编辑的 article 的缓存
	cacheArticle := cache_service.Article{ID: a.ID}
	key := cacheArticle.GetArticleKey()
	if gredis.Exists(key) {
		_, err := gredis.Delete(key)
		if err != nil {
			logging.Error(err)
		}
	}

	// 删掉所有 article list （因为目前无法判断所编辑的 article 属于哪个 list）
	gredis.LikeDeletes("ARTICLE_LIST")

	return models.EditArticle(a.ID, maps)
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Error(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageOffset: a.PageOffset,
		PageSize:   a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := models.GetArticles(a.PageOffset, a.PageSize, a.GetMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, articles, 3600)
	return articles, nil
}

func (a *Article) ExistById() (bool, error) {
	return models.ExistArticleById(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.GetMaps())
}

func (a *Article) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}
