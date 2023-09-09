package main

func GetValues(m map[string]any) []any {
	values := make([]any, 0, len(m))

	for _, v := range m {
		values = append(values, v)
	}

	return values
}
