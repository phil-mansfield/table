package table

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

type InTable struct {
	xs []map[string] float64

	ColNames []string
	Rows int
}

// NewInTable creates a table that can be read in from a file.
func NewInTable(colNames ...string) *InTable {
	t := new(InTable)

	t.ColNames = make([]string, len(colNames))
	for i := 0; i < len(colNames); i++ {
		t.ColNames[i] = colNames[i]
	}

	t.xs = nil

	return t
}

// This is suboptimal.
func (t *InTable) Access(row int, colName string) float64 {
	if t.xs == nil {
		panic("Table not yet initialized.")
	}

	return t.xs[row][colName]
}

func isComment(line string) bool {
	return len(line) > 0 && line[0] == '#'
}

func uncomment(line string) string {
	return strings.Split(line, "#")[0]
}

func getColIndices(line string) []string {
	colIndices := make([]string, 0)
	for _, word := range strings.Split(line, " ") {
		if len(word) > 0 {
			colIndices = append(colIndices, word)
		}
	}
	return colIndices
}

func findStr(strs []string, target string) bool {
	for _, str := range strs {
		if str == target {
			return true
		}
	}
	return false
}

// This is suboptimal.
func checkColNames(colIndices []string, colNames []string) {
	for _, name := range colNames {
		if !findStr(colIndices, name) {
			panic(fmt.Sprintf("name '%s' not found in input file.", name))
		}
	}
}

// Read takes the contents of the specified file and converts them
// into an InTable.
func (t *InTable) Read(fileName string) {
	str, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(str), "\n")
	if len(lines) > 0  && !isComment(lines[0]) {
		panic("Input text file does not have table col names.")
	}

	headerEnded := false
	var colIndices []string
	t.xs = make([]map[string]float64, 0)

	for i, line := range lines {
		if !isComment(line) && len(line) > 0 {
			if !headerEnded {
				colIndices = getColIndices(lines[i - 1][1:])
				checkColNames(colIndices, t.ColNames)
				headerEnded = true
			}
			
			uncom := uncomment(line)
			words := make([]string, 0)
			for _, word := range strings.Split(uncom, " ") {
				if len(word) > 0 {
					words = append(words, word)
				}
			}

			if len(words) != len(colIndices) {
				panic(fmt.Sprintf("line %d has the incorrect " +
					"number of columns", i + 1))
			}

			m := make(map[string]float64)
			for col, word := range words {
				m[colIndices[col]], err = strconv.ParseFloat(word, 64)
				if err != nil {
					panic(err)
				}
			}

			t.xs = append(t.xs, m)
		}
	}

	t.Rows = len(t.xs)
}
