package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model

	// gorm:"index" 申明该字段为索引
	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
}

func ExistArticleById(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return article.ID > 0, nil
}

func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func GetArticles(pageOffset int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	// Preload就是一个预加载器，它会执行两条 SQL
	// 分别是 SELECT * FROM blog_articles;
	// 和 SELECT * FROM blog_tag WHERE id IN (1,2,3,4);
	if err := db.Preload("Tag").Where(maps).Offset(pageOffset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	// Article有一个结构体成员是 TagID，就是外键。
	// gorm会通过 类名+ID 的方式去找到这两个类之间的关联关系

	// Article有一个结构体成员是Tag，就是我们嵌套在Article里的Tag结构体，
	// 我们可以通过Related进行关联查询
	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		CoverImageUrl: data["cover_image_url"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
	}

	if err := db.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}

	return nil
}
