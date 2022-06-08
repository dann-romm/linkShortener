package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"linkShortener/internal/appctrl"
	"linkShortener/internal/storage"
	"linkShortener/internal/storage/repository"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	StorageService storage.StorageService
}

// GetLink Handler for GET /link?short_link=<short_link>
func (s *Server) GetLink(c echo.Context) error {
	shortLink := c.QueryParam("short_link")
	if shortLink == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "short_link is required"})
	}
	fullLink, err := s.StorageService.GetLink(context.Background(), shortLink)
	if err == repository.ErrLinkNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"full_link": fullLink})
}

// SaveLink Handler for POST /link?full_link=<full_link>
func (s *Server) SaveLink(c echo.Context) error {
	fullLink := c.QueryParam("full_link")
	if fullLink == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "full_link is required"})
	}
	shortLink, err := s.StorageService.SaveNewLink(context.Background(), transformLink(fullLink))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"short_link": shortLink})
}

// function that truncates protocol, www and trailing slashes
func transformLink(link string) string {
	if len(link) >= 7 && link[0:7] == "http://" {
		link = link[7:]
	} else if len(link) >= 8 && link[0:8] == "https://" {
		link = link[8:]
	}
	if len(link) >= 4 && link[0:4] == "www." {
		link = link[4:]
	}
	lastSlash := len(link)
	for lastSlash > 0 && link[lastSlash-1] == '/' {
		lastSlash--
	}
	return link[0:lastSlash]
}

// TODO: add health check endpoint
func (s *Server) appStart(ctx context.Context, halt <-chan struct{}) error {
	e := echo.New()
	e.GET("/link", s.GetLink)
	e.POST("/link", s.SaveLink)

	// return 503 if halt is closed
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			select {
			case <-halt:
				log.Println("[server] halt channel closed, returning 503")
				return c.NoContent(http.StatusServiceUnavailable)
			default:
				return next(c)
			}
		}
	})

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		var port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("[server] starting server on port %s", port)
		errCh <- e.Start(fmt.Sprintf(":%s", port))
	}()

	select {
	case err := <-errCh:
		return err
	case <-halt:
		// TODO: wait current requests to finish
		time.Sleep(time.Second * 2)
		return nil
	case <-ctx.Done():
		return nil
	}
}

func main() {
	server := &Server{}

	var resources = appctrl.ServiceKeeper{
		Services: []appctrl.Service{
			&server.StorageService,
		},
		ShutdownTimeout: time.Second * 4,
		PingPeriod:      time.Second * 10,
	}

	var app = appctrl.Application{
		MainFunc:           server.appStart,
		Resources:          &resources,
		TerminationTimeout: time.Second * 7,
	}

	if err := app.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
