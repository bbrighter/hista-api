package hista

import (
	"time"

	"encore.dev/beta/errs"
	"encore.dev/types/option"
)

func (s *ApiTestSuite) TestGetNotes() {
	tests := map[string]struct {
		createBefore   bool
		expectedNumber int
	}{
		"0": {},
		"1": {createBefore: true, expectedNumber: 1},
	}
	for name, test := range tests {
		s.Run(name, func() {
			if test.createBefore {
				s.createTestNote()
			}
			resp, err := s.service.ListNotes(s.ctx, s.piid)
			s.NoError(err)
			s.Len(resp.Notes, test.expectedNumber)
		})
	}
}

func (s *ApiTestSuite) TestDeleteNote() {

	var err error
	err = s.service.DeleteNote(s.ctx, s.piid, 1)
	s.Error(err)

	noteId := s.createTestNote()

	err = s.service.DeleteNote(s.ctx, s.piid, noteId)
	s.NoError(err)
}

func (s *ApiTestSuite) TestPatchNote() {
	date := time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local)
	text := "text"
	tests := map[string]struct {
		useWrongId        bool
		params            NoteParams
		expectedErrorCode errs.ErrCode
	}{
		"ok, date":  {params: NoteParams{Date: option.Some(date)}},
		"ok, text":  {params: NoteParams{Text: option.Some(text)}},
		"not found": {useWrongId: true, params: NoteParams{Text: option.Some(text)}, expectedErrorCode: errs.NotFound},
		"no params": {expectedErrorCode: errs.InvalidArgument},
	}

	for name, test := range tests {
		s.Run(name, func() {
			noteId := s.createTestNote()
			if test.useWrongId {
				noteId = 1000
			}
			err := s.service.PatchNote(s.ctx, s.piid, noteId, test.params)
			s.assertErrCode(err, test.expectedErrorCode)
		})
	}
}
