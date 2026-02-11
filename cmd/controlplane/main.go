package main

import "github.com/wernsiet/morchy/cmd/controlplane/app"

func main() {
	cmd := app.NewControlPlaneCommand()
	cmd.Execute()
}
