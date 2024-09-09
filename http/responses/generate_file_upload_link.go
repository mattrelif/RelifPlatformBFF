package responses

type GenerateFileUploadLink struct {
	Link string `json:"link"`
}

func NewGenerateFileUploadLink(link string) GenerateFileUploadLink {
	return GenerateFileUploadLink{
		Link: link,
	}
}
