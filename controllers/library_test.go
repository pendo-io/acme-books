package controllers_test

import (
	"acme-books/controllers"
	"acme-books/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestTuple struct {
	name string
	run  func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock)
}

type FDatastoreQueryMock struct {
	mock.Mock
}

func (d FDatastoreQueryMock) Create(book models.Book) error {
	returnValues := d.Called(book)
	err := returnValues.Error(0)
	return err

}

func (d FDatastoreQueryMock) Change(id int, book models.Book) error {
	returnValues := d.Called(id, book)
	err := returnValues.Error(0)
	return err
}

func (d FDatastoreQueryMock) GetById(id int) (models.Book, error) {
	returnValues := d.Called(id)
	err := returnValues.Error(1)

	if err == nil {
		return returnValues.Get(0).(models.Book), nil
	} else {
		return models.Book{}, err
	}
}

func (d FDatastoreQueryMock) List(query *datastore.Query) ([]models.Book, error) {
	returnValues := d.Called(query)
	err := returnValues.Error(1)

	if err == nil {
		return returnValues.Get(0).([]models.Book), nil
	} else {
		return []models.Book{}, err
	}
}

func TestExample(t *testing.T) {
	params := make(map[string]string)
	w := httptest.NewRecorder()

	fixture := FDatastoreQueryMock{}
	fixture.
		On("List", datastore.NewQuery("Book").Order("Id")).
		Return([]models.Book{}, nil).
		Once()

	c := new(controllers.LibraryController)
	c.ListAll(params, w, &fixture)

	require.Equal(t, "[]", w.Body.String())
	require.Equal(t, http.StatusOK, w.Code)
	fixture.AssertExpectations(t)
}

func TestList(t *testing.T) {
	var tests []TestTuple
	defer func() {
		for _, test := range tests {
			w := httptest.NewRecorder()
			e := FDatastoreQueryMock{}
			t.Run(test.name, func(t *testing.T) {
				test.run(t, w, &e)
			})
		}
	}()

	tests = []TestTuple{
		{"no data", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			f.
				On("List", datastore.NewQuery("Book").Order("Id")).
				Return([]models.Book{}, nil).
				Once()

			c := new(controllers.LibraryController)
			c.ListAll(params, w, f)

			require.Equal(t, "[]", w.Body.String())
			require.Equal(t, http.StatusOK, w.Code)
			f.AssertExpectations(t)
		}},
		{"filtering", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			params["title"] = "cats"
			f.
				On("List", datastore.NewQuery("Book").Order("Id").Filter("Title =", "cats")).
				Return([]models.Book{}, nil).
				Once()

			c := new(controllers.LibraryController)
			c.ListAll(params, w, f)

			require.Equal(t, "[]", w.Body.String())
			require.Equal(t, http.StatusOK, w.Code)
			f.AssertExpectations(t)
		}},
		{"filtering an invalid value", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			params["dogs"] = "asdfsd"
			f.
				On("List", datastore.NewQuery("Book").Order("Id")).
				Return([]models.Book{}, nil).
				Once()

			c := new(controllers.LibraryController)
			c.ListAll(params, w, f)

			require.Equal(t, "[]", w.Body.String())
			require.Equal(t, http.StatusOK, w.Code)
			f.AssertExpectations(t)
		}},
	}
}

func TestBorrow(t *testing.T) {
	var tests []TestTuple
	defer func() {
		for _, test := range tests {
			w := httptest.NewRecorder()
			e := FDatastoreQueryMock{}
			t.Run(test.name, func(t *testing.T) {
				test.run(t, w, &e)
			})
		}
	}()

	tests = []TestTuple{
		{"no 'id' param", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)

			c := new(controllers.LibraryController)
			c.Borrow(params, w, f)

			require.Equal(t, "", w.Body.String())
			require.Equal(t, http.StatusBadRequest, w.Code)
			f.AssertExpectations(t)
		}},

		{"invalid 'id' param", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			params["id"] = "a"
			c := new(controllers.LibraryController)
			c.Borrow(params, w, f)

			require.Equal(t, "", w.Body.String())
			require.Equal(t, http.StatusBadRequest, w.Code)
			f.AssertExpectations(t)
		}},

		{"book already borrowed", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			params["id"] = "1"

			f.On("GetById", 1).
				Return(models.Book{
					Id:       0,
					Title:    "hello",
					Borrowed: true,
					Author:   "cat",
				}, nil).
				Once()

			c := new(controllers.LibraryController)
			c.Borrow(params, w, f)

			require.Equal(t, "", w.Body.String())
			require.Equal(t, http.StatusBadRequest, w.Code)
			f.AssertExpectations(t)
		}},

		{"success", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			params["id"] = "1"
			book := models.Book{
				Id:       0,
				Title:    "hello",
				Borrowed: false,
				Author:   "cat",
			}
			f.On("GetById", 1).
				Return(book, nil).
				Once()

			book.Borrowed = true
			f.On("Change", 1, book).
				Return(nil).
				Once()

			c := new(controllers.LibraryController)
			c.Borrow(params, w, f)

			require.Equal(t, "", w.Body.String())
			require.Equal(t, http.StatusNoContent, w.Code)
			f.AssertExpectations(t)
		}},
	}
}

func TestReturn(t *testing.T) {
	var tests []TestTuple
	defer func() {
		for _, test := range tests {
			w := httptest.NewRecorder()
			e := FDatastoreQueryMock{}
			t.Run(test.name, func(t *testing.T) {
				test.run(t, w, &e)
			})
		}
	}()

	tests = []TestTuple{
		{"no 'id' param", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)

			c := new(controllers.LibraryController)
			c.Return(params, w, f)

			require.Equal(t, "", w.Body.String())
			require.Equal(t, http.StatusBadRequest, w.Code)
			f.AssertExpectations(t)
		}},

		{"invalid 'id' param", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			params["id"] = "a"
			c := new(controllers.LibraryController)
			c.Return(params, w, f)

			require.Equal(t, "", w.Body.String())
			require.Equal(t, http.StatusBadRequest, w.Code)
			f.AssertExpectations(t)
		}},

		{"book not yet borrowed", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			params["id"] = "1"

			f.On("GetById", 1).
				Return(models.Book{
					Id:       0,
					Title:    "hello",
					Borrowed: false,
					Author:   "cat",
				}, nil).
				Once()

			c := new(controllers.LibraryController)
			c.Return(params, w, f)

			require.Equal(t, "", w.Body.String())
			require.Equal(t, http.StatusBadRequest, w.Code)
			f.AssertExpectations(t)
		}},

		{"success", func(t *testing.T, w *httptest.ResponseRecorder, f *FDatastoreQueryMock) {
			params := make(map[string]string)
			params["id"] = "1"
			book := models.Book{
				Id:       0,
				Title:    "hello",
				Borrowed: true,
				Author:   "cat",
			}
			f.On("GetById", 1).
				Return(book, nil).
				Once()

			book.Borrowed = false
			f.On("Change", 1, book).
				Return(nil).
				Once()

			c := new(controllers.LibraryController)
			c.Return(params, w, f)

			require.Equal(t, "", w.Body.String())
			require.Equal(t, http.StatusNoContent, w.Code)
			f.AssertExpectations(t)
		}},
	}
}
