package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/helmfile/vals"
	"github.com/joho/godotenv"
)

var (
	inputFile    = flag.String("i", "", "Input .env file (required)")
	outputFile   = flag.String("o", "", "Output file. If not specified, writes to stdout.")
	keepComments = flag.Bool("keep-comments", false, "Keep comments and empty lines in the output")
)

func main() {
	flag.Parse()

	if *inputFile == "" {
		fatal("Input file (-i) is required")
	}

	origEnv, err := godotenv.Read(*inputFile)
	if err != nil {
		fatal("Error loading .env file: %v", err)
	}

	envMap, err := convertToEnvMap(origEnv)
	if err != nil {
		fatal("Error converting to env map: %v", err)
	}

	output := getOutputWriter()
	defer output.Close()

	if *keepComments {
		writeWithComments(*inputFile, envMap, output)
	} else {
		writeWithoutComments(envMap, output)
	}
}

func convertToEnvMap(origEnv map[string]string) (map[string]string, error) {
	jsonData, err := json.Marshal(origEnv)
	if err != nil {
		return nil, err
	}

	m, err := convertJSONToMap(string(jsonData))
	if err != nil {
		return nil, err
	}

	envLines, err := vals.QuotedEnv(m)
	if err != nil {
		return nil, err
	}

	envMap := make(map[string]string)
	for _, line := range envLines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap, nil
}

func convertJSONToMap(jsonData string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	return data, err
}

func getOutputWriter() *os.File {
	if *outputFile == "" {
		return os.Stdout
	}

	output, err := os.Create(*outputFile)
	if err != nil {
		fatal("Error creating output file: %v", err)
	}
	return output
}

func writeWithComments(inputFile string, envMap map[string]string, output *os.File) {
	file, err := os.Open(inputFile)
	if err != nil {
		fatal("Error opening .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// コメントまたは空行の場合、そのまま出力
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			fmt.Fprintln(output, line)
			continue
		}

		// コメントでない場合、変換された環境変数を出力
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			if val, exists := envMap[key]; exists {
				fmt.Fprintf(output, "%s=%s\n", key, val)
			} else {
				// キーが envMap にない場合はそのまま出力
				fmt.Fprintln(output, line)
			}
		}
	}

	if scanner.Err() != nil {
		fatal("Error reading .env file: %v", scanner.Err())
	}
}

func writeWithoutComments(envMap map[string]string, output *os.File) {
	for key, val := range envMap {
		fmt.Fprintf(output, "%s=%s\n", key, val)
	}
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
