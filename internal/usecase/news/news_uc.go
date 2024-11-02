package uc_news

import (
	"context"
	"mime/multipart"
	"sync"

	"github.com/gosimple/slug"
	"github.com/news/helper"
	req_dto_news "github.com/news/internal/dto/request/news"
	dto_response "github.com/news/internal/dto/response"
	"github.com/news/internal/entity"
	"github.com/news/pkg"
	"gorm.io/gorm"
)

type UcNews interface {
	Create(param ParamCreate) (err error)
	Update(param ParamUpdate) (err error)
	Delete(param ParamDelete) (err error)
	AddViewingNews(param ParamAddViewingNews) (err error)
	GetNews(param ParamGetNews) (news []entity.News, err error)
	GetNewsByFilter(param ParamGetNewsByFilter) (news []entity.News, err error)
	GetNewsById(ctx context.Context, id uint) (news entity.News, err error)
	GetNewsBySlug(ctx context.Context, slug string) (news entity.News, err error)
	GetNewsMostViewed(ctx context.Context) (news []entity.News, err error)
	GetTotalPostAndViews(ctx context.Context) (totalPost, totalViews int64)
}

type ucNews struct {
	Gc pkg.Gcloud
	db *gorm.DB
}

func NewNewsUc(Gc pkg.Gcloud, DB *gorm.DB) UcNews {
	return &ucNews{Gc, DB}
}

type ParamDelete struct {
	Ctx context.Context
	Req req_dto_news.DeleteNews
}

func (u *ucNews) Delete(param ParamDelete) (err error) {
	err = u.db.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {
		oldNews := new(entity.News)
		if err := tx.First(oldNews, "id = ?", param.Req.ID).Error; err != nil {
			return err
		}
		if err := u.Gc.DeleteInStorage(param.Ctx, []string{oldNews.Cover}); err != nil {
			return err
		}
		if err := tx.Delete(oldNews, "id = ?", param.Req.ID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helper.LogsError(err)
	}
	return
}

type ParamCreate struct {
	Ctx    context.Context
	Req    req_dto_news.CreateNews
	UserID uint
}

func (u *ucNews) Create(param ParamCreate) (err error) {
	url, err := u.Gc.Upload2Storage(param.Ctx, "cover", []*multipart.FileHeader{param.Req.Cover})

	if err != nil {
		helper.LogsError(err)
		return err
	}

	err = u.db.WithContext(param.Ctx).Create(&entity.News{
		UsersID:      param.UserID,
		CategoriesID: param.Req.CategoryID,
		Slug:         slug.Make(param.Req.Title),
		Title:        param.Req.Title,
		Cover:        url[0],
		Content:      param.Req.Contents,
	}).Error

	if err != nil {
		helper.LogsError(err)
		return err
	}

	return nil
}

type ParamUpdate struct {
	Ctx    context.Context
	Req    req_dto_news.UpdateNews
	UserID uint
}

func (u *ucNews) Update(param ParamUpdate) (err error) {
	err = u.db.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {
		var oldData entity.News

		if err2 := tx.First(&oldData, "id = ?", param.Req.ID).Error; err2 != nil {
			if err2 == gorm.ErrRecordNotFound {
				return new(dto_response.Response).Err404("news")
			}
			helper.LogsError(err2)
			return err2
		}
		// delete in gc
		var newCover string
		if param.Req.Cover != nil {
			err2 := u.Gc.DeleteInStorage(param.Ctx, []string{oldData.Cover})
			if err2 != nil {
				helper.LogsError(err2)
				return err2
			}
			covers, err2 := u.Gc.Upload2Storage(param.Ctx, "cover", []*multipart.FileHeader{param.Req.Cover})
			if err2 != nil {
				helper.LogsError(err2)
				return err2
			}
			newCover = covers[0]
		}
		var slugText string
		if param.Req.Title != "" {
			slugText = slug.Make(param.Req.Title)
		}
		err2 := tx.Updates(&entity.News{
			Base: entity.Base{
				ID: param.Req.ID,
			},
			UsersID:      param.UserID,
			CategoriesID: param.Req.CategoryID,
			Cover:        newCover,
			Slug:         slugText,
			Title:        param.Req.Title,
			Content:      param.Req.Contents,
		}).Error

		if err2 != nil {
			helper.LogsError(err2)
			return err2
		}

		return nil
	})
	return
}

