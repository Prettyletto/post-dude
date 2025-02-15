package service

import (
	"fmt"

	"github.com/Prettyletto/post-dude/cmd/internal/model"
	"github.com/Prettyletto/post-dude/cmd/internal/repository"
)

type CollectionService interface {
	CreateCollection(collection *model.Collection) error
	RetrieveAllCollections() ([]model.Collection, error)
	RetrieveCollectionById(id int) (*model.Collection, error)
	UpdateCollection(id int, collection *model.Collection) (*model.Collection, error)
	DeleteCollection(id int) error
}

type collectionService struct {
	repo repository.CollectionRepository
}

func NewCollectionService(repo repository.CollectionRepository) CollectionService {
	return &collectionService{repo: repo}
}

func (s *collectionService) CreateCollection(model *model.Collection) error {
	if model.Name == "" {
		return fmt.Errorf("Name of collection cannot be empty")
	}
	return s.repo.SaveCollection(model)
}

func (s *collectionService) RetrieveAllCollections() ([]model.Collection, error) {
	return s.repo.FindAllCollections()
}

func (s *collectionService) RetrieveCollectionById(id int) (*model.Collection, error) {
	return s.repo.FindCollectionById(id)
}

func (s *collectionService) UpdateCollection(id int, collection *model.Collection) (*model.Collection, error) {
	if collection.Name == "" {
		return nil, fmt.Errorf("Name of collection cannot be empty")
	}
	err := s.repo.UpdateCollection(id, collection)
	if err != nil {
		return nil, err
	}
	return collection, err

}

func (s *collectionService) DeleteCollection(id int) error {
	return s.repo.DeleteCollection(id)
}
