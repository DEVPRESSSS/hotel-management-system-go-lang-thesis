package dto

type PriceRequest struct {
	RoomID   string `json:"room_id"`
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
}
