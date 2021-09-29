package items

type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	Id uint64 `json:"id"`
}
