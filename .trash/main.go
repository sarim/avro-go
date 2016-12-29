package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"

	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avrophonetic"
	"github.com/sarim/avro-go/avroregex"
	// "time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	avroParser := runClassic()
	regexParser := runRegex()
	dBSearch := runDict(regexParser)
	sb := runSB(dBSearch, avroParser)
	if *cpuprofile != "" {
		f, _ := os.Create(*cpuprofile)

		pprof.StartCPUProfile(f)
	}
	for i := 0; i < 1; i++ {
		// startTime := time.Now()
		// fmt.Printf("%#v\n", dBSearch.Search("sari"))
		fmt.Printf("%q\n", sb.Suggest("#sari"))
		// sb.Suggest("#sari")
	}
	pprof.StopCPUProfile()
	// fmt.Printf("Time: %s\n", time.Since(startTime))
}

func runClassic() *avroclassic.Parser {
	dataFile, _ := ioutil.ReadFile("./data/avroclassic.json")
	avroData := avroclassic.AvroData{}
	json.Unmarshal(dataFile, &avroData)
	avro := avroclassic.Parser{avroData}
	// fmt.Printf("%#v\n", avro.Parse("sarim"))
	return &avro
}

func runRegex() *avroregex.Parser {
	dataFile, _ := ioutil.ReadFile("./data/avroregex.json")
	avroData := avroregex.AvroData{}
	json.Unmarshal(dataFile, &avroData)
	avro := avroregex.Parser{avroData}
	// fmt.Printf("%#v\n", avro.Parse("a(!k)"))
	return &avro
}

func runDict(regexParser *avroregex.Parser) *avrodict.Searcher {
	dataFile, _ := ioutil.ReadFile("./data/avrodict.json")
	avroTable := avrodict.AvroTable{}
	json.Unmarshal(dataFile, &avroTable)
	avro := avrodict.Searcher{avroTable, regexParser}
	// _ = avro
	// fmt.Printf("%#v\n", avro.Search("sari"))
	return &avro
}

func runSB(dBSearch *avrodict.Searcher, avroParser *avroclassic.Parser) *avrophonetic.SuggestionBuilder {
	var autocorrectDB map[string]string
	var suffixDict map[string]string
	pref := avrophonetic.Preference{false}

	dataFile, _ := ioutil.ReadFile("./data/autocorrect.json")
	json.Unmarshal(dataFile, &autocorrectDB)

	suffixFile, _ := ioutil.ReadFile("./data/suffixdict.json")
	json.Unmarshal(suffixFile, &suffixDict)

	avrosb := avrophonetic.NewBuilder(dBSearch, autocorrectDB, avroParser, suffixDict, pref)
	// fmt.Printf("%#v\n", avrosb.Suggest("sari"))
	return avrosb
}
