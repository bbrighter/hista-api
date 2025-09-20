package repository

import (
	"log"
	"os"

	"encore.app/product_mgmt/entity"
	"github.com/goccy/go-yaml"
)

type Config struct {
	Products []entity.Product
	Apps     []entity.App
}

func ParseConfig(path string) Config {
	type yamlProduct struct {
		ID     string   `yaml:"id"`
		Name   string   `yaml:"name"`
		AppIds []string `yaml:"apps"`
	}
	type config struct {
		Apps     []entity.App  `yaml:"apps"`
		Products []yamlProduct `yaml:"products"`
	}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Cannot read path: %s", err)
	}
	var c config
	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Fatalf("Cannot parse yaml: %s", err)
	}

	appMap := make(map[string]entity.App)
	for _, app := range c.Apps {
		appMap[app.ID] = app
	}

	var products []entity.Product
	for _, p := range c.Products {
		var apps []entity.App
		for _, appId := range p.AppIds {
			apps = append(apps, appMap[appId])
		}
		products = append(products, entity.Product{
			ID:   p.ID,
			Name: p.Name,
			Apps: apps,
		})
	}
	return Config{Products: products, Apps: c.Apps}
}
