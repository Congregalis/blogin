package cache_service

import (
	"strconv"
	"strings"

	"github.com/Congregalis/gin-demo/pkg/e"
)

type Tag struct {
	ID    int
	Name  string
	State int

	PageOffset int
	PageSize   int
}

func (t *Tag) GetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}

	if t.ID > 0 {
		keys = append(keys, strconv.Itoa(t.ID))
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageOffset > 0 {
		keys = append(keys, strconv.Itoa(t.PageOffset))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}
