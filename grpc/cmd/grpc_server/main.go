package main

import (
	"context"
	"github.com/f1xend/auth/internal/app"
	"log"
)

//
//var configPath string

//func init() {
//	flag.StringVar(&configPath, "config-path", "prod.env", "path to config file")
//}

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	if err = a.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
