package helper

import (
	"fmt"
	"reflect"
)

// Change represents a single detected difference between old and new data.
type Change struct {
	OldValue interface{} `json:"oldValue"`
	NewValue interface{} `json:"newValue"`
	Label    string      `json:"label"` // Identifies the field that changed
}

// GetChanges compares two maps and returns a slice containing all detected differences.
// Each element in the slice represents a field that has changed.
func GetChanges(updatedData, oldData map[string]interface{}) []Change {
	// Initialize a slice to hold all the change structs.
	// Pre-allocate capacity slightly for potential performance gain if many changes are expected.
	changeList := make([]Change, 0, len(updatedData)) // Start with 0 length, capacity based on updatedData size

	// Iterate through each key-value pair in the updated data map.
	for key, newVal := range updatedData {
		// Retrieve the corresponding value from the old data map.
		oldVal, exists := oldData[key]

		// Check if the key is new (!exists) or if the values are different.
		if !exists || !reflect.DeepEqual(newVal, oldVal) {
			// A change is detected. Create a Change struct.
			oldValStr := fmt.Sprintf("%v", oldVal)
			newValStr := fmt.Sprintf("%v", newVal)

			change := Change{
				OldValue: oldValStr,
				NewValue: newValStr,
				Label:    key,
			}
			// Append the detected change to the list.
			changeList = append(changeList, change)
		}
	}

	// Return the slice of changes directly.
	// It will be an empty slice if no changes were found.
	return changeList
}
