package esui_test

import (
	"errors"
	"testing"

	"github.com/ariefsam/esui"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockEventstore struct {
	mock.Mock
}

func (m *mockEventstore) StoreEvent(aggregateID string, aggregateName string, eventName string, data interface{}) (err error) {
	args := m.Called(aggregateID, aggregateName, eventName, data)
	return args.Error(0)
}

func (m *mockEventstore) FetchAggregateEvents(aggregateID string, aggregateName string, fromID string) (events []esui.EsuiEvent, err error) {
	args := m.Called(aggregateID, aggregateName, fromID)
	events, _ = args.Get(0).([]esui.EsuiEvent)
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

		entityID, err := esObj.CreateEntity("user")
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

		entityID, err := esObj.CreateEntity("userx")
		require.Error(t, err)
		require.Empty(t, entityID)
	})
}

func TestGetEntity(t *testing.T) {
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	esObj := esui.NewEsui(estore, idgenerator)
	require.NotNil(t, esObj)

	t.Run("Get Entity Success", func(t *testing.T) {

		estore.On("FetchAggregateEvents", "abc123", "entity", "").Return(
			[]esui.EsuiEvent{
				{
					EventID:       "abc123",
					AggregateID:   "abc123",
					AggregateName: "entity",
					EventName:     "created",
					Data:          `{"name":"user"}`,
				},
			}, nil).Once()
		esuiEntity, err := esObj.GetEntity("abc123")
		require.NoError(t, err)
		require.Equal(t, "user", esuiEntity.Name)
	})

	t.Run("Get Entity Failed FetchAggregateEvents", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "bc123", "entity", "").Return(esui.EsuiEntity{}, errors.New("Error fetch aggregate events")).Once()
		esuiEntity, err := esObj.GetEntity("bc123")
		require.Error(t, err)
		require.Empty(t, esuiEntity)
	})
}
