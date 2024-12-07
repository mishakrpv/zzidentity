package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zzidentity/zzidentity/pkg/config"
)

func main() {
	fmt.Println("It works!")

	cfg, err := config.New(os.Getenv("CONFIG_FILE"))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	jsonConf, _ := json.Marshal(cfg)

	fmt.Println(string(jsonConf))
}
