package product

import (
	"encore.dev/config"
)

type Config struct {
	Products config.Values[ProductConfig]
	Apps     config.Values[AppConfig]
}

type ProductConfig struct {
	ID   config.String
	Name config.String
	Apps []config.String
}

type AppConfig struct {
	ID   config.String
	Name config.String
}

func (cfg *Config) GetApps(prodCfg ProductConfig) []App {
	var apps []App
	for _, a := range prodCfg.Apps {
		apps = append(apps, cfg.FindAppById(a()))
	}

	return apps
}

func (cfg *Config) ToProducts() []Product {
	var products []Product
	for _, prod := range cfg.Products() {
		products = append(products, Product{
			ID:   prod.ID(),
			Name: prod.Name(),
			Apps: cfg.GetApps(prod),
		})
	}
	return products
}

func (cfg *Config) ToApps() []App {
	var apps []App
	for _, app := range cfg.Apps() {
		apps = append(apps, App{
			ID:   app.ID(),
			Name: app.Name(),
		})
	}
	return apps
}

func (cfg *Config) FindAppById(appId string) App {
	apps := cfg.ToApps()
	for _, app := range apps {
		if app.ID == appId {
			return app
		}
	}
	return App{}
}
