package esui_test

import (
	"context"
	"strings"
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

	estore.On("FetchAggregateEvents", "proj1", "projection", "").Return([]esui.EstoreEvent{
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

func TestCreateTable(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	es := esui.NewEsui(estore, idgenerator)

	t.Run("Create Table Unknown Projection", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "proj11", "projection", "").Return([]esui.EstoreEvent{}, nil)

		err := es.CreateTable(ctx, "proj11", "table1")
		assert.Error(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "projection not found"))
	})

	t.Run("Create Table Success", func(t *testing.T) {

		estore.On("FetchAggregateEvents", "proj1", "projection", "").Return([]esui.EstoreEvent{
			{
				EventID:       "1",
				AggregateID:   "proj1",
				AggregateName: "projection",
				EventName:     "created",
				Data:          `{"name":"projection1"}`,
			},
		}, nil)

		estore.On("StoreEvent", "proj1", "projection", "table_created", esui.EsuiTableCreated{
			Name: "table1",
		}).Return(nil).Once()

		err := es.CreateTable(ctx, "proj1", "table1")
		assert.NoError(t, err)

		estore.AssertCalled(t, "StoreEvent", "proj1", "projection", "table_created", esui.EsuiTableCreated{
			Name: "table1",
		})
	})

	t.Run("Get Projection With Table", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "proj1x", "projection", "").Return([]esui.EstoreEvent{
			{
				EventID:       "1",
				AggregateID:   "proj1x",
				AggregateName: "projection",
				EventName:     "created",
				Data:          `{"name":"projection1"}`,
			},
			{
				EventID:       "2",
				AggregateID:   "proj1x",
				AggregateName: "projection",
				EventName:     "table_created",
				Data:          `{"name":"table1"}`,
			},
		}, nil)

		projection, err := es.GetProjection(ctx, "proj1x")
		assert.NoError(t, err)
		assert.EqualValues(t, "projection1", projection.Name)
		assert.EqualValues(t, "proj1x", projection.ID)
		assert.EqualValues(t, "table1", projection.Tables["table1"].Name)
		assert.EqualValues(t, "proj1x", projection.Tables["table1"].ProjectionID)
	})

}

func TestAddColumn(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	es := esui.NewEsui(estore, idgenerator)

	t.Run("Add Column Unknown Projection", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "proj11", "projection", "").Return([]esui.EstoreEvent{}, nil)

		err := es.AddColumn(ctx, "proj11", "table1", "column1", "string")
		assert.Error(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "projection not found"))
	})

	t.Run("Add Column Success", func(t *testing.T) {

		estore.On("FetchAggregateEvents", "proj1", "projection", "").Return([]esui.EstoreEvent{
			{
				EventID:       "1",
				AggregateID:   "proj1",
				AggregateName: "projection",
				EventName:     "created",
				Data:          `{"name":"projection1"}`,
			},
			{
				EventID:       "2",
				AggregateID:   "proj1",
				AggregateName: "projection",
				EventName:     "table_created",
				Data:          `{"name":"table1"}`,
			},
		}, nil)

		estore.On("StoreEvent", "proj1", "projection", "column_added", esui.EsuiColumnAdded{
			TableName:  "table1",
			ColumnName: "column1",
			ColumnType: "string",
		}).Return(nil).Once()

		err := es.AddColumn(ctx, "proj1", "table1", "column1", "string")
		assert.NoError(t, err)

		estore.AssertCalled(t, "StoreEvent", "proj1", "projection", "column_added", esui.EsuiColumnAdded{
			TableName:  "table1",
			ColumnName: "column1",
			ColumnType: "string",
		})
	})

	t.Run("Add Column Unknown Table", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "proj1", "projection", "").Return([]esui.EstoreEvent{
			{
				EventID:       "1",
				AggregateID:   "proj1",
				AggregateName: "projection",
				EventName:     "created",
				Data:          `{"name":"projection1"}`,
			},
		}, nil)

		err := es.AddColumn(ctx, "proj1", "unknown_table", "column1", "string")
		assert.Error(t, err)
		assert.True(t, strings.HasPrefix(err.Error(), "table not found"))
	})

	t.Run("Get Projection With Column", func(t *testing.T) {
		estore.On("FetchAggregateEvents", "proj123", "projection", "").Return([]esui.EstoreEvent{
			{
				EventID:       "1",
				AggregateID:   "proj123",
				AggregateName: "projection",
				EventName:     "created",
				Data:          `{"name":"projection1"}`,
			},
			{
				EventID:       "2",
				AggregateID:   "proj123",
				AggregateName: "projection",
				EventName:     "table_created",
				Data:          `{"name":"table1"}`,
			},
			{
				EventID:       "3",
				AggregateID:   "proj123",
				AggregateName: "projection",
				EventName:     "column_added",
				Data:          `{"table_name":"table1","column_name":"column1","column_type":"string"}`,
			},
		}, nil)

		projection, err := es.GetProjection(ctx, "proj123")
		assert.NoError(t, err)
		assert.EqualValues(t, "projection1", projection.Name)
		assert.EqualValues(t, "proj123", projection.ID)
		assert.EqualValues(t, "table1", projection.Tables["table1"].Name)
		assert.EqualValues(t, "proj123", projection.Tables["table1"].ProjectionID)
		assert.EqualValues(t, "column1", projection.Tables["table1"].Columns["column1"].Name)
		assert.EqualValues(t, "string", projection.Tables["table1"].Columns["column1"].Type)
	})

}

func TestAddBlockJavascriptToProjection(t *testing.T) {
	ctx := context.TODO()
	estore := &mockEventstore{}
	idgenerator := &mockIDGenerator{}
	es := esui.NewEsui(estore, idgenerator)

	data := esui.Block{
		BlockID:      "blockxxx",
		Name:         "script 1",
		Type:         "javascript",
		OrderedAfter: "",
	}

	estore.On("FetchAggregateEvents", "proj1", "projection", "").Return([]esui.EstoreEvent{
		{
			EventID:       "1",
			AggregateID:   "proj1",
			AggregateName: "projection",
			EventName:     "created",
			Data:          `{"name":"projection1"}`,
		},
	}, nil).Once()

	estore.On("StoreEvent", "proj1", "projection", "block_added", data).Return(nil).Once()

	err := es.AddBlock(ctx, "proj1", data)
	assert.NoError(t, err)

	estore.AssertCalled(t, "StoreEvent", "proj1", "projection", "block_added", data)

}
