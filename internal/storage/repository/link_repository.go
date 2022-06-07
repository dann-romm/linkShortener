package repository

import (
	"errors"
	"linkShortener/internal/storage/entity"
)

type LinkRepository interface {
	SaveLink(*entity.Link) error
	GetLink(string) (*entity.Link, error)
	GetAllLink() ([]entity.Link, error)
	UpdateLink(*entity.Link) error
	DeleteLink(string) error

	Ping() error
	Close() error
}

var (
	ErrLinkNotFound = errors.New("link not found")
)
