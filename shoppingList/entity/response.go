package entity

type IdResponse struct {
	ID uint `json:"id"`
}

func ToIdResponse(id uint) IdResponse {
	return IdResponse{ID: id}
}
