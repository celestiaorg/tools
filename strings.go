package tools

func UniqueStringSlice(input []string) []string {

	strMap := map[string]bool{}
	for _, s := range input {
		strMap[s] = true
	}

	output := []string{}
	for s := range strMap {
		output = append(output, s)
	}

	return output
}
