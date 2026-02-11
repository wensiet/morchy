package main

import "github.com/wernsiet/morchy/cmd/mctl/app"

func main() {
	cmd := app.NewMCTLCommand()
	cmd.Execute()
}
