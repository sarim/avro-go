package avrodict

import (
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/sarim/avro-go/avroregex"
	"github.com/sarim/gtre"
)

type Searcher struct {
	Table AvroTable
	Regex *avroregex.Parser
}

func (avro *Searcher) Search(enText string) []string {
	if len(enText) == 0 {
		return []string{}
	}
	lmc := unicode.ToLower(rune(enText[0]))
	var tableList []string
	switch lmc {
	case 'a':
		tableList = []string{"a", "aa", "e", "oi", "o", "nya", "y"}
	case 'b':
		tableList = []string{"b", "bh"}
	case 'c':
		tableList = []string{"c", "ch", "k"}
	case 'd':
		tableList = []string{"d", "dh", "dd", "ddh"}
	case 'e':
		tableList = []string{"i", "ii", "e", "y"}
	case 'f':
		tableList = []string{"ph"}
	case 'g':
		tableList = []string{"g", "gh", "j"}
	case 'h':
		tableList = []string{"h"}
	case 'i':
		tableList = []string{"i", "ii", "y"}
	case 'j':
		tableList = []string{"j", "jh", "z"}
	case 'k':
		tableList = []string{"k", "kh"}
	case 'l':
		tableList = []string{"l"}
	case 'm':
		tableList = []string{"h", "m"}
	case 'n':
		tableList = []string{"n", "nya", "nga", "nn"}
	case 'o':
		tableList = []string{"a", "u", "uu", "oi", "o", "ou", "y"}
	case 'p':
		tableList = []string{"p", "ph"}
	case 'q':
		tableList = []string{"k"}
	case 'r':
		tableList = []string{"rri", "h", "r", "rr", "rrh"}
	case 's':
		tableList = []string{"s", "sh", "ss"}
	case 't':
		tableList = []string{"t", "th", "tt", "tth", "khandatta"}
	case 'u':
		tableList = []string{"u", "uu", "y"}
	case 'v':
		tableList = []string{"bh"}
	case 'w':
		tableList = []string{"o"}
	case 'x':
		tableList = []string{"e", "k"}
	case 'y':
		tableList = []string{"i", "y"}
	case 'z':
		tableList = []string{"h", "j", "jh", "z"}
	}

	pattern := avro.Regex.Parse(enText)
	if len(pattern) == 0 {
		return []string{}
	}
	//TODO: Handle error here
	re := gtre.Parse([]rune(pattern))

	var count int32 = 0
	ch := make(chan string, 30)
	var wg sync.WaitGroup

	for _, tn := range tableList {
		wordChunk := avro.Table["w_"+tn]
		for i := range wordChunk {

			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {

				for _, word := range wordChunk[i] {
					if re.Match(word) {
						ch <- string(word)
						atomic.AddInt32(&count, 1)
					}
				}
				defer wg.Done()

			}(&wg, i)

		}
	}
	wg.Wait()

	finalCount := int(atomic.LoadInt32(&count))
	//fmt.
	retWords := make([]string, finalCount)

	for i := 0; i < finalCount; i++ {
		retWords[i] = <-ch
	}

	return retWords
}
