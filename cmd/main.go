package main

import (
	"log"

	_ "github.com/DopamineNone/gedis/internal/command"
)

func main() {
	application := wireApp()

	go func() {
		if err := application.Run(); err != nil {
			log.Println(err.Error())
		}
	}()

	application.ListenAndQuit()
}
