package product

type ProductService struct {
	p *ProductRepo
}

func NewProductService(p *ProductRepo) *ProductService {
	return &ProductService{p: p}
}

func (s *ProductService) Find(id string) (Product, error) {
	return s.p.Find(id)
}
