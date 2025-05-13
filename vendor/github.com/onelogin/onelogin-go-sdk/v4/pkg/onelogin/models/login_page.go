package models

type CreateSessionLoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
	Subdomain       string `json:"subdomain"`
	Fields          string `json:"fields,omitempty"`
}

type VerifyFactorRequest struct {
	DeviceID    string `json:"device_id"`
	StateToken  string `json:"state_token"`
	OtpToken    string `json:"otp_token,omitempty"`
	DoNotNotify bool   `json:"do_not_notify,omitempty"`
}
