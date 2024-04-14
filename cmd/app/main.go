package main

import "avito-test2024-spring/internal/app"

const ConfigPath = "../../configs/main"

// TODO docker + docker compose + prometheus + grafana
// TODO e2e test for get user_banner

func main() {
	app.Run(ConfigPath)
}
