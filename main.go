package main

import (
	"fmt"
	"github.com/mahdi-cpp/api-go-pkg/common_models"
	"github.com/mahdi-cpp/api-go-pkg/network"
	"log"
)

func main() {

	userControl := network.NewNetworkControl[[]common_models.User]("http://localhost:8080/api/v1/user/")

	// Make request (nil body if not needed)
	users, err := userControl.Read("list", nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Use the data
	for _, user := range *users {
		fmt.Printf("%d: %s (%s %s)\n",
			user.ID,
			user.Username,
			user.FirstName,
			user.LastName)
	}
}
