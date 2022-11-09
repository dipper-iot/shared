package util

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Merge(list1 []string, list2 []string) []string {
	list3 := list1
	for index := range list2 {
		list3 = append(list3, list2[index])
	}
	return list3
}
