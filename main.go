package main

import (
	"encoding/json"
	"fmt"
	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avrophonetic"
	"github.com/sarim/avro-go/avroregex"
	"io/ioutil"
	"time"
)

func main() {
	avroParser := runClassic()
	regexParser := runRegex()
	dBSearch := runDict(regexParser)
	avrosb := runSB(dBSearch, avroParser)
	startTime := time.Now()
    fmt.Printf("%#v\n", avrosb.Suggest("sari"))
	fmt.Printf("Time: %s\n", time.Since(startTime))
}

func runClassic() *avroclassic.Parser {
	dataFile, _ := ioutil.ReadFile("./avroclassic/avroclassic.json")
	avroData := avroclassic.AvroData{}
	json.Unmarshal(dataFile, &avroData)
	avro := avroclassic.Parser{avroData}
	// fmt.Printf("%#v\n", avro.Parse("sarim"))
	return &avro
}

func runRegex() *avroregex.Parser {
	dataFile, _ := ioutil.ReadFile("./avroregex/avroregex.json")
	avroData := avroregex.AvroData{}
	json.Unmarshal(dataFile, &avroData)
	avro := avroregex.Parser{avroData}
	// fmt.Printf("%#v\n", avro.Parse("a(!k)"))
	return &avro
}

func runDict(regexParser *avroregex.Parser) *avrodict.Searcher {
	dataFile, _ := ioutil.ReadFile("./avrodict/avrodict.json")
	avroTable := avrodict.AvroTable{}
	json.Unmarshal(dataFile, &avroTable)
	avro := avrodict.Searcher{avroTable, regexParser}
	_ = avro
	// fmt.Printf("%#v\n", avro.Search("sari"))
	return &avro
}

func runSB(dBSearch *avrodict.Searcher, avroParser *avroclassic.Parser) *avrophonetic.SuggestionBuilder {
	var autocorrectDB map[string]string
	var suffixDict map[string]string
	pref := avrophonetic.Preference{false}

	dataFile, _ := ioutil.ReadFile("./avrophonetic/autocorrect.json")
	json.Unmarshal(dataFile, &autocorrectDB)

	suffixFile, _ := ioutil.ReadFile("./avrophonetic/suffixdict.json")
	json.Unmarshal(suffixFile, &suffixDict)

	avrosb := avrophonetic.NewBuilder(dBSearch, autocorrectDB, avroParser, suffixDict, pref)
	// fmt.Printf("%#v\n", avrosb.Suggest("sari"))
	return avrosb
}
