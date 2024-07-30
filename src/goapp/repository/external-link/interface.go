package repository

import "main/model"

type ExternalLink interface {
	GetAll() ([]model.ExternalLink, error)
	GetByIsEnabled(isEnabled bool) ([]model.ExternalLink, error)
	GetByID(id int64) (*model.ExternalLink, error)
	Create(externalLink *model.ExternalLink) (*model.ExternalLink, error)
	Update(id int64, externalLink *model.ExternalLink) (*model.ExternalLink, error)
	Delete(id int64) error
}
