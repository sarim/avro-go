package avrophonetic

type CandidateSelector interface {
	Load() error
	Save() error
	Get(candidate string, suggestions []string) (int, string, bool)
	Set(candidate string, suggestion string) error
}

type InMemoryCandidateSelector struct {
	storage map[string]string
}

func (cs *InMemoryCandidateSelector) Load() error {
	cs.storage = make(map[string]string)
	return nil
}

func (cs *InMemoryCandidateSelector) Save() error {
	return nil
}

func (cs *InMemoryCandidateSelector) Get(candidate string, suggestions []string) (int, string, bool) {
	prev, ok := cs.storage[candidate]
	if ok {
		for i, v := range suggestions {
			if v == prev {
				return i, v, true
			}
		}
	}
	return 0, suggestions[0], false
}

func (cs *InMemoryCandidateSelector) Set(candidate string, suggestion string) error {
	cs.storage[candidate] = suggestion
	return nil
}

func NewInMemoryCandidateSelector(data map[string]string) *InMemoryCandidateSelector {
	s := &InMemoryCandidateSelector{}
	if data == nil {
		s.Load()
	} else {
		s.storage = data
	}
	return s
}
