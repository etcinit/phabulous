package requests

type PHIDLookupRequest struct {
	Names []string `json:"names"`
	Request
}
