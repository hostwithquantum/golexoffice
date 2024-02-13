package golexoffice_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hostwithquantum/golexoffice"
	"github.com/stretchr/testify/assert"
)

// {"requestId":"3fb21ee4-ad26-4e2f-82af-a1197af02d08","IssueList":[{"i18nKey":"invalid_value","source":"company and person","type":"validation_failure"},{"i18nKey":"missing_entity","source":"company.name","type":"validation_failure"}]}

// {"requestId":"75d4dad6-6ccb-40fd-8c22-797f2d421d98","IssueList":[{"i18nKey":"missing_entity","source":"company.vatRegistrationId","type":"validation_failure"},{"i18nKey":"missing_entity","source":"company.taxNumber","type":"validation_failure"}]}

func TestErrorResponse(t *testing.T) {
	server := errorMock()
	defer server.Close()

	lexOffice := golexoffice.NewConfig("token", nil)
	lexOffice.SetBaseUrl(server.URL)

	t.Run("errors=legacy", func(t *testing.T) {
		_, err := lexOffice.AddContact(golexoffice.ContactBody{
			Company: &golexoffice.ContactBodyCompany{
				Name:              "company",
				VatRegistrationId: "",
				TaxNumber:         "",
			},
		})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "key: missing_entity (company.vatRegistrationId): validation_failure")
		assert.ErrorContains(t, err, "key: missing_entity (company.taxNumber): validation_failure")
	})

	t.Run("errors=new", func(t *testing.T) {
		_, err := lexOffice.AddInvoice(golexoffice.InvoiceBody{})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "field: lineItems[0].unitPrice.taxRatePercentage (NOTNULL): darf nicht leer sein")
	})

}

func TestErrorNoDetails(t *testing.T) {
	server := errorMockNoDetails()
	defer server.Close()

	lexOffice := golexoffice.NewConfig("token", nil)
	lexOffice.SetBaseUrl(server.URL)

	t.Run("errors=legacy", func(t *testing.T) {
		_, err := lexOffice.AddContact(golexoffice.ContactBody{
			Company: &golexoffice.ContactBodyCompany{
				Name:              "company",
				VatRegistrationId: "",
				TaxNumber:         "",
			},
		})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "something went wrong")
	})

	t.Run("errors=new", func(t *testing.T) {
		_, err := lexOffice.AddInvoice(golexoffice.InvoiceBody{})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Something else went wrong. (406 Not Acceptable)")
	})
}

func TestRateLimit(t *testing.T) {
	rateLimitHits := 10
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rateLimitHits > 0 {
			rateLimitHits -= 1
			w.WriteHeader(http.StatusTooManyRequests)
			//nolint:errcheck
			w.Write([]byte(`{
				"status": 429,
				"error": "Too Many Requests",
				"message": "Rate limit exceeded"
			}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`)) //nolint:errcheck
		}
	}))
	defer server.Close()

	lexOffice := golexoffice.NewConfig("token", nil)
	lexOffice.SetBaseUrl(server.URL)

	t.Run("retry until ok", func(t *testing.T) {
		rateLimitHits = 2
		_, err := lexOffice.Invoice("tralalala")
		assert.NoError(t, err)
	})

	t.Run("retry until out of retries", func(t *testing.T) {
		rateLimitHits = 10
		_, err := lexOffice.Invoice("tralalala")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Rate limit exceeded")
	})
}

func errorMock() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/contacts/" {
			w.WriteHeader(http.StatusBadRequest)
			//nolint:errcheck
			w.Write([]byte(`{
				"requestId":"75d4dad6-6ccb-40fd-8c22-797f2d421d98",
				"IssueList":[
					{"i18nKey":"missing_entity","source":"company.vatRegistrationId","type":"validation_failure"},
					{"i18nKey":"missing_entity","source":"company.taxNumber","type":"validation_failure"}
				]
			}`))
			return
		}
		if r.URL.Path == "/v1/invoices" {
			w.WriteHeader(http.StatusNotAcceptable)
			//nolint:errcheck
			w.Write([]byte(`{
				"timestamp": "2017-05-11T17:12:31.233+02:00",
				"status": 406,
				"error": "Not Acceptable",
				"path": "/v1/invoices",
				"traceId": "90d78d0777be",
				"message": "Validation failed for request. Please see details list for specific causes.",
				"details": [
					{
						"violation": "NOTNULL",
						"field": "lineItems[0].unitPrice.taxRatePercentage",
						"message": "darf nicht leer sein"
					}
				]
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func errorMockNoDetails() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/contacts/" {
			w.WriteHeader(http.StatusBadRequest)
			// unclear if this can actually happen for these legacy errors
			//nolint:errcheck
			w.Write([]byte(`{
				"requestId":"75d4dad6-6ccb-40fd-8c22-797f2d421d98",
				"IssueList":[]
			}`))
			return
		}
		if r.URL.Path == "/v1/invoices" {
			w.WriteHeader(http.StatusNotAcceptable)
			// this *can* and *does* happen however (details is optional)
			//nolint:errcheck
			w.Write([]byte(`{
				"timestamp": "2017-05-11T17:12:31.233+02:00",
				"status": 406,
				"error": "Not Acceptable",
				"path": "/v1/invoices",
				"traceId": "90d78d0777be",
				"message": "Something else went wrong."
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}
