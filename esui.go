package esui

import (
	"encoding/json"
	"log"
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
type ShortID string

type EsuiEntity struct {
	Name string `json:"name"`
}

type EsuiEntityCreated struct {
	Name string `json:"name"`
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
	log.Println("err", err)
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
