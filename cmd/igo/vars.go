package main

const (
	// appName is an identifier-like name used anywhere this app needs to be identified.
	//
	// It identifies the application itself, the actual instance needs to be identified via environment
	// and other details.
	appName = "igo"

	// friendlyAppName is the visible name of the application.
	friendlyAppName = "iCloud Client in go"

	// envPrefix is prepended to environment variables when processing configuration.
	envPrefix = "igo"
)
