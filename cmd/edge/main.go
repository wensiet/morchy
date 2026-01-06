package main

import "github.com/wernsiet/morchy/cmd/edge/app"

func main() {
	cmd := app.NewEdgeCommand()
	cmd.Execute()
}
