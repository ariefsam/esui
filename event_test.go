package esui_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ariefsam/esui"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockEventstore struct {
	mock.Mock
}

func (m *mockEventstore) StoreEvent(ctx context.Context, aggregateID string, aggregateName string, eventName string, data interface{}) (err error) {
	args := m.Called(aggregateID, aggregateName, eventName, data)
	if len(args) == 0 {
		return nil
	}
	return args.Error(0)
}

func (m *mockEventstore) FetchAggregateEvents(ctx context.Context, aggregateID string, aggregateName string, fromID string) (events []esui.EstoreEvent, err error) {
	args := m.Called(aggregateID, aggregateName, fromID)
	if len(args) == 0 {
		return []esui.EstoreEvent{}, nil
	}
	events, _ = args.Get(0).([]esui.EstoreEvent)
	err = args.Error(1)
	return
}

type mockIDGenerator struct {
	mock.Mock
}

func (m *mockIDGenerator) Generate() string {
	args := m.Called()
	return args.String(0)
}

func TestNewEntity(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	esObj := esui.NewEsui(estore, idgenerator)
	require.NotNil(t, esObj)

	t.Run("Create Entity Success", func(t *testing.T) {
		idgenerator.On("Generate").Return("abc123").Once()
		expectedEntityObj := esui.EsuiEntityCreated{
			Name: "user",
		}
		estore.On("StoreEvent", "abc123", "entity", "created", expectedEntityObj).Return(nil)

		entityID, err := esObj.CreateEntity(ctx, "user")
		require.NoError(t, err)
		require.NotEmpty(t, entityID)

		estore.AssertCalled(t, "StoreEvent", "abc123", "entity", "created", mock.MatchedBy(func(data interface{}) bool {
			actualEntity, ok := data.(esui.EsuiEntityCreated)
			if !ok {
				return false
			}
			return actualEntity.Name == expectedEntityObj.Name
		}))

		require.NoError(t, err)
	})

	t.Run("Create Entity Failed Eventstore", func(t *testing.T) {
		idgenerator.On("Generate").Return("abc123").Once()
		expectedEntityObj := esui.EsuiEntityCreated{
			Name: "userx",
		}
		estore.On("StoreEvent", "abc123", "entity", "created", expectedEntityObj).Return(errors.New("Error store event"))

		entityID, err := esObj.CreateEntity(ctx, "userx")
		require.Error(t, err)
		require.Empty(t, entityID)
	})
}

func TestGetEntity(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	esObj := esui.NewEsui(estore, idgenerator)
	require.NotNil(t, esObj)

	t.Run("Get Entity Success", func(t *testing.T) {

		estore.On("FetchAggregateEvents", "abc123", "entity", "").Return(
			[]esui.EstoreEvent{
				{
					EventID:       "abc123",
					AggregateID:   "abc123",
					AggregateName: "entity",
					EventName:     "created",
					Data:          `{"name":"user"}`,
				},
			}, nil).Once()
		esuiEntity, err := esObj.GetEntity(ctx, "abc123")
		require.NoError(t, err)
		require.Equal(t, "user", esuiEntity.Name)
	})

	t.Run("Get Entity Failed FetchAggregateEvents", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "bc123", "entity", "").Return(esui.EsuiEntity{}, errors.New("Error fetch aggregate events")).Once()
		esuiEntity, err := esObj.GetEntity(ctx, "bc123")
		require.Error(t, err)
		require.Empty(t, esuiEntity)
	})
}

