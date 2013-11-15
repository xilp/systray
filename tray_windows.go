package systray

func _NewSystray(iconPath string, clientPath string) _Systray {
	return _Systray{iconPath}
}

type _Systray struct {
	iconPath string
}
