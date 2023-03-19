package main

import (
	"fmt"

	cf "github.com/essentialkaos/go-confluence/v5"
)

func main() {
	api, err := cf.NewAPIWithToken("https://wiki.mediatek.inc", "NDg0OTAwMTg2OTQ1OlZ9h/O7g7y9FRtuTJRk3uvRQ/cP")
	api.SetUserAgent("MyApp", "1.2.3")
	fmt.Println(len("NDg0OTAwMTg2OTQ1OlZ9h/O7g7y9FRtuTJRk3uvRQ/cP"))

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	content, err := api.GetContentByID(
		"1199584586", cf.ContentIDParameters{
			Version: 4,
			Expand:  []string{"space", "body.view", "version"},
		},
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("ID: %s\n", content.ID)
}
