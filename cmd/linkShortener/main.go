package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"linkShortener/internal/appctrl"
	"linkShortener/internal/linkservice"
	"linkShortener/internal/storage/repository"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	LinkService *linkservice.LinkService
}

// GetLink Handler for GET /link?short_link=<short_link>
func (s *Server) GetLink(c echo.Context) error {
	shortLink := c.QueryParam("short_link")
	if shortLink == "" {
		// short_link parameter is required
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "short_link is required"})
	}
	fullLink, err := s.LinkService.GetLink(context.Background(), shortLink)
	if err == repository.ErrLinkNotFound {
		// requested link does not exist
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
		// full_link parameter is required
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "full_link is required"})
	}
	shortLink, err := s.LinkService.SaveLink(context.Background(), fullLink)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"short_link": shortLink})
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
	server := &Server{&linkservice.LinkService{}}

	var resources = appctrl.ServiceKeeper{
		Services: []appctrl.Service{
			server.LinkService,
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
