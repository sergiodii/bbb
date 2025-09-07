package slice

import (
	"reflect"
	"sort"
	"testing"
)

func TestSlice(t *testing.T) {
	t.Run("TransformSliceToMap", func(t *testing.T) {
		// Arrange
		slice := []string{"a", "b", "c", "a"}

		// Act
		result := TransformSliceToMap(slice)

		// Assert
		expected := map[string]struct{}{
			"a": {},
			"b": {},
			"c": {},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("TransformMapToSlice", func(t *testing.T) {
		// Arrange
		m := map[string]struct{}{
			"a": {},
			"b": {},
			"c": {},
		}

		// Act
		result := TransformMapToSlice(m)

		// Assert
		expected := []string{"a", "b", "c"}
		sort.Strings(result)
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("TransformSliceToMultipleSlices", func(t *testing.T) {
		// Arrange
		slice := []int{1, 2, 3, 4, 5}

		// Act
		result := TransformSliceToMultipleSlices(slice, 2)

		// Assert
		expected := [][]int{
			{1, 2},
			{3, 4},
			{5},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("MultipliesSlice", func(t *testing.T) {
		// Arrange
		slice := []string{"a", "b"}

		// Act
		result := MultipliesSlice(slice, 3)

		// Assert
		expected := []string{"a", "a", "a", "b", "b", "b"}
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})

	t.Run("RemoveDuplicates", func(t *testing.T) {
		// Arrange
		slice := []string{"a", "b", "a", "c", "b"}

		// Act
		result := RemoveDuplicates(slice)

		// Assert
		expected := []string{"a", "b", "c"}
		sort.Strings(result)
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})
}
