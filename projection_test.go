package esui_test

import (
	"context"
	"testing"

	"github.com/ariefsam/esui"
	"github.com/stretchr/testify/assert"
)

func TestProjection(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	es := esui.NewEsui(estore, idgenerator)

	idgenerator.On("Generate").Return("xyz123").Once()
	estore.On("StoreEvent", "xyz123", "projection", "created", esui.EsuiProjectionCreated{
		Name: "projection1",
	}).Return(nil).Once()
	projID, err := es.CreateProjection(ctx, "projection1")
	assert.NoError(t, err)
	assert.EqualValues(t, "xyz123", projID)

	estore.AssertCalled(t, "StoreEvent", "xyz123", "projection", "created", esui.EsuiProjectionCreated{
		Name: "projection1",
	})

}

func TestGetProjection(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	es := esui.NewEsui(estore, idgenerator)

	estore.On("FetchAggregateEvents", "proj1", "projection", "").Return([]esui.EsuiEvent{
		{
			EventID:       "1",
			AggregateID:   "proj1",
			AggregateName: "projection",
			EventName:     "created",
			Data:          `{"name":"projection1"}`,
		},
	}, nil)

	projection, err := es.GetProjection(ctx, "proj1")
	assert.NoError(t, err)
	assert.EqualValues(t, projection.Name, "projection1")
	assert.EqualValues(t, projection.ID, "proj1")
}
