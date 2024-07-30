package service

import (
	"errors"
	"main/model"
	repositoryExternalLink "main/repository/external-link"
	"strconv"
)

type externalLinkService struct {
	repositoryExternalLink repositoryExternalLink.ExternalLink
}

func NewExternalLinkService(repositoryExternalLink repositoryExternalLink.ExternalLink) ExternalLinkService {
	return &externalLinkService{
		repositoryExternalLink: repositoryExternalLink,
	}
}

func (s externalLinkService) GetAll() ([]model.ExternalLink, error) {
	return s.repositoryExternalLink.GetAll()
}

func (s externalLinkService) GetAllEnabled() ([]model.ExternalLink, error) {
	return s.repositoryExternalLink.GetByIsEnabled(true)
}

func (s externalLinkService) GetByID(id string) (*model.ExternalLink, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repositoryExternalLink.GetByID(parsedId)
}

func (s externalLinkService) Create(externalLink *model.ExternalLink) (*model.ExternalLink, error) {
	return s.repositoryExternalLink.Create(externalLink)
}

func (s externalLinkService) Update(id string, externalLink *model.ExternalLink) (*model.ExternalLink, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repositoryExternalLink.Update(parsedId, externalLink)
}

func (s externalLinkService) Delete(id string) error {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	return s.repositoryExternalLink.Delete(parsedId)
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
