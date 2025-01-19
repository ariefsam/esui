package esui_test

import (
	"testing"

	"github.com/ariefsam/esui"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockEventstore struct {
	mock.Mock
}

func (m *mockEventstore) StoreEvent(aggregateID string, aggregateName string, eventName string, data interface{}) (err error) {
	m.Called(aggregateID, aggregateName, eventName, data)
	return
}

func TestNewEsui(t *testing.T) {
	estore := &mockEventstore{}
	esObj := esui.NewEsui(estore)
	require.NotNil(t, esObj)

	entityID, err := esObj.CreateEntity("user")
	require.NoError(t, err)
	require.NotEmpty(t, entityID)

	entityObj := esui.Entity{
		GeneralEntity: esui.GeneralEntity{
			ID:   entityID,
			Name: "user",
		},
	}

	require.Equal(t, entityObj, repoEntity.createCalled)

	err = esObj.AddEvent(entityID, "user", "created", map[string]any{
		"name": "string",
	})
	require.NoError(t, err)

}
