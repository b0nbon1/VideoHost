package main

import (

	"github.com/b0nbon1/VidFlux/app"
)

//	@title			The ultimate Streaming API
//	@version		0.1
//	@description	Todo list API using Fiber and Postgres
//	@contact.name	Bonvic Bundi
//	@license.name	MIT
//	@host			localhost:4500
//	@BasePath		/
func main() {
	// setup and run app
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}