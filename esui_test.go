package esui_test

import (
	"testing"

	"github.com/ariefsam/esui"
	"github.com/stretchr/testify/require"
)

type mockRepository[T esui.EsuiEntity] struct {
	createCalled T
	updateCalled bool
}

func (m *mockRepository[T]) Create(id string, object T) (err error) {
	m.createCalled = object
	return
}

type mockEventstore struct{}

func (m *mockEventstore) StoreEvent(aggregateID string, aggregateName string, eventName string, data interface{}) (err error) {
	return
}

func TestNewEsui(t *testing.T) {
	repoEntity := &mockRepository[esui.Entity]{}
	repoProjection := &mockRepository[esui.Projection]{}
	estore := &mockEventstore{}
	esObj := esui.NewEsui(repoEntity, repoProjection, estore)
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
