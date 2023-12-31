package cli

import (
	"fmt"
	"log"
	"os"
)

func Run() {
	if len(os.Args) != 4 {
		return
	}

	token, err := GenerateToken(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated Token:", token)
	os.Exit(0)
}
