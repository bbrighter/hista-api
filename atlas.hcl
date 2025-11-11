data "external_schema" "users_gorm" {
	program = ["env", "ENCORERUNTIME_NOPANIC=1", "go", "run", "./users/scripts/atlas-gorm-loader.go"]
}
  
env "users" {
	src = data.external_schema.users_gorm.url
	url = "postgresql://hista-api-dpc2:local@127.0.0.1:9500/users_db?sslmode=disable"
  
	migration {
		dir = "file://users/migrations"
	  	format = golang-migrate
	}

	format {
	  	migrate {
			diff = "{{ sql . \"  \" }}"
	  	}
	}
}

data "external_schema" "product_mgmt_gorm" {
	program = ["env", "ENCORERUNTIME_NOPANIC=1", "go", "run", "./product_mgmt/scripts/atlas-gorm-loader.go"]
}
  
env "product_mgmt" {
	src = data.external_schema.product_mgmt_gorm.url
	url = "postgresql://hista-api-dpc2:local@127.0.0.1:9500/product_mgmt_db?sslmode=disable"
  
	migration {
		dir = "file://product_mgmt/migrations"
	  	format = golang-migrate
	}

	format {
	  	migrate {
			diff = "{{ sql . \"  \" }}"
	  	}
	}
}

data "external_schema" "api_gorm" {
	program = ["env", "ENCORERUNTIME_NOPANIC=1", "go", "run", "./api/scripts/atlas-gorm-loader.go"]
}
env "api" {
	src = data.external_schema.api_gorm.url
	url = "postgresql://hista-api-dpc2:local@127.0.0.1:9500/hista_db?sslmode=disable"
  
	migration {
	  	dir = "file://api/migrations"
	  	format = golang-migrate
	}
  
	format {
	  	migrate {
			diff = "{{ sql . \"  \" }}"
	  	}
	}
}
  
data "external_schema" "shopping_list_gorm" {
	program = ["env", "ENCORERUNTIME_NOPANIC=1", "go", "run", "./shoppingList/scripts/atlas-gorm-loader.go"]
}
  
env "shopping_list" {
	src = data.external_schema.shopping_list_gorm.url
	url = "postgresql://hista-api-dpc2:local@127.0.0.1:9500/shopping_list?sslmode=disable"
  
	migration {
		dir = "file://shoppingList/migrations"
	  	format = golang-migrate
	}

	format {
	  	migrate {
			diff = "{{ sql . \"  \" }}"
	  	}
	}
}