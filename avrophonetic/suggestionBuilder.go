package avrophonetic

import (
	"github.com/arbovm/levenshtein"
	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodict"
	"regexp"
	"sort"
	"strings"
)

type byLevenshtein struct {
	baseWord string
	words    *[]string
}

func (l byLevenshtein) Len() int {
	return len(*l.words)
}

func (l byLevenshtein) Swap(i, j int) {
	a := *l.words
	a[i], a[j] = a[j], a[i]
}

func (l byLevenshtein) Less(i, j int) bool {
	a := *l.words

	da := levenshtein.Distance(l.baseWord, a[i])
	db := levenshtein.Distance(l.baseWord, a[j])

	return da < db
}

type splitableWord struct {
	begin  string
	middle string
	end    string
}

type correctableWord struct {
	corrected string
	exact     bool
	invalid   bool
}

type cacheableWord struct {
	base string
	eng  string
}
type Preference struct {
	DictDisabled bool
}

type Suggestion struct {
	Words         []string
	PrevSelection int
}

type SuggestionBuilder struct {
	DBSearch      *avrodict.Searcher
	AutocorrectDB map[string]string
	AvroParser    *avroclassic.Parser
	SuffixDict    map[string]string
	Pref          Preference
	tempCache     map[string]cacheableWord
	phoneticCache map[string][]string
}

func NewBuilder(a *avrodict.Searcher, b map[string]string, c *avroclassic.Parser, d map[string]string, e Preference) *SuggestionBuilder {
	sb := SuggestionBuilder{}
	sb.DBSearch = a
	sb.AutocorrectDB = b
	sb.AvroParser = c
	sb.SuffixDict = d
	sb.Pref = e
    sb.tempCache = make(map[string]cacheableWord)
    sb.phoneticCache = make(map[string][]string)
	return &sb
}

func (avro *SuggestionBuilder) getDictionarySuggestion(splitWord splitableWord) []string {
	key := strings.ToLower(splitWord.middle)

	if words, ok := avro.phoneticCache[key]; ok {
		copiedWords := make([]string, len(words))
		copy(copiedWords, words)
		return copiedWords
	} else {
		return avro.DBSearch.Search(key)
	}
}

func (avro *SuggestionBuilder) getClassicPhonetic(banglish string) string {
	return avro.AvroParser.Parse(banglish)
}

func (avro *SuggestionBuilder) correctCase(banglish string) string {
	return avro.AvroParser.FixString(banglish)
}

func (avro *SuggestionBuilder) getAutocorrect(word string, splitWord splitableWord) correctableWord {
	var corrected correctableWord

	//Search for whole match
	if aWord, ok := avro.AutocorrectDB[word]; ok {
		// [smiley rule]
		if aWord == word {
			corrected.corrected = word
			corrected.exact = true
		} else {
			corrected.corrected = avro.getClassicPhonetic(aWord)
			corrected.exact = false
		}
	} else {
		//Whole word is not present, search without padding
		correctedMiddle := avro.correctCase(splitWord.middle)
		if aWord, ok := avro.AutocorrectDB[correctedMiddle]; ok {
			corrected.corrected = avro.getClassicPhonetic(aWord)
			corrected.exact = false
		} else {
			corrected.invalid = true
		}
	}

	return corrected
}
func (avro *SuggestionBuilder) separatePadding(word string) splitableWord {
	// Mehdi: Feeling lost? Ask Rifat :D
	// Gittu: ?= lookahead error go. check later
	// re := regexp.MustCompile("(^(?::`|\\.`|[\\-\\]~!@#%&*()_=+[{}'\";<>\\/?|.,])*?(?=(?:,{2,}))|^(?::`|\\.`|[\\-\\]~!@#%&*()_=+[{}'\";<>\\/?|.,])*)(.*?(?:,,)*)((?::`|\\.`|[\\-\\]~!@#%&*()_=+[{}'\";<>\\/?|.,])*$)")

	var splitWord splitableWord
	splitWord.begin = ""    //match[1]
	splitWord.middle = word //match[2]
	splitWord.end = ""      //match[3]

	return splitWord
}

func (avro *SuggestionBuilder) sortByPhoneticRelevance(phonetic string, dictSuggestion []string) []string {
	//Copy things into a sortable interface implementation then call sort
	suggSlice := make([]string, len(dictSuggestion))
	copy(suggSlice, dictSuggestion)
	var sortAble byLevenshtein
	sortAble.baseWord = phonetic
	sortAble.words = &suggSlice

	sort.Sort(sortAble)

	return suggSlice
}

func (avro *SuggestionBuilder) isKar(input string) bool {
	if len(input) < 1 {
		return false
	}
	//TODO: implement without regex
	re := regexp.MustCompile("^[ািীুূৃেৈোৌৄ]$")
	return re.MatchString(input[0:1])
}

func (avro *SuggestionBuilder) isVowel(input string) bool {
	if len(input) < 1 {
		return false
	}
	//TODO: implement without regex
	re := regexp.MustCompile("^[অআইঈউঊঋএঐওঔঌৡািীুূৃেৈোৌ]$")
	return re.MatchString(input[0:1])
}

func (avro *SuggestionBuilder) addToTempCache(full string, base string, eng string) {
	if v, ok := avro.tempCache[full]; !ok {
		v.base = base
		v.eng = eng
		avro.tempCache[full] = v
	}
}

