package utilu

// Определяет, содержится ли строка в слайсе
func Contains(s []string, val string) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}
	return false
}
