package api

import (
	"acme-books/model"
	"bytes"
	"cloud.google.com/go/datastore"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type bookImplementationMock struct {
	mock.Mock
}

func (bi *bookImplementationMock) GetByKey(id int) (model.Book, error) {
	args := bi.Called(id)
	return args.Get(0).(model.Book), args.Error(1)
}

func (bi *bookImplementationMock) ListAll(query *datastore.Query) []model.Book {
	args := bi.Called(query)
	return args.Get(0).([]model.Book)
}

func (bi *bookImplementationMock) ChangeBorrowedStatus(id int, borrowed bool) error {
	args := bi.Called(id, borrowed)
	return args.Error(1)
}

func (bi *bookImplementationMock) AddBook(book model.Book) (model.Book, error) {
	args := bi.Called(book)
	return args.Get(0).(model.Book), args.Error(1)
}

func (bi *bookImplementationMock) Delete(id int) error {
	args := bi.Called(id)
	return args.Error(0)
}

type LibraryTestSuite struct {
	suite.Suite
	w       *httptest.ResponseRecorder
	bi      bookImplementationMock
	library Library
}

func (suite *LibraryTestSuite) SetupTest() {
	suite.w = httptest.NewRecorder()
	suite.bi = bookImplementationMock{}
	suite.library = Library{}
}

func (suite *LibraryTestSuite) TestGetByKey() {
	params := make(map[string]string)
	params["id"] = "1"
	id, _ := strconv.Atoi(params["id"])

	suite.bi.On("GetByKey", id).Return(model.Book{}, nil)
	suite.library.GetByKey(params, suite.w, &suite.bi)

	suite.Assert().Equal(http.StatusOK, suite.w.Code)
	suite.bi.AssertExpectations(suite.T())
}

func (suite *LibraryTestSuite) TestListAll() {
	r := httptest.NewRequest("GET", "http://localhost:3030/books", nil)
	query := datastore.NewQuery("Book")

	suite.bi.On("ListAll", query).Return([]model.Book{}, nil)
	suite.library.ListAll(r, suite.w, &suite.bi)

	suite.Assert().Equal(http.StatusOK, suite.w.Code)
	suite.bi.AssertExpectations(suite.T())
}

func (suite *LibraryTestSuite) TestBorrow() {
	params := make(map[string]string)
	params["id"] = "1"
	id, _ := strconv.Atoi(params["id"])

	suite.bi.On("ChangeBorrowedStatus", id, true).Return(model.Book{}, nil)
	suite.library.Borrow(params, suite.w, &suite.bi)

	suite.Assert().Equal(http.StatusNoContent, suite.w.Code)
	suite.bi.AssertExpectations(suite.T())
}

func (suite *LibraryTestSuite) TestReturn() {
	params := make(map[string]string)
	params["id"] = "1"
	id, _ := strconv.Atoi(params["id"])

	suite.bi.On("ChangeBorrowedStatus", id, true).Return(model.Book{}, nil)
	suite.library.Borrow(params, suite.w, &suite.bi)

	suite.Assert().Equal(http.StatusNoContent, suite.w.Code)
	suite.bi.AssertExpectations(suite.T())
}

func (suite *LibraryTestSuite) TestAddBook() {
	book := model.Book{Title: "Test Book", Author: "Test Author", Borrowed: false}
	buffer := new(bytes.Buffer)
	buffer.Write([]byte("{\n  \"title\": \"Test Book\",\n  \"writer\": \"Test Author\",\n  \"borrowed\": false\n}"))

	r := httptest.NewRequest("POST", "http://localhost:3030/books", buffer)

	suite.bi.On("AddBook", book).Return(model.Book{}, nil)
	suite.library.AddBook(r, suite.w, &suite.bi)

	suite.Assert().Equal(http.StatusOK, suite.w.Code)
	suite.bi.AssertExpectations(suite.T())
}

func (suite *LibraryTestSuite) TestDeleteBook() {
	params := make(map[string]string)
	params["id"] = "1"
	id, _ := strconv.Atoi(params["id"])

	suite.bi.On("Delete", id).Return(nil)
	suite.library.DeleteBook(params, suite.w, &suite.bi)

	suite.Assert().Equal(http.StatusNoContent, suite.w.Code)
	suite.bi.AssertExpectations(suite.T())
}

func TestLibrarySuite(t *testing.T) {
	suite.Run(t, new(LibraryTestSuite))
}
