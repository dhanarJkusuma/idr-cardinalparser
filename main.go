package main

import (
	"fmt"
	"idr-cardinalparser/converter"
)

func main() {
	var conv converter.NumberConverter
	conv = converter.NewNumberConverter()
	number, _ := conv.ConvertToNumber("tiga ribu lima ratus tujuh puluh lima")
	fmt.Println(number)
}
