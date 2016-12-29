package main

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"os"

	"fmt"

	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodata"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avroregex"
)

var output = `
Autocorrect Entries:    %d
Classic Patterns:       %d
Dictionary Tables:      %d
Regex Patterns:         %d
SuffixDict Entries:     %d`

func main() {
	file1, _ := ioutil.ReadFile("./autocorrect.json")
	file2, _ := ioutil.ReadFile("./avroclassic.json")
	file3, _ := ioutil.ReadFile("./avrodict.json")
	file4, _ := ioutil.ReadFile("./avroregex.json")
	file5, _ := ioutil.ReadFile("./suffixdict.json")
	var db avrodata.AvroDB
	db.Autocorrect = make(map[string]string)
	db.Classicdb = avroclassic.AvroData{}
	db.Dictdb = avrodict.AvroTable{}
	db.Regexdb = avroregex.AvroData{}
	db.Suffixdb = make(map[string]string)

	json.Unmarshal(file1, &db.Autocorrect)
	json.Unmarshal(file2, &db.Classicdb)
	json.Unmarshal(file3, &db.Dictdb)
	json.Unmarshal(file4, &db.Regexdb)
	json.Unmarshal(file5, &db.Suffixdb)

	fmt.Printf(output,
		len(db.Autocorrect),
		len(db.Classicdb.Patterns),
		len(db.Dictdb),
		len(db.Regexdb.Patterns),
		len(db.Suffixdb))

	gobFile, _ := os.Create("./compiled.gob")

	enc := gob.NewEncoder(gobFile)
	err := enc.Encode(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n[DONE]")
}
