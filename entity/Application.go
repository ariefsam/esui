package entity

type Application struct {
	ID   string
	Name string
}

type ApplicationVersion struct {
	ID            string
	ApplicationID string
	Name          string
	Version       string
	Menus         []Menu
}

type Menu struct {
	ID         string
	Title      string
	Type       string
	SubMenus   []SubMenu
	Components []Component
}

type SubMenu struct {
	ID         string
	Title      string
	Components []Component
}

type Component struct {
	ID          string
	Type        string
	WelcomeText string
}
