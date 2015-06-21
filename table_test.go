package table

import (
	"testing"
	"fmt"
)

func TestReadTable(t *testing.T) {
	cols, err := ReadTable("example_table.txt", []int{ 2, 0 }, nil)
	if err != nil { panic(err.Error()) }
	fmt.Println(cols)
}
