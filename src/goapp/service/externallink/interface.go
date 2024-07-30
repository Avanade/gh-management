package externallink

import "main/model"

type ExternalLinkService interface {
	GetAll() ([]model.ExternalLink, error)
	GetAllEnabled() ([]model.ExternalLink, error)
	GetByID(id string) (*model.ExternalLink, error)
	Create(externalLink *model.ExternalLink) (*model.ExternalLink, error)
	Update(id string, externalLink *model.ExternalLink) (*model.ExternalLink, error)
	Delete(id string) error
	Validate(externalLink *model.ExternalLink) error
}
