package convert

func Uint64ToPoint(data uint64) *uint64 {
	if data == 0 {
		return nil
	}
	return &data
}
