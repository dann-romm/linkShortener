package linkservice

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"linkShortener/internal/storage/repository"
	"testing"
)

func openDB() *gorm.DB {
	host := "localhost"
	port := "5432"
	user := "link_shortener"
	password := "h3tGFWqJRNEaTyycMITs3"
	dbname := "link_shortener"

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestLinkService_SaveLink(t *testing.T) {
	_ = LinkService{
		repo: repository.NewLinkGormRepo(openDB()),
	}

}
