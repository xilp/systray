package main

import (
	"bufio"
	"os"
	"github.com/xilp/systray"
)

func main() {
	if len(os.Args) != 3 {
		println("usage: example icon-path client-path")
		return
	}

	tray := systray.New(os.Args[1], os.Args[2], 6333)
	tray.OnClick(func() {
		println("clicked")
	})
	err := tray.Show("idle.ico", "Test systray")
	if err != nil {
		println(err.Error())
	}

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			println("Input icon file name:")
			print(">> ")
			data, _, _ := reader.ReadLine()
			line := string(data)
			if len(line) == 0 {
				break
			}
			err := tray.Show(line, line)
			if err != nil {
				println(err.Error())
			}
		}

		err = tray.Stop()
		if err != nil {
			println(err.Error())
		}
		os.Exit(0)
	}()

	err = tray.Run()
	if err != nil {
		println(err.Error())
	}
}
