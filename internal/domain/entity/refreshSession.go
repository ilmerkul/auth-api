package entity

type RefreshSession struct {
	ID           int
	UserId       int
	RefreshToken string
	UA           string
	Fingerprint  string
	ExpiresIn    string
	CreatedAt    string
}
