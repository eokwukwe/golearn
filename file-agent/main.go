package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"

	"eokwukwe/fileagent/agent"
	"eokwukwe/fileagent/tools"
)

func main() {
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	tools := []tools.ToolDefinition{
		tools.ReadFileDefinition,
		tools.ListFilesDefinition,
	}
	agent := agent.NewAgent(&client, getUserMessage, tools)
	err := agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}
}
