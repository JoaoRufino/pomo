package mainrpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/mocks"
	"github.com/snowzach/gorestapi/store"
)

func TestTaskPost(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	grs := new(mocks.GRStore)
	err := Setup(r, grs)
	assert.Nil(t, err)

	// Create Item
	i := &gorestapi.Task{
		ID:   "id",
		Name: "name",
	}

	// Mock call to item store
	grs.On("TaskSave", mock.Anytask, i).Once().Return(nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.POST("/tasks").WithJSON(i).Expect().Status(http.StatusOK).JSON().Object().Equal(i)

	// Check remaining expectations
	grs.AssertExpectations(t)

}

func TestTasksFind(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	grs := new(mocks.GRStore)
	err := Setup(r, grs)
	assert.Nil(t, err)

	// Return Item
	i := []*gorestapi.Task{
		{
			ID:   "id1",
			Name: "name1",
		},
		{
			ID:   "id2",
			Name: "name2",
		},
	}

	// Mock call to item store
	grs.On("TasksFind", mock.Anytask, mock.AnythingOfType("*queryp.QueryParameters")).Once().Return(i, int64(2), nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/tasks").Expect().Status(http.StatusOK).JSON().Object().Equal(&store.Results{Count: 2, Results: i})

	// Check remaining expectations
	grs.AssertExpectations(t)

}

func TestTaskGetByID(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	grs := new(mocks.GRStore)
	err := Setup(r, grs)
	assert.Nil(t, err)

	// Create Item
	i := &gorestapi.Task{
		ID:   "id",
		Name: "name",
	}

	// Mock call to item store
	grs.On("TaskGetByID", mock.Anytask, "1234").Once().Return(i, nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/tasks/1234").Expect().Status(http.StatusOK).JSON().Object().Equal(&i)

	// Check remaining expectations
	grs.AssertExpectations(t)

}

func TestTaskDeleteByID(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	grs := new(mocks.GRStore)
	err := Setup(r, grs)
	assert.Nil(t, err)

	// Mock call to item store
	grs.On("TaskDeleteByID", mock.Anytask, "1234").Once().Return(nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.DELETE("/tasks/1234").Expect().Status(http.StatusNoContent)

	// Check remaining expectations
	grs.AssertExpectations(t)

}
