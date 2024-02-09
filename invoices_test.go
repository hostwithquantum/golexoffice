package golexoffice_test

import (
	"testing"
    "encoding/json"

	"github.com/hostwithquantum/golexoffice"
	"github.com/stretchr/testify/assert"
)

func TestEmptyInvoiceOnlyHasRequiredProperties(t *testing.T) {
    body := golexoffice.InvoiceBody{}
	encoded, err := json.Marshal(body)
    assert.NoError(t, err)
    assert.Equal(t, `{"voucherDate":"","address":{},"lineItems":null,"totalPrice":{"currency":""},"taxConditions":{"taxType":""},"shippingConditions":{"shippingType":""}}`, string(encoded))
}
