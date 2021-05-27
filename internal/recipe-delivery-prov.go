package internal

import (
	"context"
	"os"

	jsoniter "github.com/json-iterator/go" //nolint:goimports
	"github.com/pkg/errors"
	"sber-test/pkg/models"
	"sber-test/pkg/providers" //nolint:goimports
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// NewRecipeDeliveryProviderFromFile ...
func NewRecipeDeliveryProviderFromFile(f string) providers.RecipeDeliveryProvider {
	return &recipeDeliveryProvider{
		source: f,
	}
}

// RecipeDeliveryProvider ...
type recipeDeliveryProvider struct {
	source string
}

// Provide ...
func (p *recipeDeliveryProvider) Provide(_ context.Context, consumer func(models.RecipeDelivery) error) error {
	const api = "RecipeDeliveryProvider.Provide"

	f, e := os.Open(p.source)
	if e != nil {
		return errors.Wrapf(e, "%s: open file('%s')", api, p.source)
	}

	defer f.Close() //nolint:gosec
	decoder := json.NewDecoder(f)
	var raw []models.RecipeDelivery
	if e = decoder.Decode(&raw); e != nil {
		return errors.Wrapf(e, "%s: JSON Decode", api)
	}
	for _, item := range raw {
		if e := consumer(item); e != nil {
			return e
		}
	}
	return nil
}
