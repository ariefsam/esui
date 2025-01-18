package eventstore

type Eventstore struct{}

func NewEventstore() (obj *Eventstore) {
	obj = &Eventstore{}
	return obj
}
