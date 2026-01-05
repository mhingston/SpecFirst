package cmd

var version = "dev"

func SetVersion(v string) {
	if v != "" {
		version = v
	}
}
