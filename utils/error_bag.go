package utils

import "github.com/kamva/go-libs/contracts"

type ErrorBag struct {
	errors   map[string][]string
	taggable contracts.Taggable
}

func (b ErrorBag) Append(field string, value string) {
	fieldName := b.taggable.GetTag(b.taggable ,field, "json")
	b.errors[fieldName] = append(b.errors[fieldName], value)
}
func (b ErrorBag) GetErrors() interface{} {
	return b.errors
}

func NewErrorBag(taggable contracts.Taggable) ErrorBag {
	return ErrorBag{errors: make(map[string][]string), taggable: taggable}
}
