package shoppinglist

func (s *ApiTestSuite) TestListProducts() {
	tests := map[string]struct {
		productExists bool
		useWrongPiid  bool
		expectedLen   int
	}{
		"1 entry":    {productExists: true, expectedLen: 1},
		"0 entries":  {},
		"wrong piid": {useWrongPiid: true, expectedLen: 0},
	}
	for name, test := range tests {
		s.Run(name, func() {
			if test.productExists {
				s.createProduct()
			}
			ctx := s.GetCtx(test.useWrongPiid)

			prods, err := s.service.ListProducts(ctx, s.piid)
			s.NoError(err)
			s.Len(prods.Products, test.expectedLen)
		})
	}
}
