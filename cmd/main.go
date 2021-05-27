package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"sber-test/internal"
	"sber-test/pkg/processors"
	ts "sber-test/pkg/time-slot"
)

var (
	source                            string
	reportCountPerRecipe              bool
	reportUniqueRecipeCount           bool
	reportBusiestPostcode             bool
	reportMatchedRecipe               string
	reportDeliveriesByPostcodeAndTime string
)

func init() {
	flag.StringVar(&source, "source", "", "points fo source file needs in processing")
	flag.BoolVar(&reportCountPerRecipe, "count-per-recipe", false, "reports counts per Recipe")
	flag.BoolVar(&reportUniqueRecipeCount, "unique-recipe-count", false, "reports unique Recipe count")
	flag.BoolVar(&reportBusiestPostcode, "busiest-postcode", false, "report busiest postcode")
	flag.StringVar(&reportMatchedRecipe, "find-recipes", "", "report recipes by name(s); example: --find-recipes='Potato,Veggie.Mushroom'")
	flag.StringVar(&reportDeliveriesByPostcodeAndTime,
		"deliveries-by-postcode-and-time", "",
		"deliveries by postcode and time; example: --deliveries-by-postcode-and-time='10120,10AM,3PM'")
}

func reportError(formats string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, formats, args...)
}

func reportSubjectsFromArgs() []processors.RecipeReportSubj {
	var subjects []processors.RecipeReportSubj
	if reportCountPerRecipe {
		subjects = append(subjects, processors.ReportCounterPerRecipe())
	}
	if reportUniqueRecipeCount {
		subjects = append(subjects, processors.ReportUniqueRecipes())
	}
	if reportBusiestPostcode {
		subjects = append(subjects, processors.ReportBusiestPostcode())
	}
	if len(reportMatchedRecipe) > 0 {
		var names []string
		for _, name := range strings.Split(reportMatchedRecipe, ",") {
			if name = strings.TrimSpace(name); len(name) > 0 {
				names = append(names, name)
			}
		}
		if len(names) == 0 {
			reportError("'--find-recipes' param has wrong value")
			os.Exit(1)
		}
		subjects = append(subjects, processors.ReportIfMatchedRecipes(names[0], names[1:]...))
	}
	if len(reportDeliveriesByPostcodeAndTime) > 0 {
		raw := strings.Split(reportDeliveriesByPostcodeAndTime, ",")
		const p = "--deliveries-by-postcode-and-time"
		if len(raw) != 3 {
			reportError("'%s' param has wrong value", p)
			os.Exit(1)
		}
		var from, to ts.Hour
		err := from.FromString([]byte(raw[1]))
		if err == nil {
			err = to.FromString([]byte(raw[2]))
		}
		if err != nil {
			reportError("'%s' param has wrong value cause %v", p, err)
			os.Exit(1)
		}
		subjects = append(subjects, processors.ReportDeliveryCountForPostcodeAndTime(raw[0], from, to))
	}
	return subjects
}

func main() {
	flag.Parse()
	if len(source) == 0 {
		reportError("source param is not provided")
		os.Exit(1)
	}
	subjects := reportSubjectsFromArgs()
	if len(subjects) == 0 {
		reportError("asked no any subject to report")
		os.Exit(1)
	}
	reporter := processors.NewRecipeReportProcessor(subjects[0], subjects[1:]...)
	src := internal.NewRecipeDeliveryProviderFromFile(source)
	ctx := context.Background()
	report, err := reporter.Process(ctx, src)
	if err != nil {
		reportError("%v", err)
		os.Exit(1)
	}
	var decodedResult []byte
	if decodedResult, err = json.Marshal(report); err != nil {
		reportError("%v", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", string(decodedResult))
}
