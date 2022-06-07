package persistence

import (
	"linkShortener/internal/entity"
	"linkShortener/internal/repository"
	"sync"
)

type LinkInmemoryRepo struct {
	storage map[string]*entity.Link
	mux     sync.Mutex
}

func NewLinkInmemoryRepo() *LinkInmemoryRepo {
	return &LinkInmemoryRepo{
		storage: make(map[string]*entity.Link),
	}
}

// SaveLink saves a link to the repository
func (r *LinkInmemoryRepo) SaveLink(link *entity.Link) error {
	// TODO: handle duplicate short link and handle link collisions
	r.mux.Lock()
	defer r.mux.Unlock()
	r.storage[link.ShortLink] = link
	return nil
}

// GetLink returns a link from the repository
func (r *LinkInmemoryRepo) GetLink(shortLink string) (*entity.Link, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if link, ok := r.storage[shortLink]; ok {
		return link, nil
	} else {
		return nil, repository.ErrLinkNotFound
	}
}

// GetAllLink returns all links from the repository
func (r *LinkInmemoryRepo) GetAllLink() ([]entity.Link, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	links := make([]entity.Link, 0, len(r.storage))
	for _, link := range r.storage {
		links = append(links, *link)
	}

	return links, nil
}

// UpdateLink updates a link in the repository
func (r *LinkInmemoryRepo) UpdateLink(link *entity.Link) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.storage[link.ShortLink]; ok {
		r.storage[link.ShortLink] = link
		return nil
	} else {
		return repository.ErrLinkNotFound
	}
}

// DeleteLink deletes a link from the repository
func (r *LinkInmemoryRepo) DeleteLink(shortLink string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.storage[shortLink]; ok {
		delete(r.storage, shortLink)
		return nil
	} else {
		return repository.ErrLinkNotFound
	}
}
