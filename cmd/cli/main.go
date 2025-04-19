package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/bobmaertz/ollama-agent/pkg/ollama/chat"
	"github.com/bobmaertz/ollama-agent/pkg/tools"
)

const (
	COLOR_RED   = "\033[31m"
	COLOR_GREEN = "\033[32m"
	COLOR_BLUE  = "\033[34m"
	COLOR_CYAN  = "\033[36m"
	COLOR_RESET = "\033[0m"

	defaultUrl   = "http://localhost:11434/api/chat"
	defaultModel = "mistral:7b"
)

var (
	urlFlag   string // URL to the ollama API
	modelFlag string // Model to use for the chat
)

func init() {
	// Define both flags with short and long versions
	flag.StringVar(&urlFlag, "url", defaultUrl, "the URL of the ollama server")
	flag.StringVar(&urlFlag, "u", defaultUrl, "the URL of the ollama server (shorthand)")

	flag.StringVar(&modelFlag, "model", defaultModel, "the model to use; must be installed")
	flag.StringVar(&modelFlag, "m", defaultModel, "the model to use; must be installed (shorthand)")

	flag.Usage = usage
}

func main() {
	// Parse the flags
	flag.Parse()

	// Get remaining positional arguments
	args := flag.Args()

	// Validate we have at least the "chat" subcommand
	if len(args) < 1 || args[0] != "chat" {
		usage()
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	// Create the chat client with the provided flags
	chatter, err := chat.OpenClient(urlFlag, modelFlag)
	if err != nil {
		fmt.Printf("%sError creating chat client: %v%s\n", COLOR_RED, err, COLOR_RESET)
		os.Exit(1)
	}

	fmt.Printf("%sConnected to %s using model %s%s\n", COLOR_CYAN, urlFlag, modelFlag, COLOR_RESET)

	for {
		ctx := context.TODO()

		fmt.Print(COLOR_GREEN + "<you>" + COLOR_RESET)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(COLOR_RED+"An error occurred:"+COLOR_RESET, err)
			continue
		}
		output, err := chatter.Send(ctx, input, chat.RoleUser)
		if err != nil {
			fmt.Println(COLOR_RED+"An error occurred:"+COLOR_RESET, err)
			continue
		}

		for _, tool := range output.Msg.ToolCalls {
			// Only support functions for now
			t := tool["function"]

			fmt.Printf("%s<calling function: %s with args: %v>%s\n", COLOR_BLUE, t.Name, t.Arguments, COLOR_RESET)

			f := tools.Available[t.Name]
			if f == nil {
				fmt.Printf("%stool %v not available: %v%s\n", COLOR_RED, t.Name, tools.Available, COLOR_RESET)
				continue
			}

			tool_resp, _ := f(t.Arguments)
			output, err = chatter.Send(ctx, tool_resp, chat.RoleTool)
			if err != nil {
				fmt.Printf("%sAn error occurred: %v%s\n", COLOR_RED, err, COLOR_RESET)
				break
			}
		}
		fmt.Println(COLOR_CYAN+"<agent>"+COLOR_RESET, output.Msg.Content)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of gollama:\n")
	fmt.Fprintf(os.Stderr, "  gollama [options] chat\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  chat       Chat with the LLM\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExample:\n")
	fmt.Fprintf(os.Stderr, "  gollama -m llama3 -u http://localhost:11434/api/chat chat\n")
}
