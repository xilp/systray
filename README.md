Systray
=======


## Cross platform systray (trayicon/menu extras) for golang

Instead of gui program, "go-program + systray + web-console" might be a interesting choise.


## Example

Mac:
```
cd example
go run icons/mac systray
```


## How it works

Win:  
    [your code in go] -> [win32 api call in c]

Mac:  
    [your code in go] -> [systray.Server in go] -(tcp)-> [systray.Client in objc]

Linux:  
    [your code in go] -> [systray.Server in go] -(tcp)-> [systray.Client in c]


