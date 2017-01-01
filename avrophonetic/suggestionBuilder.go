package avrophonetic

import (
	"sort"
	"strings"

	"unicode/utf8"

	"github.com/arbovm/levenshtein"
	"github.com/sarim/avro-go/avroclassic"
	"github.com/sarim/avro-go/avrodict"
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
	CandSelector  CandidateSelector
	tempCache     map[string]cacheableWord
	phoneticCache map[string][]string
}

func NewBuilder(a *avrodict.Searcher, b map[string]string, c *avroclassic.Parser, d map[string]string, e Preference, f CandidateSelector) *SuggestionBuilder {
	sb := SuggestionBuilder{}
	sb.DBSearch = a
	sb.AutocorrectDB = b
	sb.AvroParser = c
	sb.SuffixDict = d
	sb.Pref = e
	sb.tempCache = make(map[string]cacheableWord)
	sb.phoneticCache = make(map[string][]string)
	if f == nil {
		sb.CandSelector = NewInMemoryCandidateSelector(nil)
	} else {
		sb.CandSelector = f
	}
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
	// re := regexp.MustCompile("(^(?::`|\\.`|[\\-\\]~!@#%&*()_=+[{}'\";<>\\/?|.,])*?(?=(?:,{2,}))|^(?::`|\\.`|[\\-\\]~!@#%&*()_=+[{}'\";<>\\/?|.,])*)(.*?(?:,,)*)((?::`|\\.`|[\\-\\]~!@#%&*()_=+[{}'\";<>\\/?|.,])*$)")

	/*begin part:
	  start (:`)or(.`) or (any or non-gready-multiple) of these  -]~!@#%&*()_=+[{}'";<>/?|., then lookahead "," two or more comma
	      OR
	  start (:`) or (.`) or       any or multiple of these       -]~!@#%&*()_=+[{}'";<>/?.,
	*/

	/*middle part:
	  non-gready-multiple of everything   then multiple double comma ",,"
	*/

	/*last part:
	  :` or .`  or                                           -]~!@#%&*()_=+[{}'";<>/?|.,
	*/

	var splitWord splitableWord

	const part1 = ":`"
	const part2 = ".`"
	const symbols = "-]~!@#%&*()_=+[{}'\";<>/?|.,"

	var splitPrefix func(word *string)
	var splitSuffix func(word *string)

	splitPrefix = func(word *string) {
		if len(*word) == 0 {
			return
		}
		if strings.HasPrefix(*word, part1) {
			splitWord.begin += part1
			*word = (*word)[2:]
		} else if strings.HasPrefix(*word, part2) {
			splitWord.begin += part2
			*word = (*word)[2:]
		} else if strings.IndexAny(*word, symbols) == 0 {
			splitWord.begin += (*word)[0:1]
			*word = (*word)[1:]
		} else {
			return
		}
		splitPrefix(word)
	}

	splitSuffix = func(word *string) {
		if len(*word) == 0 {
			return
		}
		if strings.HasSuffix(*word, part1) {
			splitWord.end = part1 + splitWord.end
			*word = (*word)[0 : len(*word)-2]
		} else if strings.HasSuffix(*word, part2) {
			splitWord.end = part2 + splitWord.end
			*word = (*word)[0 : len(*word)-2]
		} else if lastChar := (*word)[len(*word)-1:]; strings.IndexAny(symbols, lastChar) != -1 {
			splitWord.end = lastChar + splitWord.end
			*word = (*word)[0 : len(*word)-1]
		} else {
			return
		}
		splitSuffix(word)
	}

	splitPrefix(&word)
	splitSuffix(&word)

	//TODO: Implement Split commas
	// splitComma(&word)

	splitWord.middle = word

	return splitWord
}

func (avro *SuggestionBuilder) sortByPhoneticRelevance(phonetic string, dictSuggestion []string) []string {
	// Copy things into a sortable interface implementation then call sort
	// TODO: BUG: Sorted result is inconsistant, bug in levenshtein or byLevenshtein?
	suggSlice := make([]string, len(dictSuggestion))
	copy(suggSlice, dictSuggestion)
	var sortAble byLevenshtein
	sortAble.baseWord = phonetic
	sortAble.words = &suggSlice

	sort.Sort(sortAble)

	return suggSlice
}

const karLetters = "ািীুূৃেৈোৌৄ"

func (avro *SuggestionBuilder) isKar(input rune) bool {
	// if len(input) < 1 {
	// 	return false
	// }
	return strings.ContainsRune(karLetters, input)
}

const vowelLetters = "অআইঈউঊঋএঐওঔঌৡািীুূৃেৈোৌ"

