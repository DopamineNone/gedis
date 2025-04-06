package main

import "log"

func main() {
	application := wireApp()

	go func() {
		if err := application.Run(); err != nil {
			log.Println(err.Error())
		}
	}()

	application.ListenAndQuit()
}
