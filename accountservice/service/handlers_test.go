package service

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spacycoder/go-blog/accountservice/dbclient"
	"github.com/spacycoder/go-blog/accountservice/model"
)

func init() {
	gock.InterceptClient(client)
}

func TestGetAccountWrongPath(t *testing.T) {
	Convey("Given a HTTP request for /invalid/123", t, func() {
		req := httptest.NewRequest("GET", "/invalid/123", nil)
		resp := httptest.NewRecorder()
		Convey("When the request is handled by the Router", func() {
			NewRouter().ServeHTTP(resp, req)
			Convey("Then the response should be a 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
		})
	})
}

func TestGetAccount(t *testing.T) {
	defer gock.Off()
	gock.New("http://quotes-service:8080").
		Get("/api/quote").
		MatchParam("strength", "4").
		Reply(200).
		BodyString(`{"quote":"May the source be with you. Always.","ipAddress":"10.0.0.5:8080","language":"en"}`)
	// Create a mock instance that implements the IBoltClient interface
	mockRepo := &dbclient.MockBoltClient{}
	// Declare two mock behaviours. For "123" as input, return a proper Account struct and nil as error.
	// For "456" as input, return an empty Account object and a real error.
	mockRepo.On("QueryAccount", "123").Return(model.Account{ID: "123", Name: "Person_123"}, nil)
	mockRepo.On("QueryAccount", "456").Return(model.Account{}, fmt.Errorf("Some error"))
	// Finally, assign mockRepo to the DBClient field (it's in _handlers.go_, e.g. in the same package)
	DBClient = mockRepo

	Convey("Given a HTTP request for /accounts/123", t, func() {
		req := httptest.NewRequest("GET", "/accounts/123", nil)
		resp := httptest.NewRecorder()
		Convey("When the request is handled by the Router", func() {
			NewRouter().ServeHTTP(resp, req)
			Convey("Then the response should be a 200", func() {
				So(resp.Code, ShouldEqual, 200)
				account := model.Account{}
				json.Unmarshal(resp.Body.Bytes(), &account)
				So(account.ID, ShouldEqual, "123")
				So(account.Name, ShouldEqual, "Person_123")
				So(account.Quote.Text, ShouldEqual, "May the source be with you. Always.")

			})
		})
	})

	Convey("Given a HTTP request for /accounts/456", t, func() {
		req := httptest.NewRequest("GET", "/accounts/456", nil)
		resp := httptest.NewRecorder()
		Convey("When the request is handled by the Router", func() {
			NewRouter().ServeHTTP(resp, req)
			Convey("Then the response should be a 404", func() {
				So(resp.Code, ShouldEqual, 404)
			})
		})
	})
}
