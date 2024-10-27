package main

import (
	"fmt"
	"product-service/pkg/config"
)

func main() {
	conf := config.Load()

	fmt.Println(conf.DB_NAME)
}
