package translation

import "fmt"

type Translator struct {
	dictionary map[string]string
}

func (t Translator) Translate(key string) string {
	return t.dictionary[key]
}

func (t Translator) TranslateWithParams(key string, params ...string) string {
	return fmt.Sprintf(t.dictionary[key], params)
}

func NewTranslator(dictionary ...Dictionary) Translator {
	return Translator{
		dictionary: mergeDictionary(defaultDictionary, dictionary...),
	}
}