func TestAddEventToEntity(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	esObj := esui.NewEsui(estore, idgenerator)
	require.NotNil(t, esObj)

	t.Run("Cannot Add Event if Entity Not Found", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "notfoundIDxxx", "entity", "").Return([]esui.EstoreEvent{}, nil).Once()
		err := esObj.AddEventToEntity(ctx, "notfoundIDxxx", "event_added")
		require.Error(t, err)
	})

	t.Run("Add Event To Entity Success", func(t *testing.T) {
		idgenerator.On("Generate").Return("abc123").Once()

		estore.On("FetchAggregateEvents", "abc123", "entity", "").Return(
			[]esui.EstoreEvent{
				{
					EventID:       "abc123",
					AggregateID:   "abc123",
					AggregateName: "entity",
					EventName:     "created",
					Data:          `{"name":"user"}`,
				},
			}, nil).Once()

		estore.On("StoreEvent", "abc123", "entity", "event_added", esui.EsuiEventAdded{
			Name: "user_created",
		}).Return(nil).Once()
		err := esObj.AddEventToEntity(ctx, "abc123", "user_created")
		require.NoError(t, err)

		estore.AssertCalled(t, "StoreEvent", "abc123", "entity", "event_added", mock.MatchedBy(func(data interface{}) bool {
			dataEvent, ok := data.(esui.EsuiEventAdded)
			if !ok {
				return false
			}
			return dataEvent.Name == "user_created"
		}))
	})

	t.Run("Get entity will show event", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "abc123", "entity", "").Return(
			[]esui.EstoreEvent{
				{
					EventID:       "abc123",
					AggregateID:   "abc123",
					AggregateName: "entity",
					EventName:     "created",
					Data:          `{"name":"user"}`,
				},
				{
					EventID:       "abc124",
					AggregateID:   "abc123",
					AggregateName: "entity",
					EventName:     "event_added",
					Data:          `{"name":"product_created"}`,
				},
			}, nil).Once()
		esuiEntity, err := esObj.GetEntity(ctx, "abc123")
		require.NoError(t, err)
		require.Equal(t, "user", esuiEntity.Name)
		require.Equal(t, esui.EsuiEntityEvent{}, esuiEntity.Events["product_created"])
	})

	t.Run("Failed to create already exist event on entity product", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "prod123", "entity", "").Return([]esui.EstoreEvent{
			{
				EventID:       "abc123",
				AggregateID:   "prod123",
				AggregateName: "entity",
				EventName:     "created",
				Data:          `{"name":"product"}`,
			},
			{
				EventID:       "abc124",
				AggregateID:   "prod123",
				AggregateName: "entity",
				EventName:     "event_added",
				Data:          `{"name":"product_created"}`,
			},
		}, nil).Once()

		err := esObj.AddEventToEntity(ctx, "prod123", "product_created")
		require.Error(t, err)
	})

}

func TestAddAttribute(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	esObj := esui.NewEsui(estore, idgenerator)
	require.NotNil(t, esObj)

	// aggregate: product
	// entity: prod123
	// event: product_created
	// attribute: name string
	t.Run("Add Attribute Success", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "prod123", "entity", "").Return([]esui.EstoreEvent{
			{
				EventID:       "abc123",
				AggregateID:   "prod123",
				AggregateName: "entity",
				EventName:     "created",
				Data:          `{"name":"product"}`,
			},
			{
				EventID:       "abc124",
				AggregateID:   "prod123",
				AggregateName: "entity",
				EventName:     "event_added",
				Data:          `{"name":"product_created"}`,
			},
		}, nil).Once()

		estore.On("StoreEvent", "prod123", "entity", "attribute_added", esui.EsuiAttributeAdded{
			EventName: "product_created",
			Name:      "name",
			Type:      "string",
		}).Return(nil).Once()

		err := esObj.AddAttribute(ctx, "prod123", "product_created", "name", "string")
		require.NoError(t, err)

		estore.AssertCalled(t, "StoreEvent", "prod123", "entity", "attribute_added", mock.MatchedBy(func(data interface{}) bool {
			dataEvent, ok := data.(esui.EsuiAttributeAdded)
			if !ok {
				return false
			}
			return dataEvent.Name == "name" && dataEvent.Type == "string" && dataEvent.EventName == "product_created"
		}))
	})

	t.Run("Get entity will show attribute", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "prod1234", "entity", "").Return([]esui.EstoreEvent{
			{
				EventID:       "abc123",
				AggregateID:   "prod1234",
				AggregateName: "entity",
				EventName:     "created",
				Data:          `{"name":"product"}`,
			},
			{
				EventID:       "abc124",
				AggregateID:   "prod1234",
				AggregateName: "entity",
				EventName:     "event_added",
				Data:          `{"name":"product_created"}`,
			},
			{
				EventID:       "abc125",
				AggregateID:   "prod1234",
				AggregateName: "entity",
				EventName:     "attribute_added",
				Data:          `{"event_name":"product_created","name":"name","type":"string"}`,
			},
			{
				EventID:       "abc126",
				AggregateID:   "prod1234",
				AggregateName: "entity",
				EventName:     "attribute_added",
				Data:          `{"event_name":"product_created","name":"price","type":"float"}`,
			},
		}, nil).Once()
		esuiEntity, err := esObj.GetEntity(ctx, "prod1234")
		require.NoError(t, err)
		require.Equal(t, "product", esuiEntity.Name)
		require.Equal(t, esui.AttributeType("string"), esuiEntity.Events["product_created"].Attributes["name"])
		require.Equal(t, esui.AttributeType("float"), esuiEntity.Events["product_created"].Attributes["price"])
	})

}
