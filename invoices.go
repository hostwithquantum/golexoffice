//**********************************************************
//
// Copyright (C) 2018 - 2023 J&J Ideenschmiede GmbH <info@jj-ideenschmiede.de>
//
// This file is part of golexoffice.
// All code may be used. Feel free and maybe code something better.
//
// Author: Jonas Kwiedor
//
//**********************************************************

package golexoffice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// InvoiceBody is to define body data
type InvoiceBody struct {
	Id                 string                        `json:"id,omitempty"`
	OrganizationId     string                        `json:"organizationId,omitempty"`
	CreateDate         string                        `json:"createDate,omitempty"`
	UpdatedDate        string                        `json:"updatedDate,omitempty"`
	Version            int                           `json:"version,omitempty"`
	Archived           bool                          `json:"archived,omitempty"`
	VoucherStatus      string                        `json:"voucherStatus,omitempty"`
	VoucherNumber      string                        `json:"voucherNumber,omitempty"`
	VoucherDate        string                        `json:"voucherDate"`
	DueDate            interface{}                   `json:"dueDate,omitempty"`
	Address            InvoiceBodyAddress            `json:"address"`
	LineItems          []InvoiceBodyLineItems        `json:"lineItems"`
	TotalPrice         InvoiceBodyTotalPrice         `json:"totalPrice"`
	TaxAmounts         []InvoiceBodyTaxAmounts       `json:"taxAmounts,omitempty"`
	TaxConditions      InvoiceBodyTaxConditions      `json:"taxConditions"`
	PaymentConditions  *InvoiceBodyPaymentConditions `json:"paymentConditions,omitempty"`
	ShippingConditions InvoiceBodyShippingConditions `json:"shippingConditions"`
	Title              string                        `json:"title,omitempty"`
	Introduction       string                        `json:"introduction,omitempty"`
	Remark             string                        `json:"remark,omitempty"`
	Language           string                        `json:"language,omitempty"`
}

type InvoiceBodyAddress struct {
	ContactId   string `json:"contactId,omitempty"`
	Name        string `json:"name,omitempty"`
	Supplement  string `json:"supplement,omitempty"`
	Street      string `json:"street,omitempty"`
	City        string `json:"city,omitempty"`
	Zip         string `json:"zip,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

type InvoiceBodyLineItems struct {
	Id                 string               `json:"id,omitempty"`
	Type               string               `json:"type"`
	Name               string               `json:"name"`
	Description        string               `json:"description,omitempty"`
	Quantity           interface{}          `json:"quantity,omitempty"`
	UnitName           string               `json:"unitName,omitempty"`
	UnitPrice          InvoiceBodyUnitPrice `json:"unitPrice,omitempty"`
	DiscountPercentage interface{}          `json:"discountPercentage,omitempty"`
	LineItemAmount     interface{}          `json:"lineItemAmount,omitempty"`
}

type InvoiceBodyUnitPrice struct {
	Currency          string      `json:"currency"`
	NetAmount         interface{} `json:"netAmount,omitempty"`
	GrossAmount       interface{} `json:"grossAmount,omitempty"`
	TaxRatePercentage int         `json:"taxRatePercentage"`
}

type InvoiceBodyTotalPrice struct {
	Currency                string      `json:"currency"`
	TotalNetAmount          interface{} `json:"totalNetAmount,omitempty"`
	TotalGrossAmount        interface{} `json:"totalGrossAmount,omitempty"`
	TaxRatePercentage       interface{} `json:"taxRatePercentage,omitempty"`
	TotalTaxAmount          interface{} `json:"totalTaxAmount,omitempty"`
	TotalDiscountAbsolute   interface{} `json:"totalDiscountAbsolute,omitempty"`
	TotalDiscountPercentage interface{} `json:"totalDiscountPercentage,omitempty"`
}

type InvoiceBodyTaxAmounts struct {
	TaxRatePercentage int     `json:"taxRatePercentage"`
	TaxAmount         float64 `json:"taxAmount"`
	Amount            float64 `json:"amount"`
}

type InvoiceBodyTaxConditions struct {
	TaxType     string      `json:"taxType"`
	TaxTypeNote interface{} `json:"taxTypeNote,omitempty"`
}

type InvoiceBodyPaymentConditions struct {
	PaymentTermLabel          string                               `json:"paymentTermLabel"`
	PaymentTermDuration       int                                  `json:"paymentTermDuration"`
	PaymentDiscountConditions InvoiceBodyPaymentDiscountConditions `json:"paymentDiscountConditions"`
}

type InvoiceBodyPaymentDiscountConditions struct {
	DiscountPercentage int `json:"discountPercentage"`
	DiscountRange      int `json:"discountRange"`
}

type InvoiceBodyShippingConditions struct {
	ShippingDate    string      `json:"shippingDate,omitempty"`
	ShippingEndDate interface{} `json:"shippingEndDate,omitempty"`
	ShippingType    string      `json:"shippingType"`
}

// InvoiceReturn is to decode json data
type InvoiceReturn struct {
	Id          string `json:"id"`
	ResourceUri string `json:"resourceUri"`
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
	Version     int    `json:"version"`
}

// Invoice is to get a invoice by id
func (c *Config) Invoice(id string) (InvoiceBody, error) {

	// Set config for new request
	//c := NewConfig(, token, &http.Client{})

	// Send request
	response, err := c.Send("/v1/invoices/"+id, nil, "GET", "application/json")
	if err != nil {
		return InvoiceBody{}, err
	}

	// Close request
	defer response.Body.Close()

	read, err := io.ReadAll(response.Body)
	if err != nil {
		return InvoiceBody{}, err
	}
	fmt.Println(string(read))

	// Decode data
	var decode InvoiceBody

	err = json.NewDecoder(response.Body).Decode(&decode)
	if err != nil {
		return InvoiceBody{}, err
	}

	// Return data
	return decode, nil

}

// AddInvoice is to create a invoice
func (c *Config) AddInvoice(body InvoiceBody) (InvoiceReturn, error) {

	// NOTE: we're using VoucherStatus ("open" or "draft") to determine if this
	// should be a draft invoice
	//
	// The way the API works is this: (https://developers.lexoffice.io/docs/#invoices-endpoint-create-an-invoice)
	// > Invoices transmitted via the API are created in draft mode per default. To
	// > create a finalized invoice with status open the optional query parameter
	// > finalize has to be set. The status of an invoice cannot be changed via the api.
	isOpen := body.VoucherStatus == "open"
	body.VoucherStatus = "" // unset for the request

	// Convert body
	convert, err := json.Marshal(body)
	if err != nil {
		return InvoiceReturn{}, err
	}

	// Set config for new request
	//c := NewConfig(, token, &http.Client{})

	// Send request
	url := "/v1/invoices"
	if isOpen {
		url += "?finalize=true"
	}
	response, err := c.Send(url, bytes.NewBuffer(convert), "POST", "application/json")
	if err != nil {
		return InvoiceReturn{}, err
	}

	// Close request
	defer response.Body.Close()

	// Decode data
	var decode InvoiceReturn

	err = json.NewDecoder(response.Body).Decode(&decode)
	if err != nil {
		return InvoiceReturn{}, err
	}

	// Return data
	return decode, nil

}
