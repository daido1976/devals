package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	env, err := vals.QuotedEnv(m)
	if err != nil {
		fatal("%v", err)
	}
	// TODO: 出力先は -o 引数で指定できるようにする（未指定ならば標準出力）
	for _, l := range env {
		fmt.Fprintln(os.Stdout, l)
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
