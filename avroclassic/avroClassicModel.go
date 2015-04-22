package avroclassic

import "encoding/json"
import "bytes"

type AvroPatternRuleMatch struct {
	Type  string
	Scope string 
    Value string
    Negative bool
}

type AvroPatternRule struct {
	Matches []AvroPatternRuleMatch
	Replace string
}

type AvroPattern struct {
	Find    string
	Replace string
	Rules   []AvroPatternRule
}

type AvroData struct {
	Patterns      []AvroPattern
	Vowel         string
	Consonant     string
	CaseSensitive string
}

func (m *AvroPatternRuleMatch) UnmarshalJSON(data []byte) error {
    var aux struct {
        Type string
        Scope string
        Value string
    }
    
    dec := json.NewDecoder(bytes.NewReader(data))
    if err := dec.Decode(&aux); err != nil {
        return err
    }
    m.Type = aux.Type
    m.Value = aux.Value
    m.Scope = aux.Scope
    if aux.Scope[0] == '!' {
	    m.Negative = true
	    m.Scope = aux.Scope[1:]
	}
    return nil
}