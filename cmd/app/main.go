package main

import "backend-trainee-assignment-winter-2025/internal/app"

const (
	configPath = "config/config.yaml"
)

func main() {
	app.Run(configPath)
}
