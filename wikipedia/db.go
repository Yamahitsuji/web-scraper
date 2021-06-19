package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"golang.org/x/xerrors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ArticleRecord struct {
	ID        int `gorm:"primaryKey"`
	Title     string
	URL       string
	Latitude  string
	Longitude string
	Details   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *ArticleRecord) TableName() string {
	return "article"
}

func NewDB() (*gorm.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/wikipedia?charset=utf8mb4&parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	l := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	return db.Session(&gorm.Session{
		DryRun:                 false,
		PrepareStmt:            false,
		SkipDefaultTransaction: false,
		AllowGlobalUpdate:      false,
		FullSaveAssociations:   false,
		Context:                nil,
		Logger:                 l,
		NowFunc:                nil,
	}), nil
}

func CreateArticle(article *Article) error {
	bytes, err := json.Marshal(article.Details)
	if err != nil {
		return xerrors.Errorf("[ERROR] :%w", err)
	}
	r := &ArticleRecord{
		Title:     article.Title,
		URL:       article.URL,
		Latitude:  article.Latitude,
		Longitude: article.Longitude,
		Details:   string(bytes),
	}

	db, err := NewDB()
	if err != nil {
		return xerrors.Errorf("[ERROR] :%w", err)
	}
	if err := db.Create(r).Error; err != nil {
		return xerrors.Errorf("[ERROR] :%w", err)
	}
	return nil
}

func GetFirstArticleByTitle(title string) (*Article, error) {
	db, err := NewDB()
	if err != nil {
		return nil, xerrors.Errorf("[ERROR] :%w", err)
	}
	r := new(ArticleRecord)
	tx := db.First(r, "title = ?", title)
	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return nil, nil
		}
		return nil, xerrors.Errorf("[ERROR] :%w", tx.Error)
	}
	detail := make(map[string]string)
	if err := json.Unmarshal([]byte(r.Details), &detail); err != nil {
		return nil, err
	}
	article := NewArticle(r.Title, r.URL, r.Latitude, r.Longitude, detail)
	return article, nil
}

func CreateIfNotExist(article *Article) error {
	ent, err := GetFirstArticleByTitle(article.Title)
	switch {
	case err != nil:
		return err
	case ent != nil:
		return nil
	}
	if err := CreateArticle(article); err != nil {
		return err
	}
	return nil
}
