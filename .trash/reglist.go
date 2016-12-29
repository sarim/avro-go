package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/sarim/avro-go/avroregex"
)

func main() {
	dataFile, _ := ioutil.ReadFile("./data/avroregex.json")
	avroData := avroregex.AvroData{}
	json.Unmarshal(dataFile, &avroData)

	charMap := make(map[string]bool)

	for _, pattern := range avroData.Patterns {
		for _, v := range pattern.Replace {
			charMap[string(v)] = true
		}
	}

	var charSlice []string
	for char, _ := range charMap {
		charSlice = append(charSlice, char)
	}

	sort.Strings(charSlice)

	for _, s := range charSlice {
		fmt.Println(s)
	}

}
