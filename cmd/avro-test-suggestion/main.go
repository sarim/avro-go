package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodata"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avrophonetic"
	"github.com/sarim/avro-go/avroregex"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var iteration = flag.Int("iteration", 1, "iteration count")
var word = flag.String("word", "sari", "Test Word to run through avro")
var disableDict = flag.Bool("disable-dict", false, "Disable Dictionary Suggestion")

func main() {
	flag.Parse()

	db := avrodata.NewDB()
	pref := avrophonetic.Preference{*disableDict}

	avroParser := avroclassic.Parser{db.Classicdb}
	regexParser := avroregex.Parser{db.Regexdb}
	dBSearch := avrodict.Searcher{db.Dictdb, &regexParser}
	sb := avrophonetic.NewBuilder(&dBSearch, db.Autocorrect, &avroParser, db.Suffixdb, pref)

	if *cpuprofile != "" {
		f, _ := os.Create(*cpuprofile)

		pprof.StartCPUProfile(f)
	}

	startTime := time.Now()

	var suggestion avrophonetic.Suggestion

	for i := 0; i < *iteration; i++ {
		suggestion = sb.Suggest(*word)
	}

	if *cpuprofile != "" {
		pprof.StopCPUProfile()
	}

	fmt.Printf("Time: %s\n", time.Since(startTime))
	fmt.Printf("%q\n", suggestion)
}
