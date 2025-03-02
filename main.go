package main

import (
	"demo/api"
	"fmt"
	"log"

	"github.com/ianschenck/envflag"
)

func main() {
	fmt.Println("jai baba ri")
	//secret key
	secretKey := envflag.String("SECRET_KEY", "3712394871491741947194734729473219471293847", "secret key for auth in")


	conn, err := api.NewPostgreStore()
	if err != nil {
		log.Fatalf("error at postgre store: %v", err)
		return
	}
	if errr := conn.Init(); err != nil {
		log.Fatalf("error while establishing table: %v", errr)
	}
	server := api.NewApiServer(":3000", conn, *secretKey)
	server.Run()
}