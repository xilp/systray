package systray

func _NewSystray(iconPath string, clientPath string, port int) *_Systray {
	return &_Systray{_NewSystraySvr(iconPath, clientPath + "/systray.app/Contents/MacOS/systray", port)}
}

type _Systray struct {
	*_SystraySvr
}
