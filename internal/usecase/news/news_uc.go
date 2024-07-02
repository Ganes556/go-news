package uc_news

import (
	"context"
	"mime/multipart"
	"net/url"
	"strings"
	"sync"

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
	GetNewsByTitle(ctx context.Context, title string) (news entity.News, err error)
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
		err2 := tx.Updates(&entity.News{
			Base: entity.Base{
				ID: param.Req.ID,
			},
			UsersID:      param.UserID,
			CategoriesID: param.Req.CategoryID,
			Cover:        newCover,
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
}

func (u *ucNews) AddViewingNews(param ParamAddViewingNews) (err error) {
	err = u.db.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {

		ipread := new(entity.IpRead)
		if err2 := tx.Preload("News", func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", param.IdNews).Select("id")
		}).First(ipread, "ip = ?", param.Ip).Error; err2 != nil {
			if err2 == gorm.ErrRecordNotFound {
				if err3 := tx.Create(&entity.IpRead{
					IP: param.Ip,
					News: []entity.News{
						{
							Base: entity.Base{
								ID: param.IdNews,
							},
						}},
				}).Error; err3 != nil {
					helper.LogsError(err)
					return err3
				}
			}
			if err2 != gorm.ErrRecordNotFound {
				return err2
			}
		} else {			
			if err3 := tx.Model(&entity.IpRead{
				Base: ipread.Base,
				IP: ipread.IP,
				News: ipread.News,
			}).Association("News").Append([]entity.News{
				{
					Base: entity.Base{
						ID: param.IdNews,
					},
				}}); err3 != nil {
				helper.LogsError(err)
				return err3
			}
		}

		if len(ipread.News) == 0 {
			if err2 := tx.Model(&entity.News{}).Where("id = ?", param.IdNews).Update("count_view", gorm.Expr("count_view + 1")).Error; err2 != nil {
				helper.LogsError(err)
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
		Order("id DESC")
	if param.Category != "" {
		tx = tx.InnerJoins("Categories", u.db.Where(&entity.Categories{Name: param.Category}))
	}

	if param.Title != "" {
		tx = tx.Where("title LIKE ?", "%"+param.Title+"%").Preload("Users").Preload("Categories")
	}

	if param.Next != 0 {
		tx.Where("news.id < ?", param.Next)
	}

	err = tx.Limit(int(param.Limit)).
		Find(&news).Error
	return
}

func (u *ucNews) GetNewsByTitle(ctx context.Context, title string) (news entity.News, err error) {
	var parsedTitle string
	news = entity.News{}
	parsedTitle, err = url.QueryUnescape(title)
	if err != nil {
		return
	}
	title = strings.ToLower(parsedTitle)

	err = u.db.WithContext(ctx).Preload("Categories").First(&news, "title = ?", title).Error
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
		u.db.WithContext(ctx).Model(&entity.News{}).Select("SUM(count_view) AS count_views").Row().Scan(&totalViews);
	}()
	wg.Wait()
	return
}
