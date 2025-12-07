package main

import "github.com/wernsiet/morchy/cmd/agent/app"

func main() {
	cmd := app.NewAgentCommand()
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
