package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	text := "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(text))
}
