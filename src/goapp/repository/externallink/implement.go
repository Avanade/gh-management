package externallink

import (
	"database/sql"
	"main/model"
	"main/repository"
	"time"
)

type externalLinkRepository struct {
	repository.Database
}

func NewExternalLinkRepository(db repository.Database) ExternalLinkRepository {
	return &externalLinkRepository{
		Database: db,
	}
}

func (r *externalLinkRepository) GetAll() ([]model.ExternalLink, error) {
	var externalLinks []model.ExternalLink
	rows, err := r.Query("[dbo].[usp_ExternalLink_Select]")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var externalLink model.ExternalLink

		externalLink.ID = v["Id"].(int64)
		externalLink.IconSVGPath = v["IconSVG"].(string)
		externalLink.Hyperlink = v["Hyperlink"].(string)
		externalLink.DisplayName = v["LinkName"].(string)
		externalLink.IsEnabled = v["IsEnabled"].(bool)
		externalLink.Created = v["Created"].(time.Time)
		externalLink.CreatedBy = v["CreatedBy"].(string)
		externalLink.Modified = v["Modified"].(time.Time)
		externalLink.ModifiedBy = v["ModifiedBy"].(string)

		externalLinks = append(externalLinks, externalLink)
	}

	return externalLinks, nil
}

func (r *externalLinkRepository) GetByIsEnabled(isEnabled bool) ([]model.ExternalLink, error) {
	var externalLinks []model.ExternalLink
	rows, err := r.Query("[dbo].[usp_ExternalLink_Select_ByIsEnabled]",
		sql.Named("IsEnabled", isEnabled))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, mapRow := range mapRows {
		var externalLink model.ExternalLink

		externalLink.ID = mapRow["Id"].(int64)
		externalLink.IconSVGPath = mapRow["IconSVG"].(string)
		externalLink.Hyperlink = mapRow["Hyperlink"].(string)
		externalLink.DisplayName = mapRow["LinkName"].(string)
		externalLink.IsEnabled = mapRow["IsEnabled"].(bool)
		externalLink.Created = mapRow["Created"].(time.Time)
		externalLink.CreatedBy = mapRow["CreatedBy"].(string)
		externalLink.Modified = mapRow["Modified"].(time.Time)
		externalLink.ModifiedBy = mapRow["ModifiedBy"].(string)

		externalLinks = append(externalLinks, externalLink)
	}

	return externalLinks, nil
}

func (r *externalLinkRepository) GetByID(id int64) (*model.ExternalLink, error) {
	var externalLink model.ExternalLink
	rows, err := r.Query("[dbo].[usp_ExternalLink_Select_ById]",
		sql.Named("Id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	mapRow := mapRows[0]
	externalLink.ID = mapRow["Id"].(int64)
	externalLink.IconSVGPath = mapRow["IconSVG"].(string)
	externalLink.Hyperlink = mapRow["Hyperlink"].(string)
	externalLink.DisplayName = mapRow["LinkName"].(string)
	externalLink.IsEnabled = mapRow["IsEnabled"].(bool)
	externalLink.Created = mapRow["Created"].(time.Time)
	externalLink.CreatedBy = mapRow["CreatedBy"].(string)
	externalLink.Modified = mapRow["Modified"].(time.Time)
	externalLink.ModifiedBy = mapRow["ModifiedBy"].(string)

	return &externalLink, nil
}

func (r *externalLinkRepository) Create(externalLink *model.ExternalLink) (*model.ExternalLink, error) {
	result, err := r.QueryRow("[dbo].[usp_ExternalLink_Insert]",
		sql.Named("LinkName", externalLink.DisplayName),
		sql.Named("IconSVG", externalLink.IconSVGPath),
		sql.Named("Hyperlink", externalLink.Hyperlink),
		sql.Named("IsEnabled", externalLink.IsEnabled),
		sql.Named("CreatedBy", externalLink.CreatedBy))
	if err != nil {
		return nil, err
	}
	var id int64
	var created time.Time
	err = result.Scan(&id, &created)
	if err != nil {
		return nil, err
	}
	externalLink.ID = id
	externalLink.Created = created
	return externalLink, nil
}

func (r *externalLinkRepository) Update(id int64, externalLink *model.ExternalLink) (*model.ExternalLink, error) {
	err := r.Execute("[dbo].[usp_ExternalLink_Update]",
		sql.Named("Id", id),
		sql.Named("LinkName", externalLink.DisplayName),
		sql.Named("IconSVG", externalLink.IconSVGPath),
		sql.Named("Hyperlink", externalLink.Hyperlink),
		sql.Named("IsEnabled", externalLink.IsEnabled),
		sql.Named("ModifiedBy", externalLink.ModifiedBy))
	if err != nil {
		return nil, err
	}
	return externalLink, nil
}

func (r *externalLinkRepository) Delete(id int64) error {
	err := r.Execute("[dbo].[usp_ExternalLink_Delete]",
		sql.Named("Id", id))
	if err != nil {
		return err
	}
	return nil
}
