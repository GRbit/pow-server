package main

import "github.com/GRbit/pow-server/inernal/app/client"

func main() {
	if err := client.GetWord(); err != nil {
		panic(err)
	}
}
