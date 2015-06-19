package table

import (
	"testing"
	"fmt"
)

func TestReadTable(t *testing.T) {
	cols, err := ReadTable("example_table.txt", []int{ 3, 1 }, nil)
	if err != nil { panic(err.Error()) }
	fmt.Println(cols)
}
