package main

import "github.com/wernsiet/morchy/cmd/mctl/app"

func main() {
	cmd := app.NewMCTLCommand()
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
