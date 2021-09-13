package event

type UserCreated struct {
	ID         string
	Username   string
	Email      string
	IsSysAdmin bool
}
