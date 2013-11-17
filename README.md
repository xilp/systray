Systray (Trayicon/Menu Extras)
=======


## Cross platform systray for golang

Instead of gui program, "go-program + systray + web-console" might be a interesting choise.


# Platform

Mac: avalid  
Win: avalid  
Linux: coming soon  


## Run example

Mac:
```
cd example
go run icons/mac systray
```

## How it works

Win:  
    [your code in go] -> [systray: win32 api call in go]

Mac:  
    [your code in go] -> [systray.Server in go] -(tcp)-> [systray.Client in objc]

Linux:  
    [your code in go] -> [systray.Server in go] -(tcp)-> [systray.Client in c]


