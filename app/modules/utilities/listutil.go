package utilities

func UniqueItemsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	uniques := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				uniques = append(uniques, elem)
				unique[elem] = true
			}
		}
	}
	return uniques
}

func Contains(slice1 []string, slice2 []string) bool {
	slice1Uniques := UniqueItemsOf(slice1)
	slice2Uniques := UniqueItemsOf(slice2)
	for _, value2 := range slice2Uniques {
		var slice2ItemInSlice1 = false
		for _, value1 := range slice1Uniques {
			if value1 == value2 {
				slice2ItemInSlice1 = true
			}
		}
		if !slice2ItemInSlice1 {
			return false
		}
	}
	return true
}
