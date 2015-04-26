package avrodict

import (
    "bytes"
    "encoding/json"
)

type AvroTable map[string][][]rune

func (t *AvroTable) UnmarshalJSON(data []byte) error {
    var aux map[string][]string
    
    dec := json.NewDecoder(bytes.NewReader(data))
    if err := dec.Decode(&aux); err != nil {
        return err
    }
    for k, v := range aux {
        rList := make([][]rune, len(v))
        for i, w := range v {
            rList[i] = []rune(w)
        }
        (*t)[k] = rList
    }
    return nil
}