package event

type Data struct {
	ID                 string
	ReferenceID        string
	Type               string
	UserCreated        *UserCreated
	ApplicationCreated *ApplicationCreated
}

type Repository interface {
	Store(data Data)
}

func Store(timeline Repository, data Data) (err error) {
	return
}
