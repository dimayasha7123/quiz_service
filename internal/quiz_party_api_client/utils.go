package quizApi

func boolStrOrNilToBool(s string) bool {
	if s == "true" {
		return true
	}
	return false
}
