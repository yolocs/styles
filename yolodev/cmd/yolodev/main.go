package main

import (
	"context"
	"os"

	"github.com/yolocs/styles/yolodev/internal/app"
)

func main() {
	os.Exit(app.Run(context.Background(), os.Stdin, os.Stdout, os.Stderr, os.Args[1:]))
}
