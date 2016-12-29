package main

import (
	"encoding/gob"
	// "encoding/json"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avroregex"
)

func main() {

	gobData, _ := Asset("data/compiled.gob")

	gobAsset := bytes.NewReader(gobData)
	gobFile, _ := os.Open("../data/compiled.gob")

	var db1, db2 struct {
		Autocorrect map[string]string
		Classicdb   avroclassic.AvroData
		Dictdb      avrodict.AvroTable
		Regexdb     avroregex.AvroData
		Suffixdb    map[string]string
	}

	// db.Autocorrect = make(map[string]string)
	// db.Classicdb = avroclassic.AvroData{}
	// db.Dictdb = avrodict.AvroTable{}
	// db.Regexdb = avroregex.AvroData{}
	// db.Suffixdb = make(map[string]string)

	startTime := time.Now()
	dec := gob.NewDecoder(gobAsset)
	err := dec.Decode(&db1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Time for asset gob: %v\n", time.Since(startTime))

	startTime = time.Now()
	dec2 := gob.NewDecoder(gobFile)
	err = dec2.Decode(&db2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Time for file gob: %v\n", time.Since(startTime))

	fmt.Println(len(db2.Autocorrect))

}
