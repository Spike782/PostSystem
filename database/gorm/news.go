package database

import (
	"PostSystem/database/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

func PostNews(uid int, title, content string) (int, error) {
	now := time.Now()
	news := &model.News{
		UserId:     uid,
		Title:      title,
		Content:    content,
		PostTime:   &now,
		DeleteTime: nil,
	}
	err := PostDB.Create(news).Error
	if err != nil {
		slog.Error("新闻发布失败!", "title", title, "error", err)
		return 0, errors.New("新闻发布失败，请稍后重试")
	}
	return news.Id, nil
}

func DeleteNews(id int) error {
	tx := PostDB.Model(&model.News{}).Where("id=? and delete_time is null", id).Update("delete_time", time.Now())
	if tx.Error != nil {
		slog.Error("delete news error", "id", id, "error", tx.Error)
		return errors.New("新闻删除失败，请稍后重试")
	} else {
		if tx.RowsAffected <= 0 {
			return fmt.Errorf("新闻id[%d]不存在", id)
		} else {
		}
		return nil
	}
}

func UpdateNews(id int, title, content string) error {
	tx := PostDB.Model(&model.News{}).Where("id=? and delete_time is null", id).Updates(map[string]any{"title": title, "content": content})
	if tx.Error != nil {
		slog.Error("update news error", "id", id, "error", tx.Error)
		return errors.New("新闻修改失败，请稍后重试")
	} else {
		if tx.RowsAffected <= 0 {
			return fmt.Errorf("新闻id[%d]不存在", id)
		} else {
			return nil
		}
	}
}

func GetNewsById(id int) *model.News {
	news := &model.News{Id: id}
	tx := PostDB.Select("*").Where("delete_time is null").First(news)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("get news error", "id", id, "error", tx.Error)
		}
		return nil
	}
	if news != nil {
		news.ViewPostTime = news.PostTime.Format("2006-01-02 15:04:05")
	}
	return news
}

func GetNewsByPage(pageNo, pageSize int) ([]*model.News, int) {
	var total int64
	err := PostDB.Model(model.News{}).Where("delete_time is null").Count(&total).Error
	if err != nil {
		slog.Error("get news error", "error", err)
		return nil, 0
	}
	var news []*model.News
	tx := PostDB.Select("*").Where("delete_time is null").Order("create_time desc").Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&news)
	if tx.Error != nil {
		slog.Error("get newsbyorder error", "pageNo", pageNo, "pageSize", pageSize, "error", tx.Error)
	}
	if len(news) > 0 {
		for _, new := range news {
			new.ViewPostTime = new.PostTime.Format("2006-01-02 15:04:05")
		}
	}
	return news, int(total)
}
