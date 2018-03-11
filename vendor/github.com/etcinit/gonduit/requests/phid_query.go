package requests

type PHIDQueryRequest struct {
	PHIDs []string `json:"phids"`
	Request
}
