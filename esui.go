package esui

import (
	"context"
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
	ID       ShortID `json:"projection_id"`
	Name     string  `json:"name"`
	IsActive bool    `json:"is_active"`
}

type EsuiEvent struct {
	EventID       ShortID `json:"event_id"`
	AggregateID   ShortID `json:"aggregate_id"`
	AggregateName string  `json:"aggregate_name"`
	EventName     string  `json:"event_name"`
	Data          string  `json:"data"`
}

type eventstoreDB interface {
	StoreEvent(ctx context.Context, aggregateID string, aggregateName string, eventName string, data interface{}) (err error)
	FetchAggregateEvents(ctx context.Context, aggregateID string, aggregateName string, fromID string) (events []EsuiEvent, err error)
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

func (es *Esui) CreateEntity(ctx context.Context, entityName string) (entityID ShortID, err error) {
	entityObj := EsuiEntityCreated{
		Name: entityName,
	}
	entityID = ShortID(es.idgenerator.Generate())
	err = es.eventstore.StoreEvent(ctx, string(entityID), "entity", "created", entityObj)

	if err != nil {
		logger.Println(err)
		return "", err
	}
	return
}

func (es *Esui) GetEntity(ctx context.Context, entityID ShortID) (entity EsuiEntity, err error) {
	events, err := es.eventstore.FetchAggregateEvents(ctx, string(entityID), "entity", "")
	if err != nil {
		logger.Println(err)
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
		logger.Println(err)
		return
	}
	entity.ID = entityID
	entity.Name = entityCreated.Name
}

func (entity *EsuiEntity) EventAdded(event EsuiEvent) {
	var eventAdded EsuiEventAdded
	err := json.Unmarshal([]byte(event.Data), &eventAdded)
	if err != nil {
		logger.Println(err)
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
		logger.Println(err)
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

func (es *Esui) AddEventToEntity(ctx context.Context, entityID ShortID, eventName string) (err error) {
	entity, err := es.GetEntity(ctx, entityID)
	if err != nil {
		logger.Println(err)
		return
	}

	if entity.Name == "" {
		err = errors.New("entity not found")
		logger.Println(err)
		return
	}

	if _, ok := entity.Events[eventName]; ok {
		err = errors.New("event already exist")
		logger.Println(err)
		return
	}

	dataEvent := EsuiEventAdded{
		Name: eventName,
	}
	err = es.eventstore.StoreEvent(ctx, string(entityID), "entity", "event_added", dataEvent)

	return
}

func (es *Esui) AddAttribute(ctx context.Context, entityID ShortID, eventName string, attributeName AttributeName, attributeType AttributeType) (err error) {
	err = es.eventstore.StoreEvent(ctx, string(entityID), "entity", "attribute_added", EsuiAttributeAdded{
		EventName: eventName,
		Name:      attributeName,
		Type:      attributeType,
	})

	return
}

type EsuiProjectionCreated struct {
	Name string `json:"name"`
}

func (es *Esui) CreateProjection(ctx context.Context, projectionName string) (projectionID ShortID, err error) {
	projectionObj := EsuiProjectionCreated{
		Name: projectionName,
	}
	projectionID = ShortID(es.idgenerator.Generate())
	err = es.eventstore.StoreEvent(ctx, string(projectionID), "projection", "created", projectionObj)

	if err != nil {
		logger.Println(err)
		return "", err
	}
	return
}

func (es *Esui) GetProjection(ctx context.Context, projectionID ShortID) (projection EsuiProjection, err error) {
	events, err := es.eventstore.FetchAggregateEvents(ctx, string(projectionID), "projection", "")
	if err != nil {
		logger.Println(err)
		return
	}

	proj := EsuiProjection{}
	for _, event := range events {
		switch event.EventName {
		case "created":
			proj.HandleCreated(event, projectionID)
		}
	}
	projection = proj
	return
}

func (projection *EsuiProjection) HandleCreated(event EsuiEvent, projectionID ShortID) {
	var projectionCreated EsuiProjectionCreated
	err := json.Unmarshal([]byte(event.Data), &projectionCreated)
	if err != nil {
		logger.Println(err)
		return
	}
	projection.ID = projectionID
	projection.Name = projectionCreated.Name
}

type EsuiTableCreated struct {
	Name string `json:"name"`
}

func (es *Esui) CreateTable(ctx context.Context, projectionID ShortID, tableName string) (err error) {

	projection, err := es.GetProjection(ctx, projectionID)
	if err != nil {
		logger.Println(err)
		return
	}

	if projection.Name == "" {
		err = errors.New("projection not found")
		logger.Println(err)
		return
	}

	err = es.eventstore.StoreEvent(ctx, string(projectionID), "projection", "table_created", EsuiTableCreated{
		Name: tableName,
	})

	return
}
