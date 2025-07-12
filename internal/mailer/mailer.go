package mailer

import "embed"

const (
	FromName            = "wot"
	maxRetries          = 3
	UserWelcomeTemplate = "invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
