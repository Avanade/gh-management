package externallink

import (
	"errors"
	"main/model"
	"main/repository"
	"strconv"
)

type externalLinkService struct {
	Repository *repository.Repository
}

func NewExternalLinkService(repo *repository.Repository) ExternalLinkService {
	return &externalLinkService{
		Repository: repo,
	}
}

func (s externalLinkService) GetAll() ([]model.ExternalLink, error) {
	return s.Repository.ExternalLink.Select()
}

func (s externalLinkService) GetAllEnabled() ([]model.ExternalLink, error) {
	return s.Repository.ExternalLink.SelectByIsEnabled(true)
}

func (s externalLinkService) GetByID(id string) (*model.ExternalLink, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.Repository.ExternalLink.SelectByID(parsedId)
}

func (s externalLinkService) Create(externalLink *model.ExternalLink) (*model.ExternalLink, error) {
	return s.Repository.ExternalLink.Insert(externalLink)
}

func (s externalLinkService) Update(id string, externalLink *model.ExternalLink) (*model.ExternalLink, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.Repository.ExternalLink.Update(parsedId, externalLink)
}

func (s externalLinkService) Delete(id string) error {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	return s.Repository.ExternalLink.Delete(parsedId)
}

func (s externalLinkService) Validate(externalLink *model.ExternalLink) error {
	if externalLink == nil {
		return errors.New("external link is empty")
	}
	if externalLink.DisplayName == "" {
		return errors.New("display name is required")
	}
	if externalLink.IconSVGPath == "" {
		return errors.New("icon svg path is required")
	}
	if externalLink.Hyperlink == "" {
		return errors.New("hyperlink is required")
	}
	return nil
}
