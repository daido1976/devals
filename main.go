package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/helmfile/vals"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func main() {
	// dotenv ファイルを読み込む
	// TODO: ファイル名は -i 引数で指定できるようにする
	origEnv, err := godotenv.Read(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// JSON に変換
	jsonData, err := json.Marshal(origEnv)
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	// 一時ファイルを作成する
	tempFile, err := os.CreateTemp("", "temp-output-*.json")
	if err != nil {
		log.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // 処理が終わった後、一時ファイルを削除する

	// JSONを一時ファイルに書き込む
	if _, err := tempFile.Write(jsonData); err != nil {
		log.Fatalf("Error writing JSON to temp file: %v", err)
	}
	tempFile.Close()

	m := readOrFail(tempFile.Name())

	envLines, err := vals.QuotedEnv(m)
	if err != nil {
		log.Fatalf("Error converting map to environment lines: %v", err)
	}

	envMap := make(map[string]string)
	for _, line := range envLines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	file, err := os.Open(".env")
	if err != nil {
		log.Fatalf("Error opening .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// コメントまたは空行の場合、そのまま出力
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			fmt.Println(line)
			continue
		}

		// コメントでない場合、変換された環境変数を出力
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			if val, exists := envMap[key]; exists {
				fmt.Printf("%s=%s\n", key, val)
			} else {
				// キーが envMap にない場合はそのまま出力
				fmt.Println(line)
			}
		}
	}

	if scanner.Err() != nil {
		log.Fatalf("Error reading .env file: %v", scanner.Err())
	}
}

func readNodesOrFail(f string) []yaml.Node {
	nodes, err := vals.Inputs(f)
	if err != nil {
		fatal("%v", err)
	}
	return nodes
}

func readOrFail(f string) map[string]interface{} {
	nodes := readNodesOrFail(f)
	if len(nodes) == 0 {
		fatal("no document found")
	}
	var nodeValue map[string]interface{}
	err := nodes[0].Decode(&nodeValue)
	if err != nil {
		fatal("%v", err)
	}
	return nodeValue
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
