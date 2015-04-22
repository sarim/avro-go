package avrodict

import "github.com/sarim/avro-go/avroregex"
import "regexp"
import "unicode"


type Searcher struct {
    Table AvroTable
    Regex *avroregex.Parser
}

func (avro *Searcher) Search(enText string) []string {
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
     pattern = "^" + pattern + "$"
     
     var retWords []string
    
     for _, tn := range tableList {
         table := "w_" + tn
         nextRetWords := avro.searchInSlice(pattern, table)
         retWords = append(retWords, nextRetWords...)
     }
    
    return retWords
}

func (avro *Searcher) searchInSlice(pattern string, tableName string) []string {
    var retWords []string
    
    //TODO: Handle error here
    re, _ := regexp.Compile(pattern)

    for _, word := range avro.Table[tableName] {
        if re.MatchString(word) {
            retWords = append(retWords, word)
        }
    }
    return retWords
}