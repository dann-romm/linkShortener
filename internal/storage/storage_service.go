package storage

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"linkShortener/internal/storage/entity"
	"linkShortener/internal/storage/repository"
	"os"
)

var storageType = os.Getenv("STORAGE_TYPE")

var (
	ErrWrongStorageType = errors.New("wrong storage type")
)

type StorageService struct {
	repo repository.LinkRepository
}

func (s *StorageService) Init(_ context.Context) error {
	if storageType == "postgres" {
		host := os.Getenv("POSTGRES_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("POSTGRES_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("POSTGRES_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("POSTGRES_PASSWORD")
		if password == "" {
			password = "postgres"
		}
		dbname := os.Getenv("POSTGRES_DB")
		if dbname == "" {
			dbname = "postgres"
		}
		db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)), &gorm.Config{})
		if err != nil {
			return err
		}
		s.repo = repository.NewLinkGormRepo(db)
	} else if storageType == "inmemory" {
		s.repo = repository.NewLinkInmemoryRepo()
	} else {
		return ErrWrongStorageType
	}
	return nil
}

func (s *StorageService) Ping(_ context.Context) error {
	return s.repo.Ping()
}

func (s *StorageService) Close() error {
	return s.repo.Close()
}

func (s *StorageService) SaveNewLink(_ context.Context, fullLink string) error {
	link, err := entity.NewLink(fullLink)
	if err != nil {
		return err
	}
	return s.repo.SaveLink(link)
}

func (s *StorageService) GetLink(_ context.Context, shortLink string) (string, error) {
	link, err := s.repo.GetLink(shortLink)
	return link.FullLink, err
}
