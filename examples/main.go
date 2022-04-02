package main

import (
	"fmt"

	"github.com/pikann/go-verify-apple-id-token/apple"
)

func main() {
	resp, err := apple.VerifyAppleIdToken(
        "token_test",
        "client_id",
    )

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Result: ", resp)
}
