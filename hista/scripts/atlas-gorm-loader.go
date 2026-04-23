package main

import (
	"fmt"
	"io"
	"os"

	// _ "ariga.io/atlas-go-sdk/recordriver"
	gormSchema "ariga.io/atlas-provider-gorm/gormschema"
	"encore.app/hista/internal/headaches"
	"encore.app/hista/internal/meals"
	"encore.app/hista/internal/medicines"
	"encore.app/hista/internal/notes"
	"encore.app/hista/internal/pollen"
	"encore.app/hista/internal/status"
	"encore.app/hista/internal/symptoms"
)

// Define the models to generate migrations for.
var models = []any{
	&meals.Food{},
	&meals.Ingredient{},
	&meals.Meal{},
	&meals.Template{},
	&meals.TemplateItem{},
	&symptoms.Condition{},
	&symptoms.SymptomCategory{},
	&symptoms.Symptom{},
	&symptoms.ConditionEvent{},
	&notes.Note{},
	&pollen.Pollen{},
	&pollen.PollenEvent{},
	&status.Status{},
	&headaches.Headache{},
	&medicines.Intake{},
	&medicines.Medicine{},
}

func main() {
	statements, err := gormSchema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, statements)
}
