package assert_test

import (
	"reflect"
	"testing"

	"github.com/alcortesm/queue/assert"
)

func TestSeq(t *testing.T) {
	for _, test := range []struct {
		from, to int
		expected []int
	}{
		{0, 0, []int{}},
		{0, 1, []int{0}},
		{1, 0, []int{0}},
		{0, 2, []int{0, 1}},
		{2, 0, []int{0, 1}},
		{-2, 2, []int{-2, -1, 0, 1}},
		{2, -2, []int{-2, -1, 0, 1}},
	} {
		obtained := assert.Seq(test.from, test.to)
		if !reflect.DeepEqual(obtained, test.expected) {
			t.Errorf("from %d to %d: expected %v, obtained %v",
				test.from, test.to, test.expected, obtained)
		}
	}
}
