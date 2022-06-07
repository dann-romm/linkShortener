package repository

import (
	"linkShortener/internal/storage"
	"linkShortener/internal/storage/entity"
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
// link.ShortLink can be changed if it already exists due to the uniqueness of the shortLink
func (r *LinkInmemoryRepo) SaveLink(link *entity.Link) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	tmp, err := r.GetLink(link.ShortLink)
	for err == nil {
		if tmp.FullLink == link.FullLink {
			break
		}
		link.ShortLink = entity.CreateLink(link.ShortLink)
		tmp, err = r.GetLink(link.ShortLink)
	}

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
		return &entity.Link{}, storage.ErrLinkNotFound
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
		return storage.ErrLinkNotFound
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
		return storage.ErrLinkNotFound
	}
}

func (r *LinkInmemoryRepo) Ping() error {
	return nil
}

func (r *LinkInmemoryRepo) Close() error {
	// TODO: save to disk
	return nil
}
