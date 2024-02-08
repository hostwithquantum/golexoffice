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
	Id                 string                        `json:"id"`
	OrganizationId     string                        `json:"organizationId"`
	CreateDate         string                        `json:"createDate"`
	UpdatedDate        string                        `json:"updatedDate"`
	Version            int                           `json:"version"`
	Archived           bool                          `json:"archived"`
	VoucherStatus      string                        `json:"voucherStatus"`
	VoucherNumber      string                        `json:"voucherNumber"`
	VoucherDate        string                        `json:"voucherDate"`
	DueDate            interface{}                   `json:"dueDate"`
	Address            InvoiceBodyAddress            `json:"address"`
	LineItems          []InvoiceBodyLineItems        `json:"lineItems"`
	TotalPrice         InvoiceBodyTotalPrice         `json:"totalPrice"`
	TaxAmounts         []InvoiceBodyTaxAmounts       `json:"taxAmounts"`
	TaxConditions      InvoiceBodyTaxConditions      `json:"taxConditions"`
	PaymentConditions  InvoiceBodyPaymentConditions  `json:"paymentConditions"`
	ShippingConditions InvoiceBodyShippingConditions `json:"shippingConditions"`
	Title              string                        `json:"title"`
	Introduction       string                        `json:"introduction"`
	Remark             string                        `json:"remark"`
	Language           string                        `json:"language"`
}

type InvoiceBodyAddress struct {
	ContactId   string `json:"contactId"`
	Name        string `json:"name"`
	Supplement  string `json:"supplement"`
	Street      string `json:"street"`
	City        string `json:"city"`
	Zip         string `json:"zip"`
	CountryCode string `json:"countryCode"`
}

type InvoiceBodyLineItems struct {
	Id                 string               `json:"id,omitempty"`
	Type               string               `json:"type"`
	Name               string               `json:"name"`
	Description        string               `json:"description"`
	Quantity           float64              `json:"quantity"`
	UnitName           string               `json:"unitName"`
	UnitPrice          InvoiceBodyUnitPrice `json:"unitPrice"`
	DiscountPercentage int                  `json:"discountPercentage"`
	LineItemAmount     float64              `json:"lineItemAmount"`
}

type InvoiceBodyUnitPrice struct {
	Currency          string  `json:"currency"`
	NetAmount         float64 `json:"netAmount"`
	GrossAmount       float64 `json:"grossAmount"`
	TaxRatePercentage int     `json:"taxRatePercentage"`
}

type InvoiceBodyTotalPrice struct {
	Currency                string      `json:"currency"`
	TotalNetAmount          float64     `json:"totalNetAmount"`
	TotalGrossAmount        float64     `json:"totalGrossAmount"`
	TaxRatePercentage       interface{} `json:"taxRatePercentage"`
	TotalTaxAmount          float64     `json:"totalTaxAmount"`
	TotalDiscountAbsolute   interface{} `json:"totalDiscountAbsolute"`
	TotalDiscountPercentage interface{} `json:"totalDiscountPercentage"`
}

type InvoiceBodyTaxAmounts struct {
	TaxRatePercentage int     `json:"taxRatePercentage"`
	TaxAmount         float64 `json:"taxAmount"`
	Amount            float64 `json:"amount"`
}

type InvoiceBodyTaxConditions struct {
	TaxType     string      `json:"taxType"`
	TaxTypeNote interface{} `json:"taxTypeNote"`
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
	ShippingDate    string      `json:"shippingDate"`
	ShippingEndDate interface{} `json:"shippingEndDate"`
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

	// Convert body
	convert, err := json.Marshal(body)
	if err != nil {
		return InvoiceReturn{}, err
	}

	// Set config for new request
	//c := NewConfig(, token, &http.Client{})

	// Send request
	response, err := c.Send("/v1/invoices", bytes.NewBuffer(convert), "POST", "application/json")
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
