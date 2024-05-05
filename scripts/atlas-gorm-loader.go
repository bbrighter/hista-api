package main

import (
	"fmt"
	"io"
	"os"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
	"encore.app/meals"
	"encore.app/states"
)

// Define the models to generate migrations for.
var models = []any{
	&meals.Food{},
	&meals.Ingredient{},
	&meals.Meal{},
	&states.Condition{},
	&states.SymptomCategory{},
	&states.Symptom{},
	&states.State{},
}

func main() {
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
