package notes

func (notes Notes) toResp() NotesResp {
	var resps []NoteResp
	for _, note := range notes {
		var resp NoteResp = note.toResp()
		resps = append(resps, resp)
	}
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
