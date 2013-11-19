package systray

func New(iconPath string, clientPath string, port int) *Systray {
	return &Systray{_NewSystray(iconPath, clientPath, port)}
}

type Systray struct {
	*_Systray
}
