package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"strings"

	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodata"
	"github.com/sarim/avro-go/avrodict"
	"github.com/sarim/avro-go/avrophonetic"
	"github.com/sarim/avro-go/avroregex"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var iteration = flag.Int("iteration", 1, "iteration count")
var words = flag.String("words", "sari", "Comma seperated test words to run through avro")
var prevSelectCandidate = flag.String("prev-select-candidate", "", "[candidate] for testing Previously Selected Candidate")
var prevSelectWord = flag.String("prev-select-word", "", "[word] for testing Previously Selected Candidate")
var disableDict = flag.Bool("disable-dict", false, "Disable Dictionary Suggestion")

func main() {
	flag.Parse()

	db := avrodata.NewDB()
	pref := avrophonetic.Preference{*disableDict}

	avroParser := avroclassic.Parser{db.Classicdb}
	regexParser := avroregex.Parser{db.Regexdb}
	dBSearch := avrodict.Searcher{db.Dictdb, &regexParser}
	candSelector := avrophonetic.NewInMemoryCandidateSelector(nil)

	if *prevSelectCandidate != "" && *prevSelectWord != "" {
		candSelector.Set(*prevSelectCandidate, *prevSelectWord)
	}

	sb := avrophonetic.NewBuilder(&dBSearch, db.Autocorrect, &avroParser, db.Suffixdb, pref, candSelector)

	if *cpuprofile != "" {
		f, _ := os.Create(*cpuprofile)

		pprof.StartCPUProfile(f)
	}

	wordList := strings.Split(*words, ",")

	for _, word := range wordList {

		startTime := time.Now()

		var suggestion avrophonetic.Suggestion

		for i := 0; i < *iteration; i++ {
			suggestion = sb.Suggest(word)
		}

		fmt.Printf("Time: %s\n", time.Since(startTime))
		fmt.Printf("%q\n", suggestion)
	}

	if *cpuprofile != "" {
		pprof.StopCPUProfile()
	}

}
