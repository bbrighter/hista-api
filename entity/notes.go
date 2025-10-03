package entity

import (
	"sort"
	"time"

	"encore.dev/types/uuid"
)

type Note struct {
	ID   uint
	Date time.Time
	Text string
	PIID uuid.UUID `gorm:"type:uuid;index"`
}

func (n *Note) SetPiid(id uuid.UUID) {
	n.PIID = id
}

type Notes []*Note

type NotesResp struct {
	Notes []NoteResp `json:"notes"`
}

type NoteResp struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
	Text string    `json:"text"`
}

func (notes Notes) ToResp() NotesResp {
	var resps = []NoteResp{}
	for _, note := range notes {
		var resp NoteResp = note.ToResp()
		resps = append(resps, resp)
	}
	sort.Slice(resps, func(i, j int) bool {
		return resps[i].Date.Sub(resps[j].Date) > 0
	})
	return NotesResp{Notes: resps}
}

func (note Note) ToResp() NoteResp {
	return NoteResp{
		ID:   note.ID,
		Date: note.Date,
		Text: note.Text,
	}
}
