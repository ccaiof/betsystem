package auth

type Identity struct {
	Subject    string
	ProviderID string
	Roles      []string
}
