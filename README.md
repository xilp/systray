Systray
=======


## Cross platform systray (trayicon/menu extras) for golang

Win:  
    [your code in go] -> [win32 api call in c]

Mac:  
    [your code in go] -> [systray.Server in go] -(tcp)-> [systray.Client in objc]

Linux:  
    [your code in go] -> [systray.Server in go] -(tcp)-> [systray.Client in c]


## Example

Mac:
```
cd example
go run icons/mac systray
```
