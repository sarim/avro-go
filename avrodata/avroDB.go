//go:generate avro-data-generate
//go:generate go-bindata -nomemcopy -pkg avrodata -o ./compiled.go ./compiled.gob

package avrodata

import (
	"bytes"
	"encoding/gob"

	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avroregex"
)

type AvroDB struct {
	Autocorrect map[string]string
	Classicdb   avroclassic.AvroData
	Dictdb      avrodict.AvroTable
	Regexdb     avroregex.AvroData
	Suffixdb    map[string]string
}

func NewDB() *AvroDB {
	gobData, _ := Asset("compiled.gob")
	gobAsset := bytes.NewReader(gobData)
	var db AvroDB
	dec := gob.NewDecoder(gobAsset)
	err := dec.Decode(&db)
	if err != nil {
		panic(err)
	}
	return &db
}
