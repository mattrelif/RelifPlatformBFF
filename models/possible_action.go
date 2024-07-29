package models

import "relif/bff/entities"

type PossibleAction struct {
	Name         string   `bson:"_id"`
	AllowedRoles []string `bson:"allowed_roles"`
}

func (action *PossibleAction) ToEntity() entities.PossibleAction {
	return entities.PossibleAction{
		Name:         action.Name,
		AllowedRoles: action.AllowedRoles,
	}
}
