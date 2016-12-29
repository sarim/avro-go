package main

import (
	"encoding/gob"
	"encoding/json"

	_ "github.com/alecthomas/binary"
	_ "github.com/gogo/protobuf/proto"
	vitessbson "github.com/youtube/vitess/go/bson"
	// "fmt"
	"io/ioutil"
	"os"

	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avroregex"
)

type DB struct {
	Autocorrect map[string]string
	Classicdb   avroclassic.AvroData
	Dictdb      avrodict.AvroTable
	Regexdb     avroregex.AvroData
	Suffixdb    map[string]string
}

func main() {

	file1, _ := ioutil.ReadFile("./data/autocorrect.json")
	file2, _ := ioutil.ReadFile("./data/avroclassic.json")
	file3, _ := ioutil.ReadFile("./data/avrodict.json")
	file4, _ := ioutil.ReadFile("./data/avroregex.json")
	file5, _ := ioutil.ReadFile("./data/suffixdict.json")
	var db DB
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

	jsonData, _ := json.Marshal(db)
	bsonData, _ := vitessbson.Marshal(db)

	jsonFile, _ := os.Create("./data/compiled.json")
	bsonFile, _ := os.Create("./data/compiled.bson")
	binFile, _ := os.Create("./data/compiled.gob")

	enc := gob.NewEncoder(binFile)
	enc.Encode(db)

	jsonFile.Write(jsonData)
	bsonFile.Write(bsonData)

}
