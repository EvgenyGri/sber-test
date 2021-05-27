package models

import (
	ts "sber-test/pkg/time-slot"
)

// RecipeDelivery ...
type RecipeDelivery struct {
	Postcode string      `json:"postcode"`
	Recipe   string      `json:"recipe"`
	Delivery ts.Delivery `json:"delivery"`
}
