package tests

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect/v2"
	"github.com/phone_book/internal/lib"
	apiResp "github.com/phone_book/internal/lib/api/response"
	"github.com/phone_book/internal/store"
)

const (
	host = "localhost:1234"
)

func generatePersonsData() []store.Person {
	countPerson := lib.RandRange(2, 12)
	persons := make([]store.Person, 0, countPerson)
	for range countPerson {
		phoneInt, _ := strconv.Atoi(gofakeit.Phone())
		person := store.InitPersonEntry(gofakeit.Name(), gofakeit.LastName(), phoneInt)
		persons = append(persons, *person)

	}
	return persons
}

func TestPhoneBookInsertSearchRemoveList(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())
	phone := gofakeit.Phone()
	phoneInt, _ := strconv.Atoi(phone)
	insertPerson := *store.InitPersonEntry(gofakeit.Name(), gofakeit.LastName(), phoneInt)
	personsList := append(make([]interface{}, 0), insertPerson)
	e.POST("/insert").
		WithJSON(insertPerson).
		Expect().
		Status(200).
		JSON().Object().IsEqual(apiResp.OK("new record added successfully"))
	numberNotExist := gofakeit.Phone()
	numberNotInt := gofakeit.Name()
	e.GET("/search/" + numberNotInt).Expect().
		Status(200).
		JSON().Object().IsEqual(apiResp.Error(fmt.Sprintf("Phone number %s is incorrect", numberNotInt)))
	e.GET("/search/" + numberNotExist).Expect().
		Status(200).
		JSON().Object().Value("data").IsNull()
	e.GET("/search/" + phone).Expect().
		Status(200).
		JSON().Object().IsEqual(apiResp.WithData(insertPerson))
	e.GET("/search/"+phone).WithQuery("start_with", "1").Expect().
		Status(200).
		JSON().Object().IsEqual(apiResp.WithDataList(personsList))

	e.DELETE("/remove/" + numberNotInt).Expect().
		Status(200).
		JSON().Object().IsEqual(apiResp.Error(fmt.Sprintf("Phone number %s is incorrect", numberNotInt)))

	e.DELETE("/remove/" + numberNotExist).Expect().
		Status(200).
		JSON().Object().IsEqual(apiResp.Error(fmt.Sprintf("Person with number %s not found", numberNotExist)))

	e.DELETE("/remove/" + phone).Expect().
		Status(200).
		JSON().Object().IsEqual(apiResp.OK(fmt.Sprintf("Record with number %s deleted", phone)))

}
