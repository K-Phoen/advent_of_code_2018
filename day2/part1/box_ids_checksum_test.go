package main

import (
	"testing"
)

type textSample struct {
	boxId           string
	hasTwoLetters   bool
	hasThreeLetters bool
}

func TestHasNSimilarLetters(t *testing.T) {
	samples := []textSample{
		textSample{"abcdef", false, false},
		textSample{"bababc", true, true},
		textSample{"abbcde", true, false},
		textSample{"abcccd", false, true},
		textSample{"aabcdd", true, false},
		textSample{"abcdee", true, false},
		textSample{"ababab", false, true},
	}

	for _, sample := range samples {
		hasTwoLetters := hasNSimilarLetters(sample.boxId, 2)
		hasThreeLetters := hasNSimilarLetters(sample.boxId, 3)

		if hasTwoLetters != sample.hasTwoLetters {
			if sample.hasTwoLetters {
				t.Errorf("Expected two similar letters to be found in '%s'", sample.boxId)
			} else {
				t.Errorf("Did NOT expect two similar letters to be found in '%s'", sample.boxId)
			}
		}

		if hasThreeLetters != sample.hasThreeLetters {
			if sample.hasThreeLetters {
				t.Errorf("Expected three similar letters to be found in '%s'", sample.boxId)
			} else {
				t.Errorf("Did NOT expect three similar letters to be found in '%s'", sample.boxId)
			}
		}
	}
}
