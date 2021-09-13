package command

import "github.com/ariefsam/esui/event"

type Event interface {
	Store(data event.Data) (err error)
}
