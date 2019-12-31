package typemap

func pathEqual(path1, path2 []int32) bool {
	if len(path1) != len(path2) {
		return false
	}
	for i, v := range path1 {
		if path2[i] != v {
			return false
		}
	}
	return true
}

