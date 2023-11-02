package entity

import "github.com/google/uuid"

type Entity struct {
	ID       string       `bson:"_id,omitempty" validate:"required"`
	Creation *TimeEvent   `validate:"required"`
	Updated  []*TimeEvent `validate:"required,gte=1"`
	Deleted  *TimeEvent   `validate:"omitempty,required"`
}

func (entity *Entity) SetCreated(userID string) {
	event := newTimeEvent(userID)

	entity.Creation = event
	entity.Updated = []*TimeEvent{event}
}

func (entity *Entity) SetUpdated(userID string) {
	entity.Updated = append(entity.Updated, newTimeEvent(userID))
}

func (entity *Entity) SetDeleted(userID string) {
	event := newTimeEvent(userID)

	entity.Deleted = newTimeEvent(userID)
	entity.Updated = append(entity.Updated, event)
}

func (entity *Entity) IsDeleted() bool {
	return entity.Deleted != nil
}

func (entity *Entity) WasCreatedBy(userID string) bool {
	if entity.Creation == nil {
		return false
	}

	return entity.Creation.UserID == userID
}

func New(userID string) Entity {
	e := Entity{ID: uuid.NewString()}
	e.SetCreated(userID)
	return e
}
