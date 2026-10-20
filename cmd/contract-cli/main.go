package main

import (
	"context"
	"fmt"
	"os"

	"cn.qfei/contract-cli/internal/cli"
	"cn.qfei/contract-cli/internal/invocation"
)

func main() {
	if handled, err := invocation.RunInspectionHelper(os.Args[1:], os.Stdout); handled {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	app := cli.New(cli.Options{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	if err := app.Run(context.Background(), os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
