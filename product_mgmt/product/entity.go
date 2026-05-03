package product

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Apps []App  `json:"apps"`
}

type App struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