type ParamAddViewingNews struct {
	Ctx    context.Context
	Ip     string
	IdNews uint
	IdCategory uint
}

func (u *ucNews) AddViewingNews(param ParamAddViewingNews) (err error) {
	err = u.db.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {

		newIpread := entity.IpRead{
			IP: param.Ip,
			IpReadable: []entity.IpReadable{
				{
					OwnerID: param.IdCategory,
					OwnerType: "categories",
				},
				{
					OwnerID: param.IdNews,
					OwnerType: "news",
				},
			},
		}

		if err2 := tx.First(&entity.IpRead{}, "ip = ?", param.Ip).Error; err2 != nil {
			if err2 == gorm.ErrRecordNotFound {
				if err3 := tx.Create(&newIpread).Error; err3 != nil {
					helper.LogsError(err)
					return err3
				}
			}
			if err2 != gorm.ErrRecordNotFound {
				return err2
			}
		}
		
		return nil
	})
	return
}

type ParamGetNews struct {
	Ctx   context.Context
	Next  uint
	Limit uint
}

func (u *ucNews) GetNews(param ParamGetNews) (news []entity.News, err error) {
	news = []entity.News{}
	if param.Limit <= 0 {
		param.Limit = 10
	}
	tx := u.db.WithContext(param.Ctx).Omit("content").Order("id DESC").Preload("Users").Preload("Categories")

	if param.Next != 0 {
		tx.Where("news.id < ?", param.Next)
	}

	err = tx.Limit(int(param.Limit)).
		Find(&news).Error
	return
}
func (u *ucNews) GetNewsById(ctx context.Context, id uint) (news entity.News, err error) {
	news = entity.News{}
	err = u.db.WithContext(ctx).First(&news, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = new(dto_response.Response).Err404("news")
		}
		helper.LogsError(err)
	}
	return
}

type ParamGetNewsByFilter struct {
	Ctx      context.Context
	Category string
	Title    string
	Next     uint
	Limit    uint
}

func (u *ucNews) GetNewsByFilter(param ParamGetNewsByFilter) (news []entity.News, err error) {
	news = []entity.News{}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	tx := u.db.WithContext(param.Ctx).
		Omit("content").
		Order("news.id DESC")
	if param.Category != "" {
		tx = tx.InnerJoins("Categories", u.db.Where(&entity.Categories{Name: param.Category}))
	}

	if param.Title != "" {
		tx = tx.Where("title LIKE ?", "%"+param.Title+"%").Preload("Categories")
	}

	if param.Next != 0 {
		tx.Where("news.id < ?", param.Next)
	}

	err = tx.Limit(int(param.Limit)).Preload("Users").
		Find(&news).Error
	return
}

func (u *ucNews) GetNewsBySlug(ctx context.Context, slug string) (news entity.News, err error) {

	news = entity.News{}
	err = u.db.WithContext(ctx).Preload("Categories").First(&news, "slug = ?", slug).Error
	if err == gorm.ErrRecordNotFound {
		err = new(dto_response.Response).Err404("news")
	}

	return
}

func (u *ucNews) GetNewsMostViewed(ctx context.Context) (news []entity.News, err error) {
	err = u.db.WithContext(ctx).Order("count_view DESC").Preload("Categories").Limit(10).Find(&news).Error
	if err != nil {
		helper.LogsError(err)
	}
	return
}

func (u *ucNews) GetTotalPostAndViews(ctx context.Context) (totalPost, totalViews int64) {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		u.db.WithContext(ctx).Find(&entity.News{}).Select("id").Count(&totalPost)
	}()
	go func() {
		defer wg.Done()
		u.db.WithContext(ctx).Model(&entity.News{}).Select("SUM(count_view) AS count_views").Row().Scan(&totalViews)
	}()
	wg.Wait()
	return
}
