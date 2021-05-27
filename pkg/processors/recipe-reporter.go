package processors

import (
	"context"
	"sort"
	"strings"

	"github.com/pkg/errors" //nolint:goimports
	"sber-test/pkg/models"
	"sber-test/pkg/providers"
	ts "sber-test/pkg/time-slot" //nolint:goimports
)

/*//
{
    "unique_recipe_count": 15,
    "count_per_recipe": [
        {
            "recipe": "Mediterranean Baked Veggies",
            "count": 1
        },
        {
            "recipe": "Speedy Steak Fajitas",
            "count": 1
        },
        {
            "recipe": "Tex-Mex Tilapia",
            "count": 3
        }
    ],
    "busiest_postcode": {
        "postcode": "10120",
        "delivery_count": 1000
    },
    "count_per_postcode_and_time": {
        "postcode": "10120",
        "from": "11AM",
        "to": "3PM",
        "delivery_count": 500
    },
    "match_by_name": [
        "Mediterranean Baked Veggies", "Speedy Steak Fajitas", "Tex-Mex Tilapia"
    ]
}
*/

type (
	countPerRecipe struct {
		Recipe string `json:"recipe"`
		Count  int    `json:"count"`
	}

	busiestPostcode struct {
		Postcode      string `json:"postcode"`
		DeliveryCount int    `json:"delivery_count"`
	}

	countPerPostcodeAndTime struct {
		Postcode      string  `json:"postcode"`
		From          ts.Hour `json:"from"`
		To            ts.Hour `json:"to"`
		DeliveryCount int     `json:"delivery_count"`
	}

	// RecipeProcessorReport отчёт
	RecipeProcessorReport struct {
		UniqueRecipeCount       *int                     `json:"unique_recipe_count,omitempty"`
		CountPerRecipe          []countPerRecipe         `json:"count_per_recipe,omitempty"`
		BusiestPostcode         *busiestPostcode         `json:"busiest_postcode,omitempty"`
		CountPerPostcodeAndTime *countPerPostcodeAndTime `json:"count_per_postcode_and_time,omitempty"`
		RecipesMatchedByName    []string                 `json:"match_by_name,omitempty"`
	}

	// RecipeReportProcessor тот кто нам отчёт сделает
	RecipeReportProcessor struct {
		reporters []RecipeReportSubj
	}

	// RecipeReportSubj ,,,
	RecipeReportSubj interface {
		consume(models.RecipeDelivery)
		fillReport(*RecipeProcessorReport)
	}
)

// ReportUniqueRecipes Подсчитать число уникальных "recipe name"
func ReportUniqueRecipes() RecipeReportSubj {
	return &uniqueRecipeCounter{counter: make(map[string]struct{})}
}

// ReportBusiestPostcode Найти "postcode" с наибольшим числом доаставок
func ReportBusiestPostcode() RecipeReportSubj {
	return &busiestPostcodeReporter{
		postalCodeCounter: make(map[string]map[ts.Delivery]struct{}),
	}
}

//ReportIfMatchedRecipes Перечислить "recipe name" (в алфавитном порядке), которые содержат в своём имени одно из следующих слов
func ReportIfMatchedRecipes(name string, optional ...string) RecipeReportSubj {
	return &recipeMatchByName{
		names: append(append([]string(nil), name), optional...),
		res:   make(map[string]struct{}),
	}
}

//ReportCounterPerRecipe подсчитать число вхождений каждого уникального "recipe name" (с алфавитной сортировкой по "recipe name")
func ReportCounterPerRecipe() RecipeReportSubj {
	return &counterPerRecipe{
		counter: make(map[string]int),
	}
}

//ReportDeliveryCountForPostcodeAndTime Найти число доставок для "postcode", которые происходили во временном промежутке
func ReportDeliveryCountForPostcodeAndTime(postCode string, from, to ts.Hour) RecipeReportSubj {
	ret := new(counterPerPostcodeAndTime)
	ret.Postcode = postCode
	ret.From, ret.To = from, to
	return ret
}

// NewRecipeReportProcessor Recipe репортер
func NewRecipeReportProcessor(r RecipeReportSubj, optional ...RecipeReportSubj) *RecipeReportProcessor {
	ret := new(RecipeReportProcessor)
	ret.reporters = append(append(ret.reporters, r), optional...)
	return ret
}

