package util

import "samhofi.us/x/keybase/v2/types/chat1"

// ConvIDInSlice searches a slice for a given chat1.ConvIDStr and returns true if found
func ConvIDInSlice(needle chat1.ConvIDStr, haystack []chat1.ConvIDStr) bool {
	for _, c := range haystack {
		if c == needle {
			return true
		}
	}
	return false
}
