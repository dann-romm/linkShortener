package persistence

import (
	"gorm.io/gorm"
	"linkShortener/internal/entity"
)

type LinkGormRepo struct {
	db *gorm.DB
}

// NewLinkGormRepo creates a new LinkGormRepo
func NewLinkGormRepo(db *gorm.DB) *LinkGormRepo {
	return &LinkGormRepo{
		db: db,
	}
}

// SaveLink saves a link to the repository
func (r *LinkGormRepo) SaveLink(link *entity.Link) error {
	// TODO: handle duplicate short link and handle link collisions
	return r.db.Create(link).Error
}

// GetLink returns a link from the repository
func (r *LinkGormRepo) GetLink(shortLink string) (*entity.Link, error) {
	var link entity.Link
	if err := r.db.Where("short_link = ?", shortLink).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

// GetAllLink returns all links from the repository
func (r *LinkGormRepo) GetAllLink() ([]entity.Link, error) {
	var links []entity.Link
	if err := r.db.Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

// UpdateLink updates a link in the repository
func (r *LinkGormRepo) UpdateLink(link *entity.Link) error {
	return r.db.Save(link).Error
}

// DeleteLink deletes a link from the repository
func (r *LinkGormRepo) DeleteLink(shortLink string) error {
	return r.db.Where("short_link = ?", shortLink).Delete(&entity.Link{}).Error
}
