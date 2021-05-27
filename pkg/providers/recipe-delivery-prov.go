package providers

import (
	"context"

	"sber-test/pkg/models"
)

// RecipeDeliveryProvider ...
type RecipeDeliveryProvider interface {
	Provide(context.Context, func(models.RecipeDelivery) error) error
}
