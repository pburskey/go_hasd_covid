package utility

func SliceRangeContainsOnlyEmpties(aSlice []string) bool {
	if aSlice == nil || len(aSlice) == 0 {
		return false
	}
	var hasEmpty bool = true
	for i := 0; hasEmpty && i < len(aSlice); i++ {
		aString := aSlice[i]
		if len(aString) > 0 {
			hasEmpty = false
		}
	}

	return hasEmpty
}

func SliceRangeContainsNonEmptyValue(aSlice []string) bool {
	if aSlice == nil || len(aSlice) == 0 {
		return false
	}
	var itDoes bool = false
	for i := 0; !itDoes && i < len(aSlice); i++ {
		aString := aSlice[i]
		if len(aString) > 0 {
			itDoes = true
		} else {
			itDoes = false
		}
	}

	return itDoes
}

func CountEmptyValuesIn(aSlice []string) int {
	if aSlice == nil || len(aSlice) == 0 {
		return 0
	}
	var count int = 0
	for _, aString := range aSlice {

		if len(aString) == 0 {
			count++
		}
	}

	return count
}

func CountNonEmptyValuesIn(aSlice []string) int {
	if aSlice == nil || len(aSlice) == 0 {
		return 0
	}
	var count int = 0
	for _, aString := range aSlice {

		if len(aString) > 0 {
			count++
		}
	}

	return count
}
