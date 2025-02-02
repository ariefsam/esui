package esui

type Application struct {
	ID          ShortID
	Name        string
	Entity      map[EntityID]Entity
	Projections map[ProjectionID]Projection
	Version     ShortID
}

type Entity struct {
	Name   string
	Events map[ShortID]Event
}

type Event struct {
	Name      string
	Attribute map[AttributeName]AttributeType
}

type ProjectionID ShortID
type EntityID ShortID
type EntityEventName string

type Projection struct {
	ID          ShortID
	Name        string
	SubscribeTo map[EntityID]map[EntityEventName]bool
	Tables      []Table
	Blocks      []Block
}

type Table struct {
	Name    string
	Columns map[AttributeName]AttributeType
}
