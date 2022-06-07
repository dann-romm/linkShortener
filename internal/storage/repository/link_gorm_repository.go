package repository

import (
	"gorm.io/gorm"
	"linkShortener/internal/storage"
	"linkShortener/internal/storage/entity"
	"log"
)

type LinkGormRepo struct {
	db *gorm.DB
}

// NewLinkGormRepo creates a new LinkGormRepo
func NewLinkGormRepo(db *gorm.DB) *LinkGormRepo {
	_ = db.AutoMigrate(&entity.Link{})
	return &LinkGormRepo{
		db: db,
	}
}

// SaveLink saves a link to the repository
// link.ShortLink can be changed if it already exists due to the uniqueness of the shortLink
func (r *LinkGormRepo) SaveLink(link *entity.Link) error {
	tmp, err := r.GetLink(link.ShortLink)
	for err == nil {
		if tmp.FullLink == link.FullLink {
			return nil
		}
		link.ShortLink = entity.CreateLink(link.ShortLink)
		tmp, err = r.GetLink(link.ShortLink)
	}
	return r.db.Create(link).Error
}

// GetLink returns a link from the repository
func (r *LinkGormRepo) GetLink(shortLink string) (*entity.Link, error) {
	var link entity.Link
	if err := r.db.Where("short_link = ?", shortLink).First(&link).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &entity.Link{}, storage.ErrLinkNotFound
		}
		return &entity.Link{}, err
	}
	return &link, nil
}

// GetAllLink returns all links from the repository
func (r *LinkGormRepo) GetAllLink() ([]entity.Link, error) {
	var links []entity.Link
	if err := r.db.Find(&links).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []entity.Link{}, storage.ErrLinkNotFound
		}
		return []entity.Link{}, err
	}
	return links, nil
}

// UpdateLink updates a link in the repository
func (r *LinkGormRepo) UpdateLink(link *entity.Link) error {
	err := r.db.Save(link).Error
	if err == gorm.ErrRecordNotFound {
		return storage.ErrLinkNotFound
	}
	return err
}

// DeleteLink deletes a link from the repository
func (r *LinkGormRepo) DeleteLink(shortLink string) error {
	err := r.db.Where("short_link = ?", shortLink).Delete(&entity.Link{}).Error
	if err == gorm.ErrRecordNotFound {
		return storage.ErrLinkNotFound
	}
	return err
}

func (r *LinkGormRepo) Ping() error {
	if db, err := r.db.DB(); err != nil {
		return err
	} else {
		return db.Ping()
	}
}

func (r *LinkGormRepo) Close() error {
	if db, err := r.db.DB(); err != nil {
		return err
	} else {
		log.Println("Closing database connection")
		return db.Close()
	}
}
