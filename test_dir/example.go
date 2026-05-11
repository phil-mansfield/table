package main

import (
	"fmt"
	"github.com/phil-mansfield/table"
)

func main() {
	t := table.TextFile("../example_table.txt")
	cols := t.ReadFloat64s([]int{ 0, 2 })
	fmt.Printf("1st column: %.2f\n", cols[0])
	fmt.Printf("3rd column: %.2f\n", cols[1])

	// Setup in the case where the file is super huge.
	t = table.TextFile("../example_table.txt")
	buf := make([][]float64, 2) // One buffer for each column
	blocks := t.Blocks() // Number of pieces the table has been broken into
	for block := 0; block < blocks; block++ {
		fmt.Printf("Block: %d\n", block)
		cols = t.ReadFloat64Block([]int{0, 2}, block, buf)
		// buf now contains target columns for this block. They will be
		// replaced if this is called again, so either complete you
		// calculations entirely before reading another block or append the
		// data to another array.
		fmt.Printf("1st column: %.2f\n", cols[0])
		fmt.Printf("3rd column: %.2f\n", cols[1])
	}
}
