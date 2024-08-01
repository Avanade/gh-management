package externallink

import "main/model"

type ExternalLinkRepository interface {
	Select() ([]model.ExternalLink, error)
	SelectByIsEnabled(isEnabled bool) ([]model.ExternalLink, error)
	SelectByID(id int64) (*model.ExternalLink, error)
	Insert(externalLink *model.ExternalLink) (*model.ExternalLink, error)
	Update(id int64, externalLink *model.ExternalLink) (*model.ExternalLink, error)
	Delete(id int64) error
}
