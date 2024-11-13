//go:build tools
// +build tools

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/air-verse/air"
	_ "github.com/joho/godotenv"
	_ "github.com/labstack/echo/v4"
	_ "go.mongodb.org/mongo-driver/mongo"
)
