package main

import (
	"fmt"

	"github.com/AliBa1/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = c.SetUser("ali")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c, err = config.Read()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.Print()
}
