package main

import (
	"fmt"
	"log"

	"github.com/edwin-lin777/credit-card/internal/catalog"
)

func main() {
	c, err := catalog.LoadCard("cards/scotia-momentum.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", c)

}
