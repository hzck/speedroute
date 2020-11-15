// Package main creates the web server.
package main

func main() {
	a := App{}

	logDefer := a.InitLogFile()
	defer logDefer()

	a.InitConfigFile()

	dbDefer := a.InitDB()
	defer dbDefer()

	a.InitRoutes()

	a.Run()
}
