// Package cache handles everything that has to do with the generation of the cache for the Search function.
package cache

// <---------------------------------------------------------------------------------------------------->

import (
	"strings"
)

// <---------------------------------------------------------------------------------------------------->

// Be aware that the following is slightly unhinged and insane for a minimal gain

// Encode takes in a string and return an 8 byte array. The array should be viewed as a 64 long bit chain.
// The first 60 bit depending on if they're flipped or not indecate, whether a certain character is inside of the origin string (at least once).
// The 60 characters are all ascii chars (except for upper case letters).
// This allows us to simply compare two byte arrays on if a string has all the characters as needded for the search string later one,
// if that is not the case we can just skip that string and save having to do a full sub string search
func Encode(input string) [8]byte {
	// charMap with all relevant chars and byte position and increase value for bit flip
	charMap := map[rune][]uint8{
		'a': {0, 1}, 'b': {0, 2}, 'c': {0, 4}, 'd': {0, 8}, 'e': {0, 16},
		'f': {0, 32}, 'g': {0, 64}, 'h': {0, 128}, 'i': {1, 1}, 'j': {1, 2},
		'k': {1, 4}, 'l': {1, 8}, 'm': {1, 16}, 'n': {1, 32}, 'o': {1, 64},
		'p': {1, 128}, 'q': {2, 1}, 'r': {2, 2}, 's': {2, 4}, 't': {2, 8},
		'u': {2, 16}, 'v': {2, 32}, 'w': {2, 64}, 'x': {2, 128}, 'y': {3, 1},
		'z': {3, 2}, '0': {3, 4}, '1': {3, 8}, '2': {3, 16}, '3': {3, 32},
		'4': {3, 64}, '5': {3, 128}, '6': {4, 1}, '7': {4, 2}, '8': {4, 4},
		'9': {4, 8}, '!': {4, 16}, '#': {4, 32}, '$': {4, 64}, '%': {4, 128},
		'&': {5, 1}, '\'': {5, 2}, '(': {5, 4}, ')': {5, 8}, '+': {5, 16},
		',': {5, 32}, '-': {5, 64}, '.': {5, 128}, ';': {6, 1}, '=': {6, 2},
		'@': {6, 4}, '[': {6, 8}, ']': {6, 16}, '^': {6, 32}, '_': {6, 64},
		'`': {6, 128}, '{': {7, 1}, '}': {7, 2}, '~': {7, 4},
	}
	output := [8]byte{}

	// loop over the chars of the input string
	for _, char := range strings.ToLower(input) {
		// check if the char is (still) inside our charMap
		if bitFlipValue, ok := charMap[char]; ok {
			// if the char was still inside bit flip the output at the correct position
			output[bitFlipValue[0]] += bitFlipValue[1]
			// remove the char from the charMap to not accidentally add the bit flip value again
			delete(charMap, char)
		}
	}

	return output
}

// CompareBytes checks if all required letters form the search string are inside the searched string
func CompareBytes(searchBytes [8]byte, compareBytes [8]byte) bool {
	// Check if all flipped bits in searchBytes are also flipped in compareBytes
	for index := range searchBytes {
		if searchBytes[index]&^compareBytes[index] != 0 {
			// if we ever reach here it means that one of relevant bits wasn't flipped
			return false
		}
	}

	return true
}
