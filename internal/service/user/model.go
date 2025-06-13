package user

type User struct {
	Phone             string
	PersonExternalRef string
	Token             string
	Language          string
	Status            string
	UserID            int
	CountryID         int8
}

const (
	_deleted = "deleted"
)
