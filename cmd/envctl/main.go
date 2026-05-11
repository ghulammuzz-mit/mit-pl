package main

import (
	"mit/platform/internal/envctl/controller"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	controller.Execute()
}
