package translation

import (
	"net/url"
	"fmt"
	"github.com/kamva/go-libs/contracts"
)

type ValidationTranslator struct {
	dictionary Dictionary
	taggable   contracts.Taggable
}

func (t ValidationTranslator) Translate(field string, tag string, params string) string {
	query := t.taggable.GetTag(field, "translation")
	mapping, _ := url.ParseQuery(query)
	key := mapping[tag]

	if params != "" {
		return fmt.Sprintf(t.dictionary[key[0]], params)
	}

	if key != nil {
		return t.dictionary[key[0]]
	} else {
		return fmt.Sprintf("%s validation failed", field)
	}
}

func NewValidationTranslator(taggable contracts.Taggable, dictionary Dictionary) ValidationTranslator {
	return ValidationTranslator{
		taggable:   taggable,
		dictionary: dictionary,
	}
}
