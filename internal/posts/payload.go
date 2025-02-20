package posts

type PostCreateRequest struct {
	Title string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	ImageURL string `json:"image_url" validate:"url"`
	Published bool `json:"published" validate:"required"`
}
  
type PostCreateDB struct {
	Title string 
	Content string 
	ImageURL string 
	Published bool
	UserID uint
}

type PostUploadRequest struct {
	Title string `json:"title"`
	Content string `json:"content"`
	ImageURL string `json:"image_url,omitempty" validate:"omitempty,url"`
}