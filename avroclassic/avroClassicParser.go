package avroclassic

import (
	"strings"
    "unicode"
)

type Parser struct {
	Data AvroData
}

func (avro *Parser) Parse(input string) string {
	fixed := avro.FixString(input);
	output := "";
	for cur := 0; cur < len(fixed); cur++ {
		start := cur
        end := cur + 1
        prev := start - 1
		matched := false

		for _, pattern := range avro.Data.Patterns {
			end = cur + len(pattern.Find);
			if end <= len(fixed) && fixed[start:end] == pattern.Find {
				prev = start - 1;
				if len(pattern.Rules) > 0 {
				    for _, rule := range pattern.Rules {
						replace := true

						chk := 0

						for _, match := range rule.Matches {
						
							if match.Type == "suffix" {
								chk = end
							} else /* Prefix */ {
								chk = prev;
							}
						
							if match.Scope == "punctuation" { // Beginning
								if ! (
										((chk < 0) && (match.Type == "prefix")) || 
										((chk >= len(fixed)) && (match.Type == "suffix")) || 
										avro.isPunctuation(fixed[chk])) != match.Negative {
									replace = false
									break
								}
							} else if match.Scope == "vowel" { // Vowel
								if ! (
										((chk >= 0 && (match.Type == "prefix")) || 
											(chk < len(fixed) && (match.Type == "suffix"))) && 
                                            avro.isVowel(fixed[chk])) != match.Negative {
									replace = false
									break
								}
							} else if match.Scope == "consonant" { // Consonant
								if ! (
										(
											(chk >= 0 && (match.Type == "prefix")) || 
											(chk < len(fixed) && match.Type == "suffix")) && 
										avro.isConsonant(fixed[chk])) != match.Negative {
									replace = false
									break
								}
							} else if match.Scope == "exact" { // Exact
								var s, e int
								if match.Type == "suffix" {
									s = end
									e = end + len(match.Value)
								} else { // Prefix
									s = start - len(match.Value)
									e = start
								}
								if ! avro.isExact(match.Value, fixed, s, e, match.Negative) {
									replace = false
									break
								}
							}
						}
					
						if replace {
							output += rule.Replace
							cur = end - 1
							matched = true
							break
						}
					
					}
				}
				if matched == true {
				    break
				}
				
				// Default
				output += pattern.Replace;
				cur = end - 1
				matched = true
				break
			}
		}
		
		if !matched {
			output += string(fixed[cur])
		}
	}
	return output
}

func (avro *Parser) FixString(input string) string {
	fixed := ""
	for _, v := range input {
		cChar := string(v)
		if avro.isCaseSensitive(v) {
			fixed += cChar
		} else {
			fixed += strings.ToLower(cChar)
		}
	}
	return fixed
}

func (avro *Parser) isVowel(c byte) bool {
	return strings.ContainsRune(avro.Data.Vowel, unicode.ToLower(rune(c)))
}

func (avro *Parser) isConsonant(c byte) bool {
	return strings.ContainsRune(avro.Data.Consonant, unicode.ToLower(rune(c)))
}

func (avro *Parser) isPunctuation(c byte) bool {
	return (!(avro.isVowel(c) || avro.isConsonant(c)))
}

func (avro *Parser) isExact(needle string, heystack string, start int, end int, not bool) bool {
	return ((start >= 0 && end < len(heystack) && (heystack[start:end] == needle)) != not)
}

func (avro *Parser) isCaseSensitive(c rune) bool {
	return strings.ContainsRune(avro.Data.CaseSensitive, unicode.ToLower(c))
}
