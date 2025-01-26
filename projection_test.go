package esui_test

import (
	"testing"

	"github.com/ariefsam/esui"
	"github.com/stretchr/testify/assert"
)

func TestProjection(t *testing.T) {
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	es := esui.NewEsui(estore, idgenerator)

	idgenerator.On("Generate").Return("xyz123").Once()
	estore.On("StoreEvent", "xyz123", "projection", "created", esui.EsuiProjectionCreated{
		Name: "test",
	}).Return(nil).Once()
	projID, err := es.CreateProjection("test")
	assert.NoError(t, err)
	assert.EqualValues(t, "xyz123", projID)
}
