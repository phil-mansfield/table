package table

import (
	"fmt"
	"io/ioutil"
	"strings"
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
	
	lines := strings.Split(str, "\n")
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
	p := &parser{
		lines: lines, cols: cols, colIdxs: colIdxs,
		delim: string(delim), comm: comm,
	}

	for _, col := range colIdxs {
		if col > p.maxCol { p.maxCol = col }
	}

	return p
}

type parser struct {
	lines []string
	cols [][]float64
	colIdxs []int
	maxCol int
	delim string
	comm rune
}

func (p *parser) parseLine(stringLine, floatLine int) (bool, error) {
	line := p.lines[stringLine]
	lineEnd := 0
	for _, r := range line {
		if r == p.comm { break }
		lineEnd++
	}
	
	tokens := strings.Split(line[0: lineEnd], " ")
	dst := 0
	for src, tok := range tokens {
		if tok != "" {
			tokens[dst] = tokens[src]
			dst++
		}
	}
	tokens = tokens[0: dst]
	
	if len(tokens) == 0 { return false, nil }

	if len(tokens) <= p.maxCol { 
		return false, fmt.Errorf (
			"Line %d of source file comtains %d columns, but was expecting %d.",
			stringLine, len(tokens), p.maxCol + 1,
		)
	}

	for i, j := range p.colIdxs {
		t := tokens[j]
		var err error
		if p.cols[i][floatLine], err = strconv.ParseFloat(t, 64); err != nil {
			return false, fmt.Errorf(
				"On line %d of souce file: %s", stringLine, err.Error(),
			)
		}
	}

	return true, nil
}
