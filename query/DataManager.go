package entity

type EventManager struct {
	ID       string
	Name     string
	Property map[string]string
}

type CommandManager struct {
	ID       string
	Name     string
	Property map[string]string
}

type ProjectionManager struct {
	ID            string
	Name          string
	EntityManager []EntityManager
}

type EntityManager struct {
	ID       string
	Name     string
	Property map[string]string
}
