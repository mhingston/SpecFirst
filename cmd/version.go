package cmd

var version = "v0.5.0"

func SetVersion(v string) {
	if v != "" {
		version = v
		rootCmd.Version = v
	}
}