func (avro *SuggestionBuilder) isVowel(input rune) bool {
	// if len(input) < 1 {
	// 	return false
	// }
	return strings.ContainsRune(vowelLetters, input)
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
						cacheRightChar, _ := utf8.DecodeLastRuneInString(cacheItem)
						suffixLeftChar, _ := utf8.DecodeRuneInString(suffix)
						if avro.isVowel(cacheRightChar) && avro.isKar(suffixLeftChar) {
							fullWord = cacheItem + "\u09df" + suffix // \u09df = B_Y
							tempSlice = append(tempSlice, fullWord)
							avro.addToTempCache(fullWord, cacheItem, key)
						} else {
							if cacheRightChar == '\u09ce' { // \u09ce = b_Khandatta
								fullWord = avro.trimLastRune(cacheItem) + "\u09a4" + suffix // \u09a4 = b_T
								tempSlice = append(tempSlice, fullWord)
								avro.addToTempCache(fullWord, cacheItem, key)
							} else if cacheRightChar == '\u0982' { // \u0982 = b_Anushar
								fullWord = avro.trimLastRune(cacheItem) + "\u0999" + suffix // \u09a4 = b_NGA
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

func (avro *SuggestionBuilder) trimLastRune(str string) string {
	ustring := []rune(str)
	return string(ustring[0 : len(ustring)-1])
}

func (avro *SuggestionBuilder) getPreviousSelection(splitWord splitableWord, suggestionWords []string) int {
	word := splitWord.middle
	prevIndex, _, prevFound := avro.CandSelector.Get(word, suggestionWords)
	if !prevFound {
		//Full word was not found, try checking without suffix
		_len := len(word)
		if _len >= 2 {
			for j := 1; j < _len; j++ {
				testSuffix := strings.ToLower(word[_len-j:])
				suffix, ok := avro.SuffixDict[testSuffix]
				if ok {

					key := word[0 : len(word)-len(testSuffix)]
					_, _prevKeyWord, _prevFound := avro.CandSelector.Get(key, suggestionWords)

					if _prevFound {
						//Get possible words for key

						kwRightChar, _ := utf8.DecodeLastRuneInString(_prevKeyWord)
						suffixLeftChar, _ := utf8.DecodeRuneInString(suffix)

						selectedWord := ""

						if avro.isVowel(kwRightChar) && avro.isKar(suffixLeftChar) {
							selectedWord = _prevKeyWord + "\u09df" + suffix // \u09df = B_Y
						} else {
							if kwRightChar == '\u09ce' { // \u09ce = b_Khandatta
								selectedWord = avro.trimLastRune(_prevKeyWord) + "\u09a4" + suffix // \u09a4 = b_T
							} else if kwRightChar == '\u0982' { // \u0982 = b_Anushar
								selectedWord = avro.trimLastRune(_prevKeyWord) + "\u0999" + suffix // \u09a4 = b_NGA
							} else {
								selectedWord = _prevKeyWord + suffix
							}
						}

						for i, v := range suggestionWords {
							if v == selectedWord {
								//Save this referrence
								avro.CandSelector.Set(word, selectedWord)
								return i
							}
						}
					}
				}
			}
		}
	}
	return prevIndex
}

func (avro *SuggestionBuilder) joinSuggestion(autoCorrect correctableWord, dictSuggestion []string, phonetic string, splitWord splitableWord) Suggestion {
	var words []string

	if avro.Pref.DictDisabled {
		words = []string{splitWord.begin + phonetic + splitWord.end}
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

		//array_append_unique implemented by these two anonymous functions

		func() {
			if len(words) == 0 {
				words = sortedWords
				return
			}
			for i := range sortedWords {
				if sortedWords[i] == words[0] {
					words = sortedWords
					if i > 0 {
						words[i], words[0] = words[0], words[i]
					}
					return
				}
			}
			words = append(words, sortedWords...)
		}()

		/* 3rd Item: Classic Avro Phonetic */
		func() {
			for i := range words {
				if words[i] == phonetic {
					return
				}
			}
			words = append(words, phonetic)
		}()

		suggestion := Suggestion{}

		//Is there any previous custom selection of the user?
		suggestion.PrevSelection = avro.getPreviousSelection(splitWord, words)

		//Add padding to all, except exact autocorrect
		for i := range words {
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

func (avro *SuggestionBuilder) StringCommitted(word string, candidate string) {
	if !avro.Pref.DictDisabled {

		//If it is called, user made the final decision here
		//Check and save selection without suffix if that is not present

		cacheWord, ok := avro.tempCache[candidate]
		if ok {
			//Don't overwrite existing value
			if !avro.CandSelector.Has(cacheWord.eng) {
				avro.CandSelector.Set(cacheWord.eng, cacheWord.base)
			}
		} else {
			//This code deviates from JS implementation where THIS ELSE block's logic is handled seperately in `updateCandidateSelection` which is not implemented here
			avro.CandSelector.Set(word, candidate)
		}
		avro.CandSelector.Save()
	}
}
