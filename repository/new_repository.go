package repository

type Repository struct{}

func NewRepository() (obj *Repository) {
	obj = &Repository{}
	return obj
}
