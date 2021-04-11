package main

import (
	"fmt"
	"os"
)

var d = make(chan int)

func init() {
	close(d)
}
func main1() {
	r, error := os.Stat("/home/binnz/Desktop/sales.json")
	fmt.Println(r)
	fmt.Println(error)
}
