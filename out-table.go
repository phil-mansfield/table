/*
Package table implements a simple ASCII table of floating point numbers.
*/
package table

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type HeaderState int

const (
	KeepHeader HeaderState = iota
	RemoveHeader
)

const (
	minColWidth = 13
)

type OutTable struct {
	colWidth int
	cols     int

	userHeader string
	colNames   string

	floatFmt string

	rows []string
}

// New creates a new table appropriately names columns.
func NewOutTable(colNames ...string) *OutTable {
	t := new(OutTable)
	t.colWidth = minColWidth

	for _, name := range colNames {
		if len(name) > t.colWidth {
			t.colWidth = len(name)
		}
	}

	t.cols = len(colNames)

	strFmt := make([]string, t.cols)
	floatFmt := make([]string, t.cols)

	for i := range colNames {
		strFmt[i] = fmt.Sprintf("%%%ds", t.colWidth)
		floatFmt[i] = fmt.Sprintf("%%%d.%dg", t.colWidth, t.colWidth - 5)
	}

	t.floatFmt = "  " + strings.Join(floatFmt, " ")

	is := make([]interface{}, len(colNames))
	for i := range colNames {
		is[i] = interface{}(colNames[i])
	}
	t.colNames = fmt.Sprintf(strings.Join(strFmt, " "), is...)

	t.rows = make([]string, 0)

	return t
}

// SetHeader adds a series of '#'-prefixed lines to be written in between
// the customary sed string and the column names at the head of the file.
// If you want to  have a completely custom header, make the first line a
// comment and then pass RemoveComment to Write and Print.
func (t *OutTable) SetHeader(header string) { t.userHeader = header + "\n" }

// AddRow adds a row of floats to the table.  This row must have the same
// length as the column names given to the last call to SetColumns.
func (t *OutTable) AddRow(xs ...float64) {
	is := make([]interface{}, len(xs))
	for i := range xs {
		is[i] = interface{}(xs[i])
	}

	t.rows = append(t.rows, fmt.Sprintf(t.floatFmt, is...))
}

func convertToComments(s string) []string {
	rows := strings.Split(s, "\n")
	for i, row := range rows {
		rows[i] = "# " + row
	}

	return rows
}

// Comment places s on the next line of the table, with each line prefixed
// by a '#'.
func (t *OutTable) Comment(s string) {
	comments := convertToComments(s)
	for _, comment := range comments {
		t.rows = append(t.rows, comment)
	}
}

// InlineComment adds s to the end of the last printed line, prefixed
// by a '#'.  If s has newlines in it, those lines are put after the
// current line and prefixed with '#'.  If there are no rows in the table
// yet, a call InlineComment is the same as a call to Comment.
func (t *OutTable) InlineComment(s string) {
	if len(t.rows) == 0 {
		t.rows = append(t.rows, "")
	} else {
		t.rows[len(t.rows)-1] += " "
	}

	comments := convertToComments(s)
	t.rows[len(t.rows)-1] += comments[0]
	for _, comment := range comments[1:len(comments)] {
		t.rows = append(t.rows, comment)
	}
}

func (t *OutTable) header() string {
	sedStr := "(use the command \"$ sed -e '/^#/d' -e 's/#.*$//' <filename>\" to remove comments)\n"
	rows := convertToComments(sedStr + t.userHeader + t.colNames)
	return strings.Join(rows, "\n")
}

// Write writes the contents of the table to the named file.  The file
// is created if it does not exist.  The header is not written if h is set
// to RemoveHeader.  This is useful if writing multiple tables to the same
// log file.
func (t *OutTable) Write(h HeaderState, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	// Only cowards check calls to Close.
	defer file.Close()

	w := bufio.NewWriter(file)

	switch h {
	case KeepHeader:
		fmt.Fprintln(w, t.header())
	}

	for _, row := range t.rows {
		fmt.Fprintln(w, row)
	}

	w.Flush()

	return nil
}

// Print prints the contents of the table to stdout.  The header is not
// printed if h is set to RemoveHeader.  This is useful if printing
// multiple tables at once.
func (t *OutTable) Print(h HeaderState) {
	switch h {
	case KeepHeader:
		fmt.Println(t.header())
	}

	for _, row := range t.rows {
		fmt.Println(row)
	}
}
