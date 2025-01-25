package esui

import (
	"encoding/json"
	"errors"
)

type Esui struct {
	eventstore eventstoreDB
	idgenerator
}

type idgenerator interface {
	Generate() string
}

type attributeName string
type attributeType string

func (atype attributeType) Validate() error {
	if atype != "string" && atype != "int" {
		return errors.New("Invalid attribute type")
	}
	return nil
}

type ShortID string

type EsuiEntity struct {
	Name string `json:"name"`
}

type EsuiEntityCreated struct {
	Name string `json:"name"`
}

type EsuiAttributeAdded struct {
	Name attributeName `json:"name"`
	Type attributeType `json:"type"`
}

type EsuiProjection struct {
	Name string `json:"name"`
}

type EsuiEvent struct {
	EventID       ShortID `json:"event_id"`
	AggregateID   ShortID `json:"aggregate_id"`
	AggregateName string  `json:"aggregate_name"`
	EventName     string  `json:"event_name"`
	Data          string  `json:"data"`
}

type eventstoreDB interface {
	StoreEvent(aggregateID string, aggregateName string, eventName string, data interface{}) (err error)
	FetchAggregateEvents(aggregateID string, aggregateName string, fromID string) (events []EsuiEvent, err error)
}

func NewEsui(
	eventstore eventstoreDB,
	idgenerator idgenerator,
) (obj *Esui) {
	obj = &Esui{
		eventstore:  eventstore,
		idgenerator: idgenerator,
	}
	return obj
}

func (es *Esui) CreateEntity(entityName string) (entityID ShortID, err error) {
	entityObj := EsuiEntityCreated{
		Name: entityName,
	}
	entityID = ShortID(es.idgenerator.Generate())
	err = es.eventstore.StoreEvent(string(entityID), "entity", "created", entityObj)

	if err != nil {
		return "", err
	}
	return
}

func (es *Esui) GetEntity(entityID ShortID) (entity EsuiEntity, err error) {
	events, err := es.eventstore.FetchAggregateEvents(string(entityID), "entity", "")
	if err != nil {
		return
	}

	for _, event := range events {
		switch event.EventName {
		case "created":
			var entityCreated EsuiEntityCreated
			err = json.Unmarshal([]byte(event.Data), &entityCreated)
			if err != nil {
				return
			}
			entity.Name = entityCreated.Name
			continue
		}
	}
	return
}
