package translation

func Translate(key string, dictionary ...Dictionary) string {
	return NewTranslator(dictionary...).Translate(key)
}

func mergeDictionary(defaultDictionary Dictionary, overwrites ...Dictionary) Dictionary {
	for _, overwrite := range overwrites {
		for key, value := range overwrite {
			defaultDictionary[key] = value
		}
	}

	return defaultDictionary
}
