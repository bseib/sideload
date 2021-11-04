package util

func MapOfString(collection []string, f func(string) string) []string {
	newCollection := make([]string, len(collection))
	for i, v := range collection {
		newCollection[i] = f(v)
	}
	return newCollection
}
