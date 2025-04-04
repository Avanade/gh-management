package model

import "time"

type AdoOrganizationRequest struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Purpose   string    `json:"purpose"`
	Created   time.Time `json:"created"`
	CreatedBy string    `json:"createdBy"`
}
