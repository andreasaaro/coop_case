package mastodon

type Account struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type MastodonData struct {
	ID        string  `json:"id,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	Content   string  `json:"content,omitempty"`
	Account   Account `json:"account,omitempty"`
}
