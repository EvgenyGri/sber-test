package processors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert" //nolint:goimports
	"sber-test/pkg/models"
	ts "sber-test/pkg/time-slot" //nolint:goimports
)

func TestReportIfMatchedRecipes(t *testing.T) {
	wantedName := []string{"Potato", "Veggie", "Mushroom"}
	data := []models.RecipeDelivery{
		{Recipe: "Ink"},
		{Recipe: "B Potato"},
		{Recipe: "A Veggie"},
		{Recipe: "C Mushroom"},
	}
	expected := []string{"A Veggie", "B Potato", "C Mushroom"}
	rep := ReportIfMatchedRecipes(wantedName[0], wantedName[1:]...)
	for _, item := range data {
		rep.consume(item)
	}
	var report RecipeProcessorReport
	rep.fillReport(&report)
	assert.NotNil(t, report.RecipesMatchedByName)
	assert.Equal(t, expected, report.RecipesMatchedByName)
}

func TestReportUniqueRecipes(t *testing.T) {
	data := []models.RecipeDelivery{
		{Recipe: "Ink"},
		{Recipe: "Ink"},
		{Recipe: "B Potato"},
		{Recipe: "A Veggie"},
		{Recipe: "C Mushroom"},
		{Recipe: "C Mushroom"},
	}
	rep := ReportUniqueRecipes()
	for _, item := range data {
		rep.consume(item)
	}
	var report RecipeProcessorReport
	rep.fillReport(&report)
	assert.NotNil(t, report.UniqueRecipeCount)
	assert.Equal(t, 4, *report.UniqueRecipeCount)
}

func TestReportBusiestPostcode(t *testing.T) {
	data := []models.RecipeDelivery{ //nolint:dupl
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Monday, 10, 15)},
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Thursday, 10, 15)},
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Thursday, 18, 22)},
		{Recipe: "B Potato", Postcode: "2", Delivery: ts.ConstructDelivery(time.Wednesday, 8, 15)},
		{Recipe: "A Veggie", Postcode: "3", Delivery: ts.ConstructDelivery(time.Saturday, 11, 20)},
		{Recipe: "C Mushroom", Postcode: "4", Delivery: ts.ConstructDelivery(time.Saturday, 11, 20)},
		{Recipe: "C Mushroom", Postcode: "5", Delivery: ts.ConstructDelivery(time.Friday, 11, 22)},
	}

	rep := ReportBusiestPostcode()
	for _, item := range data {
		rep.consume(item)
	}
	var report RecipeProcessorReport
	rep.fillReport(&report)
	assert.NotNil(t, report.BusiestPostcode)
	assert.Equal(t, "1", report.BusiestPostcode.Postcode)
	assert.Equal(t, 3, report.BusiestPostcode.DeliveryCount)
}

func TestReportCounterPerRecipe(t *testing.T) {
	data := []models.RecipeDelivery{ //nolint:dupl
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Monday, 10, 15)},
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Thursday, 10, 15)},
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Thursday, 18, 22)},
		{Recipe: "B Potato", Postcode: "2", Delivery: ts.ConstructDelivery(time.Wednesday, 8, 15)},
		{Recipe: "A Veggie", Postcode: "3", Delivery: ts.ConstructDelivery(time.Saturday, 11, 20)},
		{Recipe: "C Mushroom", Postcode: "4", Delivery: ts.ConstructDelivery(time.Saturday, 11, 20)},
		{Recipe: "C Mushroom", Postcode: "5", Delivery: ts.ConstructDelivery(time.Friday, 11, 22)},
	}
	rep := ReportCounterPerRecipe()
	for _, item := range data {
		rep.consume(item)
	}
	var report RecipeProcessorReport
	rep.fillReport(&report)
	assert.NotNil(t, report.CountPerRecipe)
	expected := []countPerRecipe{
		{Recipe: "A Veggie", Count: 1},
		{Recipe: "B Potato", Count: 1},
		{Recipe: "C Mushroom", Count: 2},
		{Recipe: "Ink", Count: 3},
	}
	assert.Equal(t, expected, report.CountPerRecipe)
}

func TestReportDeliveryCountForPostcodeAndTime(t *testing.T) {
	data := []models.RecipeDelivery{ //nolint:dupl
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Monday, 10, 15)},
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Thursday, 10, 15)},
		{Recipe: "Ink", Postcode: "1", Delivery: ts.ConstructDelivery(time.Thursday, 18, 22)},
		{Recipe: "B Potato", Postcode: "2", Delivery: ts.ConstructDelivery(time.Wednesday, 8, 15)},
		{Recipe: "A Veggie", Postcode: "3", Delivery: ts.ConstructDelivery(time.Saturday, 11, 20)},
		{Recipe: "C Mushroom", Postcode: "4", Delivery: ts.ConstructDelivery(time.Saturday, 11, 20)},
		{Recipe: "C Mushroom", Postcode: "5", Delivery: ts.ConstructDelivery(time.Friday, 11, 22)},
	}
	rep := ReportDeliveryCountForPostcodeAndTime("1", ts.Hour(9), ts.Hour(19))
	for _, item := range data {
		rep.consume(item)
	}
	var report RecipeProcessorReport
	rep.fillReport(&report)
	assert.NotNil(t, report.CountPerPostcodeAndTime)
	assert.Equal(t, "1", report.CountPerPostcodeAndTime.Postcode)
	assert.Equal(t, 2, report.CountPerPostcodeAndTime.DeliveryCount)
}
