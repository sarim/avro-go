package main

import (
	"encoding/gob"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	_ "github.com/alecthomas/binary"
	_ "github.com/gogo/protobuf/proto"
	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avroregex"
)

func main() {

	gobFile, _ := os.Open("./data/compiled.gob")

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
	binFile, _ := ioutil.ReadFile("./data/compiled.bson")
	bson.Unmarshal(binFile, &db1)
	fmt.Printf("Time for bson: %v\n", time.Since(startTime))

	startTime = time.Now()
	dec := gob.NewDecoder(gobFile)
	err := dec.Decode(&db2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Time for gob: %v\n", time.Since(startTime))

	fmt.Println(len(db2.Autocorrect))

}
