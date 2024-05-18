package notes

import (
	"sort"
)

func (notes Notes) toResp() NotesResp {
	var resps = []NoteResp{}
	for _, note := range notes {
		var resp NoteResp = note.toResp()
		resps = append(resps, resp)
	}
	sort.Slice(resps, func(i, j int) bool {
		return resps[i].Date.Sub(resps[j].Date) > 0
	})
	return NotesResp{Notes: resps}
}

func (note Note) toResp() NoteResp {
	var resp = NoteResp{
		ID:   note.ID,
		Date: note.Date,
		Text: note.Text,
	}
	return resp
}
