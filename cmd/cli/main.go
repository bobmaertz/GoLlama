package main

import (
	"bufio"
	"context"
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
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	//TODO: Proper Error handling
	chatter, _ := chat.OpenClient("http://sojourner:11434/api/chat", "mistral:7b")
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
			// TOOD: Only support\ functions for now
			t := tool["function"]

			fmt.Printf("%s<calling function: %s with args: %v>%s\n", COLOR_BLUE, t.Name, t.Arguments, COLOR_RESET)

			f := tools.Available[t.Name]
			if f == nil {
				//TODO: Cleanup
				fmt.Printf("%stool %v not available: %v%s\n", COLOR_RED, t.Name, tools.Available, COLOR_RESET)
				continue
			}

			//TODO: Cleanup error handling
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
