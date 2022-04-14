package tag_service

import (
	"encoding/json"

	"github.com/Congregalis/gin-demo/models"
	"github.com/Congregalis/gin-demo/pkg/gredis"
	"github.com/Congregalis/gin-demo/pkg/logging"
	"github.com/Congregalis/gin-demo/service/cache_service"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageOffset int
	PageSize   int
}

func (t *Tag) Add() error {
	// 删掉所有 tag list （因为目前无法判断所新增的 tag 属于哪个 list）
	gredis.LikeDeletes("TAG_LIST")
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Delete() error {
	// 删掉所有 tag list （因为目前无法判断所新增的 tag 属于哪个 list）
	gredis.LikeDeletes("TAG_LIST")
	return models.DeleteTag(t.ID)
}

func (t *Tag) Edit() error {
	data := map[string]interface{}{
		"modified_by": t.ModifiedBy,
		"name":        t.Name,
	}

	if t.State >= 0 {
		data["state"] = t.State
	}

	// 删掉所有 tag list （因为目前无法判断所编辑的 tag 属于哪个 list）
	gredis.LikeDeletes("TAG_LIST")
	return models.EditTag(t.ID, data)
}

func (t *Tag) GetAll() ([]*models.Tag, error) {
	var cacheTags []*models.Tag

	cache := cache_service.Tag{
		State: t.State,

		PageOffset: t.PageOffset,
		PageSize:   t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageOffset, t.PageSize, t.GetMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) ExistById() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.GetMaps())
}

func (t *Tag) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
