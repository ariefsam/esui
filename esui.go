package esui

import (
	"encoding/json"
	"errors"

	"github.com/ariefsam/esui/logger"
)

type Esui struct {
	eventstore eventstoreDB
	idgenerator
}

type idgenerator interface {
	Generate() string
}

type AttributeName string
type AttributeType string

func (atype AttributeType) Validate() error {
	if atype != "string" && atype != "int" {
		return errors.New("Invalid attribute type")
	}
	return nil
}

type ShortID string

type EsuiEntity struct {
	ID     ShortID                    `json:"entity_id"`
	Name   string                     `json:"name"`
	Events map[string]EsuiEntityEvent `json:"events"`
}

type EsuiEntityEvent struct {
	Attributes map[AttributeName]AttributeType `json:"attribute"`
}

type EsuiEntityCreated struct {
	Name string `json:"name"`
}

type EsuiEventAdded struct {
	Name string `json:"name"`
}

type EsuiAttributeAdded struct {
	EventName string        `json:"event_name"`
	Name      AttributeName `json:"name"`
	Type      AttributeType `json:"type"`
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
			entity.Created(event, entityID)
		case "event_added":
			entity.EventAdded(event)
		case "attribute_added":
			entity.AttributeAdded(event)
		}
	}

	return
}

func (entity *EsuiEntity) Created(event EsuiEvent, entityID ShortID) {
	var entityCreated EsuiEntityCreated
	err := json.Unmarshal([]byte(event.Data), &entityCreated)
	if err != nil {
		return
	}
	entity.ID = entityID
	entity.Name = entityCreated.Name
}

func (entity *EsuiEntity) EventAdded(event EsuiEvent) {
	var eventAdded EsuiEventAdded
	err := json.Unmarshal([]byte(event.Data), &eventAdded)
	if err != nil {
		return
	}
	if entity.Events == nil {
		entity.Events = make(map[string]EsuiEntityEvent)
	}
	entity.Events[eventAdded.Name] = EsuiEntityEvent{}
}

func (entity *EsuiEntity) AttributeAdded(event EsuiEvent) {
	var attributeAdded EsuiAttributeAdded
	err := json.Unmarshal([]byte(event.Data), &attributeAdded)
	if err != nil {
		return
	}
	if entity.Events == nil {
		entity.Events = make(map[string]EsuiEntityEvent)
	}
	if entity.Events[attributeAdded.EventName].Attributes == nil {
		entity.Events[attributeAdded.EventName] = EsuiEntityEvent{
			Attributes: make(map[AttributeName]AttributeType),
		}
	}
	entity.Events[attributeAdded.EventName].Attributes[attributeAdded.Name] = attributeAdded.Type
}

func (es *Esui) AddEventToEntity(entityID ShortID, eventName string) (err error) {
	entity, err := es.GetEntity(entityID)
	if err != nil {
		logger.Println(err)
		return
	}

	if entity.Name == "" {
		err = errors.New("entity not found")
		return
	}

	if _, ok := entity.Events[eventName]; ok {
		err = errors.New("event already exist")
		return
	}

	dataEvent := EsuiEventAdded{
		Name: eventName,
	}
	err = es.eventstore.StoreEvent(string(entityID), "entity", "event_added", dataEvent)

	return
}

func (es *Esui) AddAttribute(entityID ShortID, eventName string, attributeName AttributeName, attributeType AttributeType) (err error) {
	err = es.eventstore.StoreEvent(string(entityID), "entity", "attribute_added", EsuiAttributeAdded{
		EventName: eventName,
		Name:      attributeName,
		Type:      attributeType,
	})

	return
}