func (avro *SuggestionBuilder) clearTempCache() {
	avro.tempCache = make(map[string]cacheableWord)
}

func (avro *SuggestionBuilder) clearDuplicate(data []string) []string {
	//TODO: Implement
	return data
}

func (avro *SuggestionBuilder) addSuffix(splitWord splitableWord) []string {
	var tempSlice []string
	var fullWord string
	word := strings.ToLower(splitWord.middle)
	var rSlice []string

	if v, ok := avro.phoneticCache[word]; ok {
        rSlice = make([]string, len(v))        
		copy(rSlice, v)
	}

	avro.clearTempCache()

	if len(word) > 1 {
		for j, _ := range word {
			var testSuffix = word[j+1:]

			if suffix, ok := avro.SuffixDict[testSuffix]; ok {
				key := word[0 : len(word)-len(testSuffix)]
				if vSlice, ok := avro.phoneticCache[key]; ok {
					for _, cacheItem := range vSlice {
						cacheRightChar := cacheItem[len(cacheItem)-1 : len(cacheItem)]
						suffixLeftChar := suffix[0:1]
						if avro.isVowel(cacheRightChar) && avro.isKar(suffixLeftChar) {
							fullWord = cacheItem + "\u09df" + suffix // \u09df = B_Y
							tempSlice = append(tempSlice, fullWord)
							avro.addToTempCache(fullWord, cacheItem, key)
						} else {
							if cacheRightChar == "\u09ce" { // \u09ce = b_Khandatta
								fullWord = cacheItem[0:len(cacheItem)-1] + "\u09a4" + suffix // \u09a4 = b_T
								tempSlice = append(tempSlice, fullWord)
								avro.addToTempCache(fullWord, cacheItem, key)
							} else if cacheRightChar == "\u0982" { // \u0982 = b_Anushar
								fullWord = cacheItem[0:len(cacheItem)-1] + "\u0999" + suffix // \u09a4 = b_NGA
								tempSlice = append(tempSlice, fullWord)
							} else {
								fullWord = cacheItem + suffix
								tempSlice = append(tempSlice, fullWord)
								avro.addToTempCache(fullWord, cacheItem, key)
							}
						}
					}

					rSlice = append(rSlice, tempSlice...)
				}
			}
		}
	}
    return rSlice
}

func (avro *SuggestionBuilder) getPreviousSelection(splitWord splitableWord, suggestionWords []string) int {
	//TODO: implement
	return 0
}

func (avro *SuggestionBuilder) joinSuggestion(autoCorrect correctableWord, dictSuggestion []string, phonetic string, splitWord splitableWord) Suggestion {
	var words []string

	if avro.Pref.DictDisabled {
		words = append(words, splitWord.begin+phonetic+splitWord.end)
		return Suggestion{words, 0}
	} else {

		/* 1st Item: Autocorrect */
		if autoCorrect.invalid == false {
			words = append(words, autoCorrect.corrected)
			//Add autocorrect entry to dictSuggestion for suffix support
			if !autoCorrect.exact {
				dictSuggestion = append(dictSuggestion, autoCorrect.corrected)
			}
		}

		/* 2rd Item: Dictionary Avro Phonetic */
		//Update Phonetic Cache
		cacheKey := strings.ToLower(splitWord.middle)
		if _, ok := avro.phoneticCache[cacheKey]; !ok {
			if len(dictSuggestion) > 0 {
				copiedSuggestion := make([]string, len(dictSuggestion))
				copy(copiedSuggestion, dictSuggestion)
				avro.phoneticCache[cacheKey] = copiedSuggestion
			}
		}
        
		//Add Suffix
		dictSuggestionWithSuffix := avro.addSuffix(splitWord)

		sortedWords := avro.sortByPhoneticRelevance(phonetic, dictSuggestionWithSuffix)
		words = append(words, sortedWords...)

		/* 3rd Item: Classic Avro Phonetic */
		words = append(words, phonetic)

		words = avro.clearDuplicate(words)

		suggestion := Suggestion{}

		//Is there any previous custom selection of the user?
		suggestion.PrevSelection = avro.getPreviousSelection(splitWord, words)

		//Add padding to all, except exact autocorrect
		for i, _ := range words {
			if autoCorrect.exact {
				if autoCorrect.corrected != words[i] {
					words[i] = splitWord.begin + words[i] + splitWord.end
				}
			} else {
				words[i] = splitWord.begin + words[i] + splitWord.end
			}
		}

		suggestion.Words = words
		return suggestion

	}

}

func (avro *SuggestionBuilder) Suggest(word string) Suggestion {
	//Seperate begining and trailing padding characters, punctuations etc. from whole word
	splitWord := avro.separatePadding(word)

	//Convert begining and trailing padding text to phonetic Bangla
	splitWord.begin = avro.getClassicPhonetic(splitWord.begin)
	splitWord.end = avro.getClassicPhonetic(splitWord.end)

	//Convert the word to Bangla using 3 separate methods
	phonetic := avro.getClassicPhonetic(splitWord.middle)

	if !avro.Pref.DictDisabled {
		dictSuggestion := avro.getDictionarySuggestion(splitWord)
		autoCorrect := avro.getAutocorrect(word, splitWord)
		return avro.joinSuggestion(autoCorrect, dictSuggestion, phonetic, splitWord)
	} else {
		return avro.joinSuggestion(correctableWord{}, nil, phonetic, splitWord)
	}
}
