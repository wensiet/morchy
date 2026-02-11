package main

import "github.com/wernsiet/morchy/cmd/agent/app"

func main() {
	cmd := app.NewAgentCommand()
	cmd.Execute()
}
