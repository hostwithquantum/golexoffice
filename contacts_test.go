package golexoffice_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hostwithquantum/golexoffice"
	"github.com/stretchr/testify/assert"
)

func TestAddContact(t *testing.T) {
	server := lexOfficeMock()
	defer server.Close()

	config := golexoffice.NewConfig("api-key", nil)
	config.SetBaseUrl(server.URL)

	resp, err := config.AddContact(golexoffice.ContactBody{
		Roles: golexoffice.ContactBodyRoles{
			Customer: golexoffice.ContactBodyCustomer{},
		},
		Person: golexoffice.ContactBodyPerson{
			FirstName: "Thomas",
			LastName:  "Mustermann",
		},
		Note: "golexoffice",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.NotEmpty(t, resp.CreatedDate)
	assert.NotEmpty(t, resp.UpdatedDate)
	assert.NotEmpty(t, resp.ResourceUri)
}

func TestGetContacts(t *testing.T) {
	server := lexOfficeMock()
	defer server.Close()

	config := golexoffice.NewConfig("api-key", nil)
	config.SetBaseUrl(server.URL)

	t.Run("mock=company", func(t *testing.T) {
		resp, err := config.Contact("c73d5f78-847e-49d8-aa58-c6d95c5c9cb5")
		assert.NoError(t, err)

		assert.Equal(t, "c73d5f78-847e-49d8-aa58-c6d95c5c9cb5", resp.Id)
		assert.Equal(t, 10001, resp.Roles.Customer.Number)
		assert.Equal(t, 70003, resp.Roles.Vendor.Number)
		assert.Equal(t, "Beispiel GmbH", resp.Company.Name)
	})

	t.Run("mock=person", func(t *testing.T) {
		resp, err := config.Contact("e9066f04-8cc7-4616-93f8-ac9ecc8479c8")
		assert.NoError(t, err)

		assert.Equal(t, "e9066f04-8cc7-4616-93f8-ac9ecc8479c8", resp.Id)
		assert.Equal(t, 10308, resp.Roles.Customer.Number)
		assert.Equal(t, "Frau", resp.Person.Salutation)
		assert.Equal(t, "Inge", resp.Person.FirstName)
		assert.Equal(t, "Musterfrau", resp.Person.LastName)
	})
}

func lexOfficeMock() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if r.URL.Path == "/v1/contacts/c73d5f78-847e-49d8-aa58-c6d95c5c9cb5" {
				//nolint:errcheck
				w.Write([]byte(`{
					"id":"c73d5f78-847e-49d8-aa58-c6d95c5c9cb5",
					"organizationId":"67c8c57b-6d07-4bdd-b579-55240d3c2df5",
					"version":1,
					"roles":{
						"customer":{"number":10001},
						"vendor":{"number":70003}
					},
					"company":{
						"name":"Beispiel GmbH",
						"vatRegistrationId":"DE123456789",
						"contactPersons":[{
							"salutation":"Herr",
							"firstName":"Thomas",
							"lastName":"Mustermann",
							"primary":true,
							"emailAddress": "thomas@example.org"
						}]},
						"addresses":{
							"billing":[{
								"street":"Straße 1",
								"zip":"10111",
								"city":"Berlin",
								"countryCode":"DE"
							}]
						},
						"archived":false}`))
				return
			}

			if r.URL.Path == "/v1/contacts/e9066f04-8cc7-4616-93f8-ac9ecc8479c8" {
				//nolint:errcheck
				w.Write([]byte(`{
					"id": "e9066f04-8cc7-4616-93f8-ac9ecc8479c8",
					"organizationId": "aa93e8a8-2aa3-470b-b914-caad8a255dd8",
					"version": 0,
					"roles": {
					  "customer": {
						"number": 10308
					  }
					},
					"person": {
					  "salutation": "Frau",
					  "firstName": "Inge",
					  "lastName": "Musterfrau"
					},
					"note": "Notizen",
					"archived": false
				}`))
				return
			}
		}
		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			if r.URL.Path == "/v1/contacts/" {
				//nolint:errcheck
				w.Write([]byte(`{
					"id": "66196c43-baf3-4335-bfee-d610367059db",
					"resourceUri": "https://api.lexoffice.io/v1/contacts/66196c43-bfee-baf3-4335-d610367059db",
					"createdDate": "2016-06-29T15:15:09.447+02:00",
					"updatedDate": "2016-06-29T15:15:09.447+02:00",
					"version": 1
				}`))
				return
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		//nolint:errcheck
		w.Write([]byte(fmt.Sprintf("not found (%s): %s", r.Method, r.RequestURI)))
	}))
}
