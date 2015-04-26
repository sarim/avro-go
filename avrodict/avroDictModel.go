package avrodict

import (
    "bytes"
    "encoding/json"
    "math"
)

type AvroTable map[string][][][]rune

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
        (*t)[k] = chunkArray(rList)
    }
    return nil
}

func chunkArray(array [][]rune) [][][]rune {
    chunkSize := int(math.Ceil(float64(len(array)) / float64(30)))
	numOfChunks := int(math.Ceil(float64(len(array)) / float64(chunkSize)))
	output := make([][][]rune, numOfChunks)

	for i := 0; i < numOfChunks; i++ {
		start := i * chunkSize
		var length int
		x := len(array) - start
		if x < chunkSize {
			length = x
		} else {
			length = chunkSize
		}

		temp := make([][]rune, length)

		copy(temp[0:length], array[start:start+length])

		output[i] = temp
	}

	return output
}