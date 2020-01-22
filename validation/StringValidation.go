package validation

func IsEmpty(data string) bool{
	if len(data) <= 0 {
		return true
	} else {
		return false
	}
}
