package shoppinglist

func (s *ApiTestSuite) TestGetMoments() {
	tests := map[string]struct {
		useOutdatedETag bool
	}{
		"ok":           {},
		"not modified": {useOutdatedETag: true},
	}

	for name, test := range tests {
		s.Run(name, func() {
			listId := s.createList()
			s.createItem(listId)
			etag := s.etag
			if test.useOutdatedETag {
				etag = "123"
			}
			resp, err := s.service.GetMoments(s.ctx, s.piid, MomentsParams{IfNoneMatch: etag})
			s.NoError(err)
			if test.useOutdatedETag {
				s.Equal(resp.Status, 200)
				s.Len(resp.Items, 1)
				s.Equal(listId, resp.ListId)
				s.Len(resp.Products, 1)
				s.NotEqual(resp.ETag, "0")
			} else {
				s.Equal(resp.Status, 304)
				s.Len(resp.Items, 0)
				s.EqualValues(0, resp.ListId)
			}

		})
	}
}
