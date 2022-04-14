package main

import (
	"fmt"
	"github.com/andewx/dieselfluid/dslapp"
)

func main() {
	electronApplication, err := dslapp.New()
	if err != nil {
		fmt.Printf("Failed to launch application %s", err)
		return

	}

	err = electronApplication.Init() /*Waiting Worker Threat*/

	if err != nil {
		fmt.Printf("Failed to launch window %s", err)
		return
	}
}
