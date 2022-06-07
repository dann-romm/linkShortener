package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"linkShortener/internal/appctl"
	"linkShortener/internal/storage"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	StorageService storage.StorageService
}

func (s *Server) GetLink(c echo.Context) error {
	shortLink := c.QueryParam("short_link")
	if shortLink == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "short_link is required"})
	}
	fullLink, err := s.StorageService.GetLink(context.Background(), shortLink)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"full_link": fullLink})
}

func (s *Server) SaveLink(c echo.Context) error {
	fullLink := c.FormValue("full_link")
	if fullLink == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "full_link is required"})
	}
	shortLink, err := s.StorageService.SaveNewLink(context.Background(), fullLink)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"short_link": shortLink})
}

func (s *Server) appStart(ctx context.Context, halt <-chan struct{}) error {
	e := echo.New()
	e.GET("/link", s.GetLink)
	e.POST("/link", s.SaveLink)
	//e.HTTPErrorHandler = func(err error, c echo.Context) {
	//
	//
	//	if err == echo.ErrNotFound {
	//		_ = c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	//	} else {
	//		_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	//	}
	//}

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return e.Start(fmt.Sprintf(":%s", port))
}

func main() {
	server := &Server{}

	var resources = appctl.ServiceKeeper{
		Services: []appctl.Service{
			&server.StorageService,
		},
		ShutdownTimeout: time.Second * 10,
		PingPeriod:      time.Second * 5,
	}

	var app = appctl.Application{
		MainFunc:           server.appStart,
		Resources:          &resources,
		TerminationTimeout: time.Second * 10,
	}

	if err := app.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
