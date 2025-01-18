package esui

import "github.com/ariefsam/esui/idgenerator"

type Esui struct {
	entityRepository     repository[Entity]
	projectionRepository repository[Projection]
	eventstore           eventstoreDB
}

type GeneralEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type attributeName string
type attributeType string

type Entity struct {
	GeneralEntity
}

type Projection struct {
	GeneralEntity
}
type EsuiEntity interface {
	Entity | Projection
}

type repository[T EsuiEntity] interface {
	Create(id string, object T) (err error)
}

type eventstoreDB interface {
	StoreEvent(aggregateID string, aggregateName string, eventName string, data interface{}) (err error)
}

func NewEsui(
	entityRepository repository[Entity],
	projectionRepository repository[Projection],
	eventstore eventstoreDB,
) (obj *Esui) {
	obj = &Esui{
		entityRepository:     entityRepository,
		projectionRepository: projectionRepository,
		eventstore:           eventstore,
	}
	return obj
}

func (es *Esui) CreateEntity(entityName string) (id string, err error) {
	entityObj := Entity{
		GeneralEntity: GeneralEntity{
			ID:   idgenerator.Generate(),
			Name: entityName,
		},
	}
	id = entityObj.ID
	err = es.entityRepository.Create(entityObj.ID, entityObj)
	if err != nil {
		return
	}

	return
}
