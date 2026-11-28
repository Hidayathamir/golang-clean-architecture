package dto

// SendVerificationEmailRequest is the input to the email sender. It is
// deliberately provider-agnostic: the SMTP implementation maps it to a message.
type SendVerificationEmailRequest struct {
	To    string
	Name  string
	Token string
}
