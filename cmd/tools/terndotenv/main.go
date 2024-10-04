package main

import (
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	database := os.Getenv("DATABASE")

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--port",
		port,
		"--host",
		host,
		"--user",
		user,
		"--password",
		password,
		"--database",
		database,
	)

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
