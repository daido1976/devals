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

	// JSONをファイルに書き込む
	if err := os.WriteFile("output.json", jsonData, 0644); err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	m := readOrFail("output.json")

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
