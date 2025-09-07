package slice

func TransformSliceToMap(slice []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range slice {
		result[item] = struct{}{}
	}
	return result
}

func TransformMapToSlice(m map[string]struct{}) []string {
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}

func TransformSliceToMultipleSlices[T any](slice []T, n int) [][]T {
	if n <= 0 {
		return nil
	}
	var result [][]T
	for i := 0; i < len(slice); i += n {
		end := i + n
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}
	return result
}

func MultipliesSlice[T any](slice []T, n int) []T {
	if n <= 0 {
		return nil
	}
	result := make([]T, 0, len(slice)*n)
	for _, item := range slice {
		for i := 0; i < n; i++ {
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
