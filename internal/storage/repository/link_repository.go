package repository

import (
	"context"
	"errors"
	"linkShortener/internal/storage/entity"
)

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
)

type LinkRepository interface {
	SaveLink(*entity.Link) error
	GetLink(string) (*entity.Link, error)
	GetAllLink() ([]entity.Link, error)
	UpdateLink(*entity.Link) error
	DeleteLink(string) error

	Ping(ctx context.Context) error
	Close() error
}
