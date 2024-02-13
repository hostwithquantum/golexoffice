package golexoffice_test

import (
	"encoding/json"
	"testing"

	"github.com/hostwithquantum/golexoffice"
	"github.com/stretchr/testify/assert"
)

func TestEmptyInvoiceOnlyHasRequiredProperties(t *testing.T) {
	body := golexoffice.InvoiceBody{}
	encoded, err := json.Marshal(body)
	assert.NoError(t, err)
	assert.JSONEq(t, `
    {
      "voucherDate": "",
      "address": {},
      "lineItems": null,
      "totalPrice": {
        "currency": ""
      },
      "taxConditions": {
        "taxType": ""
      },
      "shippingConditions": {
        "shippingType": ""
      }
    }`, string(encoded))
}

func TestPriceOfZeroIsNotOmitted(t *testing.T) {
	body := golexoffice.InvoiceBody{
		TotalPrice: golexoffice.InvoiceBodyTotalPrice{
			TotalGrossAmount: 0,
		},
		LineItems: []golexoffice.InvoiceBodyLineItems{
			{
				UnitPrice: golexoffice.InvoiceBodyUnitPrice{
					NetAmount: 0.0,
				},
			},
			{
				UnitPrice: golexoffice.InvoiceBodyUnitPrice{
					GrossAmount: 0.0,
				},
			}},
	}
	encoded, err := json.Marshal(body)
	assert.NoError(t, err)
	assert.JSONEq(t, `
    {
      "voucherDate": "",
      "address": {},
      "lineItems": [
          {
            "name": "",
            "type": "",
            "unitPrice": {
              "currency": "",
              "taxRatePercentage": 0,
              "netAmount": 0
            }
          },
          {
            "name": "",
            "type": "",
            "unitPrice": {
              "currency": "",
              "taxRatePercentage": 0,
              "grossAmount": 0
            }
          }
      ],
      "totalPrice": {
        "currency": "",
        "totalGrossAmount": 0
      },
      "taxConditions": {
        "taxType": ""
      },
      "shippingConditions": {
        "shippingType": ""
      }
    }`, string(encoded))
}