// Process обработаеи и получим-ка отчётец
func (rp *RecipeReportProcessor) Process(ctx context.Context, provider providers.RecipeDeliveryProvider) (RecipeProcessorReport, error) {
	const api = "RecipeProcessor.Process"

	var report RecipeProcessorReport
	err := provider.Provide(ctx, func(delivery models.RecipeDelivery) error {
		for _, rep := range rp.reporters {
			rep.consume(delivery)
		}
		return nil
	})
	if err != nil {
		return report, errors.Wrap(err, api)
	}
	for _, rep := range rp.reporters {
		rep.fillReport(&report)
	}
	return report, nil
}

// ---------------------------------------- IMPL -------------------------------------

type uniqueRecipeCounter struct {
	counter map[string]struct{}
}

func (r *uniqueRecipeCounter) consume(item models.RecipeDelivery) {
	r.counter[item.Recipe] = struct{}{}
}

func (r *uniqueRecipeCounter) fillReport(rep *RecipeProcessorReport) {
	n := len(r.counter)
	rep.UniqueRecipeCount = &n
}

type recipeMatchByName struct {
	names []string
	res   map[string]struct{}
}

func (r *recipeMatchByName) consume(item models.RecipeDelivery) {
	for i := range r.names {
		if strings.Contains(item.Recipe, r.names[i]) {
			r.res[item.Recipe] = struct{}{}
			return
		}
	}
}

func (r *recipeMatchByName) fillReport(rep *RecipeProcessorReport) {
	res := make([]string, 0, len(r.res))
	for s := range r.res {
		res = append(res, s)
	}
	if len(res) > 0 {
		sort.Slice(res, func(i, j int) bool {
			l, r := res[i], res[j]
			return strings.Compare(l, r) < 0
		})
		rep.RecipesMatchedByName = res
	}
}

type counterPerRecipe struct {
	counter map[string]int
}

func (r *counterPerRecipe) consume(item models.RecipeDelivery) {
	r.counter[item.Recipe] = r.counter[item.Recipe] + 1
}

func (r *counterPerRecipe) fillReport(rep *RecipeProcessorReport) {
	items := make([]countPerRecipe, 0, len(r.counter))
	for name, c := range r.counter {
		items = append(items, countPerRecipe{Recipe: name, Count: c})
	}
	if len(items) > 0 {
		sort.Slice(items, func(i, j int) bool {
			l, r := items[i], items[j]
			return strings.Compare(l.Recipe, r.Recipe) < 0
		})
		rep.CountPerRecipe = items
	}
}

type busiestPostcodeReporter struct {
	postalCodeCounter map[string]map[ts.Delivery]struct{}
}

func (r *busiestPostcodeReporter) consume(item models.RecipeDelivery) {
	counter := r.postalCodeCounter[item.Postcode]
	if counter == nil {
		counter = make(map[ts.Delivery]struct{})
		r.postalCodeCounter[item.Postcode] = counter
	}
	counter[item.Delivery] = struct{}{}
}

func (r *busiestPostcodeReporter) fillReport(rep *RecipeProcessorReport) {
	if len(r.postalCodeCounter) == 0 {
		return
	}
	var busiest busiestPostcode
	first := true
	for p, c := range r.postalCodeCounter {
		if first {
			busiest.Postcode = p
			busiest.DeliveryCount = len(c)
			first = false
		} else if busiest.DeliveryCount < len(c) {
			busiest.Postcode = p
			busiest.DeliveryCount = len(c)
		}
	}
	if len(busiest.Postcode) > 0 {
		rep.BusiestPostcode = &busiest
	}
}

type counterPerPostcodeAndTime struct {
	countPerPostcodeAndTime
}

func (r *counterPerPostcodeAndTime) consume(item models.RecipeDelivery) {
	matched := item.Postcode == r.Postcode &&
		r.From <= item.Delivery.From &&
		item.Delivery.To <= r.To
	if matched {
		r.DeliveryCount++
	}
}

func (r *counterPerPostcodeAndTime) fillReport(rep *RecipeProcessorReport) {
	if r.DeliveryCount > 0 {
		rep.CountPerPostcodeAndTime = &r.countPerPostcodeAndTime
	}
}
