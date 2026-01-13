package cmd

var version = "v0.6.9"

func SetVersion(v string) {
	if v != "" {
		version = v
		rootCmd.Version = v
	}
}
