package linkservice

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"linkShortener/internal/storage/entity"
	"linkShortener/internal/storage/repository"
	"log"
	"os"
)

var storageType = os.Getenv("STORAGE_TYPE")

var (
	ErrWrongStorageType = errors.New("wrong storage type")
)

type LinkService struct {
	repo repository.LinkRepository
}

func (s *LinkService) Init(_ context.Context) error {
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
		log.Println("[LinkService] connected to postgres")
		s.repo = repository.NewLinkGormRepo(db)
	} else if storageType == "inmemory" {
		s.repo = repository.NewLinkInmemoryRepo()
	} else {
		return ErrWrongStorageType
	}
	return nil
}

func (s *LinkService) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *LinkService) Close() error {
	return s.repo.Close()
}

func (s *LinkService) SaveLink(_ context.Context, fullLink string) (string, error) {
	// truncate protocol, www and trailing slash
	fullLink = transformLink(fullLink)

	log.Printf("[LinkService] saving new link %s", fullLink)

	link, err := entity.NewLink(fullLink)
	if err != nil {
		return "", err
	}

	tmp, err := s.repo.GetLink(link.ShortLink)
	for err == nil {
		if tmp.FullLink == link.FullLink {
			return link.ShortLink, nil
		}
		link.ShortLink = entity.CreateLink(link.ShortLink)
		tmp, err = s.repo.GetLink(link.ShortLink)
	}
	if err != repository.ErrLinkNotFound {
		return "", err
	}

	err = s.repo.SaveLink(link)
	if err != nil {
		return "", err
	}
	return link.ShortLink, nil
}

func (s *LinkService) GetLink(_ context.Context, shortLink string) (string, error) {
	log.Printf("[LinkService] getting link for %s", shortLink)
	link, err := s.repo.GetLink(shortLink)
	if err != nil {
		return "", err
	}
	return link.FullLink, nil
}

//tmp, err := r.GetLink(link.ShortLink)
//for err == nil {
//if tmp.FullLink == link.FullLink {
//return nil
//}
//link.ShortLink = entity.CreateLink(link.ShortLink)
//tmp, err = r.GetLink(link.ShortLink)
//}
