package main

import "github.com/wernsiet/morchy/cmd/controlplane/app"

func main() {
	cmd := app.NewControlPlaneCommand()
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
