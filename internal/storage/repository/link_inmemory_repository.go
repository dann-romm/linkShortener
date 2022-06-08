package repository

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"linkShortener/internal/storage/entity"
	"sync"
)

var (
	linksFilename = "data/inmemory/links.gob"
)

type LinkInmemoryRepo struct {
	storage map[string]*entity.Link
	mux     sync.Mutex
}

func NewLinkInmemoryRepo() *LinkInmemoryRepo {
	r := &LinkInmemoryRepo{
		storage: make(map[string]*entity.Link),
	}
	b, err := ioutil.ReadFile(linksFilename)
	if err != nil {
		return r
	}
	d := gob.NewDecoder(bytes.NewBuffer(b))
	_ = d.Decode(&r.storage)
	return r
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
		return &entity.Link{}, ErrLinkNotFound
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
		return ErrLinkNotFound
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
		return ErrLinkNotFound
	}
}

func (r *LinkInmemoryRepo) Ping() error {
	return nil
}

func (r *LinkInmemoryRepo) Close() error {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	if err := e.Encode(r.storage); err != nil {
		return err
	}
	err := ioutil.WriteFile(linksFilename, b.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}
