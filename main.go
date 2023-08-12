package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// dotenv ファイルを読み込む
	env, err := godotenv.Read(".env.example")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// JSON に変換
	jsonData, err := json.Marshal(env)
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}
