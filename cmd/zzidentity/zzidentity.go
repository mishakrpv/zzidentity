package main

import (
	"encoding/json"
	"fmt"

	"github.com/zzidentity/zzidentity/pkg/config"
)

func main() {
	fmt.Println("It works!")

	cfg, err := config.Load("settings.yaml")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	jsonConf, _ := json.Marshal(cfg)

	fmt.Println(string(jsonConf))
}
