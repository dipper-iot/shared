package convert

func StringToPoint(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
