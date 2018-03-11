package requests

// PasteCreateRequest represents a request to paste.create.
type PasteCreateRequest struct {
	Content  string `json:"content"`  // required
	Title    string `json:"title"`    // optional
	Language string `json:"language"` // optional
	Request
}
