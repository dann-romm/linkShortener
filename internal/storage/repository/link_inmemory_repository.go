package repository

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"linkShortener/internal/storage/entity"
	"log"
	"sync"
)

var (
	linksFilename = "data/inmemory/links.gob"
)

type LinkInmemoryRepo struct {
	storage map[string]*entity.Link
	mux     sync.RWMutex
}

func NewLinkInmemoryRepo() *LinkInmemoryRepo {
	log.Println("[inmemory] loading links from file")
	r := &LinkInmemoryRepo{
		storage: make(map[string]*entity.Link),
	}
	b, err := ioutil.ReadFile(linksFilename)
	if err != nil {
		log.Println("[inmemory] file not found")
		return r
	}
	d := gob.NewDecoder(bytes.NewBuffer(b))
	err = d.Decode(&r.storage)
	if err != nil {
		log.Println("[inmemory] error loading links from file")
	}
	return r
}

// SaveLink saves a link to the repository
// link.ShortLink can be changed if it already exists due to the uniqueness of the shortLink
func (r *LinkInmemoryRepo) SaveLink(link *entity.Link) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.saveLink(link)
}

// GetLink returns a link from the repository
func (r *LinkInmemoryRepo) GetLink(shortLink string) (*entity.Link, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.getLink(shortLink)
}

// GetAllLink returns all links from the repository
func (r *LinkInmemoryRepo) GetAllLink() ([]entity.Link, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.getAllLink()
}

// UpdateLink updates a link in the repository
func (r *LinkInmemoryRepo) UpdateLink(link *entity.Link) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.updateLink(link)
}

// DeleteLink deletes a link from the repository
func (r *LinkInmemoryRepo) DeleteLink(shortLink string) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.deleteLink(shortLink)
}

func (r *LinkInmemoryRepo) saveLink(link *entity.Link) error {
	tmp, err := r.getLink(link.ShortLink)
	for err == nil {
		if tmp.FullLink == link.FullLink {
			break
		}
		link.ShortLink = entity.CreateLink(link.ShortLink)
		tmp, err = r.getLink(link.ShortLink)
	}
	r.storage[link.ShortLink] = link
	return nil
}

func (r *LinkInmemoryRepo) getLink(shortLink string) (*entity.Link, error) {
	if link, ok := r.storage[shortLink]; ok {
		return link, nil
	} else {
		return &entity.Link{}, ErrLinkNotFound
	}
}

func (r *LinkInmemoryRepo) getAllLink() ([]entity.Link, error) {
	links := make([]entity.Link, 0, len(r.storage))
	for _, link := range r.storage {
		links = append(links, *link)
	}
	return links, nil
}

func (r *LinkInmemoryRepo) updateLink(link *entity.Link) error {
	if _, ok := r.storage[link.ShortLink]; !ok {
		return ErrLinkNotFound
	}
	r.storage[link.ShortLink] = link
	return nil
}

func (r *LinkInmemoryRepo) deleteLink(shortLink string) error {
	if _, ok := r.storage[shortLink]; !ok {
		return ErrLinkNotFound
	}
	delete(r.storage, shortLink)
	return nil
}

func (r *LinkInmemoryRepo) Ping() error {
	return nil
}

func (r *LinkInmemoryRepo) Close() error {
	r.mux.Lock()
	defer r.mux.Unlock()
	log.Println("[inmemory] saving links to file")
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	if err := e.Encode(r.storage); err != nil {
		log.Println("[inmemory] error serializing map of links")
		return err
	}
	err := ioutil.WriteFile(linksFilename, b.Bytes(), 0644)
	if err != nil {
		log.Println("[inmemory] error saving links to file")
		return err
	}
	return nil
}
