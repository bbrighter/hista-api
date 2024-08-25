package main

import (
	"fmt"
	"io"
	"os"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
	"encore.app/entity"
	"encore.app/internalAuth"
	"encore.app/notes"
	"encore.app/pollen"
	"encore.app/symptoms"
)

// Define the models to generate migrations for.
var models = []any{
	&entity.Food{},
	&entity.Ingredient{},
	&entity.Meal{},
	&symptoms.Condition{},
	&symptoms.SymptomCategory{},
	&symptoms.Symptom{},
	&symptoms.ConditionEvent{},
	&internalAuth.Token{},
	&internalAuth.User{},
	&notes.Note{},
	&pollen.Pollen{},
	&pollen.PollenEvent{},
}

func main() {
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
