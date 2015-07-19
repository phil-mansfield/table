package table

import (
	"runtime"

	"strings"
	"fmt"
	"io/ioutil"
	"strconv"
)

type ReadTableOptions struct {
	Delimiter rune
	Comment rune
}

var DefaultReadTableOptions = ReadTableOptions{
	Delimiter: ' ',
	Comment: '#',
}

func ReadTable(
	fileName string, colIdxs []int,
	opt *ReadTableOptions,
) ([][]float64, error) {
	if opt == nil { opt = &DefaultReadTableOptions }
	delim := opt.Delimiter
	comm := opt.Comment

	bs, err := ioutil.ReadFile(fileName)
	if err != nil { return nil, err }

	str := string(bs)
	runtime.GC()
	
	lines := make([]string, strings.Count(str, "\n"))
	n := splitInPlace(str, '\n', lines)
	lines = lines[:n]

	cols := make([][]float64, len(colIdxs))
	for i := range cols { cols[i] = make([]float64, len(lines)) }

	colLine := 0
	p := newParser(lines, cols, colIdxs, delim, comm)
	for i := range lines {
		if fullLine, err := p.parseLine(i, colLine); err != nil {
			return nil, err
		} else if fullLine {
			colLine++
		}
	}

	for i := range cols { cols[i] = cols[i][0: colLine] }	
	return cols, nil
}

func newParser(
	lines []string,
	cols [][]float64, colIdxs []int,
	delim, comm rune,
) *parser {
	maxCol := 0
	for _, col := range colIdxs {
		if col > maxCol { maxCol = col }
	}

	p := &parser{
		lines: lines, cols: cols, colIdxs: colIdxs,
		delim: uint8(delim), comm: uint8(comm),
		tokens: make([]string, maxCol + 1),
	}

	return p
}

type parser struct {
	lines, tokens []string
	cols [][]float64
	colIdxs []int
	delim, comm uint8
}

func (p *parser) parseLine(stringLine, floatLine int) (bool, error) {
	line := p.lines[stringLine]
	lineEnd := 0
	lline := len(line)
	for i := 0; i < lline; i++ {
		if line[i] == p.comm { break }
		lineEnd++
	}
	
	p.tokens = p.tokens[:cap(p.tokens)]
	n := splitInPlace(line[0: lineEnd], p.delim, p.tokens)
	p.tokens = p.tokens[:n]

	if len(p.tokens) == 0 { return false, nil }

	if len(p.tokens) < cap(p.tokens) { 
		return false, fmt.Errorf (
			"Line %d of source file contains %d columns, but was expecting %d.",
			stringLine, len(p.tokens), cap(p.tokens),
		)
	}

	for i, j := range p.colIdxs {
		t := p.tokens[j]
		var err error
		if p.cols[i][floatLine], err = strconv.ParseFloat(t, 64); err != nil {
			return false, fmt.Errorf(
				"On line %d of souce file: %s", stringLine, err.Error(),
			)
		}
	}

	return true, nil
}

func splitInPlace(s string, sep uint8, out []string) int {
	start := 0
	na := 0
	
	ls := len(s)
	lout := len(out)

	// i -> index into 
	for i := 0; i < ls; i++ {
		if na == lout { break }

		if s[i] == sep {
			if start != i {
				out[na] = s[start : i]
				na++
			}
			start = i + 1
		}
	}
	if na != lout && start < ls {
		out[na] = s[start:]
		na++
	}
	return na
}
