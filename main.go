package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/helmfile/vals"
	"github.com/joho/godotenv"
)

func main() {
	inputFile := flag.String("i", "", "Input dotenv format file (required)")
	outputFile := flag.String("o", "", "Output file. If not specified, writes to stdout.")
	keepComments := flag.Bool("keep-comments", false, "Keep comments and empty lines in the output")

	flag.Parse()

	if *inputFile == "" {
		fatal("Input file (-i) is required")
	}

	input, err := os.ReadFile(*inputFile)
	if err != nil {
		fatal("Error opening input file: %v", err)
	}

	output := getOutputWriter(*outputFile)
	defer output.Close()

	err = run(string(input), output, *keepComments)
	if err != nil {
		fatal("Error converting file: %v", err)
	}
}

// Gets the output writer based on the specified file path or defaults to stdout.
func getOutputWriter(outputFile string) *os.File {
	if outputFile == "" {
		return os.Stdout
	}

	output, err := os.Create(outputFile)
	if err != nil {
		fatal("Error creating output file: %v", err)
	}
	return output
}

// TODO: Currently, the 'input' is of type 'string', which may potentially lead to higher memory consumption. Consider changing it to a more memory-efficient type.
func run(input string, output io.Writer, keepComments bool) error {
	envMap, err := godotenv.Unmarshal(input)
	if err != nil {
		return fmt.Errorf("error loading dotenv format file: %v", err)
	}

	resolvedEnvMap, err := resolveRefInEnvMap(envMap)
	if err != nil {
		return fmt.Errorf("error converting to env map: %v", err)
	}

	if keepComments {
		writeWithComments(input, resolvedEnvMap, output)
	} else {
		writeWithoutComments(resolvedEnvMap, output)
	}

	return nil
}

// Resolves references in the original environment map.
func resolveRefInEnvMap(envMap map[string]string) (map[string]string, error) {
	// NOTE: Cast the map to map[string]any for vals.QuotedEnv.
	data := make(map[string]any)
	for key, value := range envMap {
		data[key] = value
	}

	envLines, err := vals.QuotedEnv(data)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, line := range envLines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m, nil
}

// Writes to the output with original comments and empty lines preserved.
func writeWithComments(input string, envMap map[string]string, output io.Writer) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()

		// If it's a comment or an empty line, output as is.
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			fmt.Fprintln(output, line)
			continue
		}

		// If not a comment, output the converted environment variable.
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			if val, exists := envMap[key]; exists {
				fmt.Fprintf(output, "%s=%s\n", key, val)
			} else {
				// If the key is not in envMap, output as is.
				fmt.Fprintln(output, line)
			}
		}
	}

	if scanner.Err() != nil {
		fatal("Error reading .env file: %v", scanner.Err())
	}
}

// Writes to the output without preserving the original comments and empty lines.
func writeWithoutComments(envMap map[string]string, output io.Writer) {
	for key, val := range envMap {
		fmt.Fprintf(output, "%s=%s\n", key, val)
	}
}

// Exits the program with an error message.
func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
