package repository

import (
	"context"
	"gorm.io/gorm"
	"linkShortener/internal/storage/entity"
	"log"
)

type LinkGormRepo struct {
	db *gorm.DB
}

// NewLinkGormRepo creates a new LinkGormRepo
func NewLinkGormRepo(db *gorm.DB) *LinkGormRepo {
	log.Println("[LinkGormRepo] migrating database")
	_ = db.AutoMigrate(&entity.Link{})
	return &LinkGormRepo{
		db: db,
	}
}

// SaveLink saves a link to the repository
func (r *LinkGormRepo) SaveLink(link *entity.Link) error {
	if _, err := r.GetLink(link.ShortLink); err == nil {
		return ErrLinkAlreadyExists
	}
	return r.db.Create(link).Error
}

// GetLink returns a link from the repository
func (r *LinkGormRepo) GetLink(shortLink string) (*entity.Link, error) {
	var link entity.Link
	if err := r.db.Find(&link, "short_link = ?", shortLink).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &entity.Link{}, ErrLinkNotFound
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
			return []entity.Link{}, ErrLinkNotFound
		}
		return []entity.Link{}, err
	}
	return links, nil
}

// UpdateLink updates a link in the repository
func (r *LinkGormRepo) UpdateLink(link *entity.Link) error {
	err := r.db.Save(link).Error
	if err == gorm.ErrRecordNotFound {
		return ErrLinkNotFound
	}
	return err
}

// DeleteLink deletes a link from the repository
func (r *LinkGormRepo) DeleteLink(shortLink string) error {
	err := r.db.Where("short_link = ?", shortLink).Delete(&entity.Link{}).Error
	if err == gorm.ErrRecordNotFound {
		return ErrLinkNotFound
	}
	return err
}

func (r *LinkGormRepo) Ping(ctx context.Context) error {
	if db, err := r.db.DB(); err != nil {
		return err
	} else {
		return db.PingContext(ctx)
	}
}

func (r *LinkGormRepo) Close() error {
	if db, err := r.db.DB(); err != nil {
		return err
	} else {
		log.Println("[LinkGormRepo] closing database connection")
		return db.Close()
	}
}
