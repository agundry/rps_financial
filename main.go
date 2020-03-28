package main

import (
	"./app"
)

func main() {
	application := app.App{}
	application.Initialize(
		"root",
		"my-secret-pw",
		"localhost:13306",
		"rps")

	application.Run(":8080")
}
