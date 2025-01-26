# esui
```
Esui Package  
│
├── Structs  
│   ├── Esui  
│   │   ├── eventstore: eventstoreDB  
│   │   └── idgenerator: idgenerator  
│   ├── EsuiEntity  
│   │   ├── ID: ShortID  
│   │   ├── Name: string  
│   │   └── Events: map[string]EsuiEntityEvent  
│   ├── EsuiEntityEvent  
│   │   └── Attributes: map[AttributeName]AttributeType  
│   ├── EsuiEntityCreated  
│   │   └── Name: string  
│   ├── EsuiEventAdded  
│   │   └── Name: string  
│   ├── EsuiAttributeAdded  
│   │   ├── EventName: string  
│   │   ├── Name: AttributeName  
│   │   └── Type: AttributeType  
│   ├── EsuiProjection  
│   │   └── Name: string  
│   ├── EsuiEvent  
│   │   ├── EventID: ShortID  
│   │   ├── AggregateID: ShortID  
│   │   ├── AggregateName: string  
│   │   ├── EventName: string  
│   │   └── Data: string  
│   └── EsuiProjectionCreated  
│       └── Name: string  
│
├── Interfaces  
│   ├── idgenerator  
│   │   └── Generate() string  
│   └── eventstoreDB  
│       ├── StoreEvent(aggregateID, aggregateName, eventName, data) error  
│       └── FetchAggregateEvents(aggregateID, aggregateName, fromID) ([]EsuiEvent, error)  
│
├── Methods  
│   ├── NewEsui(eventstore eventstoreDB, idgenerator idgenerator) *Esui  
│   ├── Esui.CreateEntity(entityName string) (ShortID, error)  
│   ├── Esui.GetEntity(entityID ShortID) (EsuiEntity, error)  
│   ├── Esui.AddEventToEntity(entityID ShortID, eventName string) error  
│   ├── Esui.AddAttribute(entityID ShortID, eventName string, attributeName AttributeName, attributeType AttributeType) error  
│   ├── Esui.CreateProjection(projectionName string) (ShortID, error)  
│   ├── EsuiEntity.Created(event EsuiEvent, entityID ShortID)  
│   ├── EsuiEntity.EventAdded(event EsuiEvent)  
│   └── EsuiEntity.AttributeAdded(event EsuiEvent)  
│
├── Types  
│   ├── ShortID: string  
│   ├── AttributeName: string  
│   └── AttributeType: string  
│       └── Validate() error  
│
└── Dependencies  
    └── github.com/ariefsam/esui/logger
```
