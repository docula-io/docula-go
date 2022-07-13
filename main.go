package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/docula-io/docula/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := cmd.Execute(ctx); err != nil {
		fmt.Println("failed to execute command:", err)
	}
}
