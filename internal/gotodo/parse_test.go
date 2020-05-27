package gotodo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCompleteValidMark(t *testing.T) {
	assert.Equal(t, true, isComplete("x"))
}

func TestIsCompleteNoMark(t *testing.T) {
	assert.Equal(t, false, isComplete("Hello, world."))
}

func TestIsCompleteInvalidMark(t *testing.T) {
	assert.Equal(t, false, isComplete("X"))
	assert.Equal(t, false, isComplete("y"))
	assert.Equal(t, false, isComplete("-"))
	assert.Equal(t, false, isComplete(" "))
	assert.Equal(t, false, isComplete("[x]"))
	assert.Equal(t, false, isComplete("x "))
	assert.Equal(t, false, isComplete(" x"))
}

func TestParsePriority(t *testing.T) {
	assert.Equal(t, 1, parsePriority("A"))
	assert.Equal(t, 2, parsePriority("B"))
	assert.Equal(t, 3, parsePriority("C"))
	assert.Equal(t, 27, parsePriority("AA"))
	assert.Equal(t, 28, parsePriority("AB"))
	assert.Equal(t, 29, parsePriority("AC"))
	assert.Equal(t, 53, parsePriority("BA"))
	assert.Equal(t, 54, parsePriority("BB"))
	assert.Equal(t, 55, parsePriority("BC"))
	assert.Equal(t, 2056, parsePriority("CAB"))
}

func TestIsPriority(t *testing.T) {
	assert.Equal(t, true, isPriority("(A)"))
	assert.Equal(t, true, isPriority("(A)"))
	assert.Equal(t, true, isPriority("(Aa)"))
	assert.Equal(t, true, isPriority("(aA)"))
	assert.Equal(t, false, isPriority("(A"))
	assert.Equal(t, false, isPriority("(AA"))
	assert.Equal(t, false, isPriority("A)"))
}

func TestParseProjectTags(t *testing.T) {
	args := []string{"hello,", "world!", "+project", "+test", "+test"}
	projects, contexts := parseTags(args)

	assert.Equal(t, 2, len(projects))
	assert.Equal(t, 0, len(contexts))

	_, p1 := projects["hello"]
	_, p2 := projects["world!"]
	_, p3 := projects["project"]
	_, p4 := projects["test"]

	assert.Equal(t, false, p1)
	assert.Equal(t, false, p2)
	assert.Equal(t, true, p3)
	assert.Equal(t, true, p4)
}

func TestParseContextTags(t *testing.T) {
	args := []string{"hello,", "world!", "@project", "@test", "@test"}
	projects, contexts := parseTags(args)

	assert.Equal(t, 2, len(contexts))
	assert.Equal(t, 0, len(projects))

	_, c1 := contexts["hello"]
	_, c2 := contexts["world!"]
	_, c3 := contexts["project"]
	_, c4 := contexts["test"]

	assert.Equal(t, false, c1)
	assert.Equal(t, false, c2)
	assert.Equal(t, true, c3)
	assert.Equal(t, true, c4)
}

func TestParseTags(t *testing.T) {
	args := []string{"hello,", "world!", "+project", "@test"}
	projects, contexts := parseTags(args)

	assert.Equal(t, 1, len(projects))
	assert.Equal(t, 1, len(contexts))

	_, p1 := projects["hello"]
	_, p2 := projects["world!"]
	_, p3 := projects["project"]
	_, p4 := projects["test"]
	_, c1 := contexts["hello"]
	_, c2 := contexts["world!"]
	_, c3 := contexts["project"]
	_, c4 := contexts["test"]

	assert.Equal(t, false, p1)
	assert.Equal(t, false, p2)
	assert.Equal(t, true, p3)
	assert.Equal(t, false, p4)
	assert.Equal(t, false, c1)
	assert.Equal(t, false, c2)
	assert.Equal(t, false, c3)
	assert.Equal(t, true, c4)

}

func TestUnparsePriority(t *testing.T) {
	assert.Equal(t, "A", unparsePriority(1))
	assert.Equal(t, "B", unparsePriority(2))
	assert.Equal(t, "C", unparsePriority(3))
	assert.Equal(t, "AA", unparsePriority(27))
	assert.Equal(t, "AB", unparsePriority(28))
	assert.Equal(t, "AC", unparsePriority(29))
	assert.Equal(t, "BA", unparsePriority(53))
	assert.Equal(t, "BB", unparsePriority(54))
	assert.Equal(t, "BC", unparsePriority(55))
	assert.Equal(t, "CAB", unparsePriority(2056))
}

func TestParseDate(t *testing.T) {
	assert.Equal(t, true, parseDate("2020-01-01").Valid)
	assert.Equal(t, true, parseDate("1969-12-31").Valid)
	assert.Equal(t, true, parseDate("2055-12-13").Valid)
}

func TestParseDateInvalid(t *testing.T) {
	assert.Equal(t, false, parseDate("2020/01/01").Valid)
	assert.Equal(t, false, parseDate("2020/1/1").Valid)
	assert.Equal(t, false, parseDate("2020-1-1").Valid)
	assert.Equal(t, false, parseDate("2019-13-31").Valid)
	assert.Equal(t, false, parseDate("2019-12-32").Valid)
	assert.Equal(t, false, parseDate("2006-01-02T15:04:05-0700").Valid)
}
