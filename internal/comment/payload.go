package comment

type CommentPostRequest struct {
	Text string `json:"text" validate:"required"`
}

type CommentPostUpdateRequest struct {
	Text string `json:"text" validate:"required"`
}