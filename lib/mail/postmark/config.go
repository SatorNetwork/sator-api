package postmark

import "github.com/keighl/postmark"

type (
	// Config struct
	Config struct {
		ProductName    string
		ProductURL     string
		SupportURL     string
		SupportEmail   string
		CompanyName    string
		CompanyAddress string
		FromEmail      string
		FromName       string
	}

	postmarkClient interface {
		SendTemplatedEmail(email postmark.TemplatedEmail) (postmark.EmailResponse, error)
	}
)
