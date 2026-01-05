package main

import "specfirst/cmd"

var version = "0.5.0"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
