The `table` package is a library for reading ASCII tables willed with numbers.
This code has been bounced around between my HPC projects for the past ten
years or so and will work well on extremely large tables.

Basic usages
------------

Most of the time (i.e., you have a modest-sized table of ASCII floating point numbers),
the default configuration will be fine.

```Go
package	main

import (
    "fmt"
    "github.com/phil-mansfield/table"
)

func main() {
    t := table.TextFile("example_table.txt")
    cols := t.ReadFloat64s([]int{ 0, 2 }) // column indices
    fmt.Printf("1st column: %.2f\n", cols[0])
    fmt.Printf("3rd column: %.2f\n", cols[1])
}
```

There are corresponding methods `ReadFloat32s` and `ReadInts` which return
`[][]float32` and `[][]int` arrays, respectively.

Other data sources
------------------

You can also read tables which you read from stdin or which have
been stored as a string.
```Go
t1 := table.TextFile("example_table.txt")
t2 := table.Stdin()
t3 := table.Text(stringVariable) 
```

All three tables work identically, so we'll just focus on `TextFile()` below.

Configuration and Default Parameters
------------------------------------

Table parsing can be changed by modifying a default configuration file. This can
be done by passing a `TextConfig` struct as the second argument. The code snippet
below allows the table reader to parse CSVs, i.e., comma separated tables.
```Go
cfg := new(table.TextConfig)
copy(cfg, table.DefaultConfig)
cfg.Separator = ','
t := table.TextFile("example_table.txt", cfg)
```

(Don't tell anyone, sometimes it's easier to just modify the default config
struct directly)
```Go
table.DefaultConfig.Separator = ','
t := table.TextFile("example_table.txt")
```

The full set of configuration parameters and their default values are shown below

```Go
DefaultConfig = TextConfig{
	Separator: ' ',
	Comment: '#',
	SkipLines: 0,
	ColumnNames: map[string]int{},
	
	MaxBlockSize: 2 * 1<<30,
	MaxLineSize: 1<<20,
}
```

* `Separator` - The character used to separate columns within a single line.
* `Comment` - The character used to start comments. Comments (including full-line comments) will
  be removed before parsing.
* `SkipLines` - Sometimes tables don't use comment characters in their header. In these cases,
  set `SkipLines` to the number of header lines that need to be skipped.
* `ColumnNames` - For large tables that you read from frequently, you may want to name the columns.
  `ColumnNames` lets you map strings onto column indices. After this is set, 
* `MaxBlockSize` - Parameter used when reading very large tables (see below).
* `MaxLineSize` - Parameter used when reading very large tables (see below).

Reading Really Large Tables
---------------------------

This package is good at reading very large tables (e.g., if you have have a multi-terabyte halo
catalog represented as a CSV for some horrifying reason. But if you want to do this, you need to
do more work

```Go
t := table.TextFile("example_table.txt")
buf := make([][]float64, 2)	// One buffer for each column                    
blocks := t.Blocks() // Number of pieces the table has been broken into      
for	block := 0;	block < blocks; block++ {
    fmt.Printf("Block: %d\n", block)
    cols = t.ReadFloat64Block([]int{0, 2}, block, buf)
    // buf now contains target columns for this block. They will be          
    // replaced if this is called again, so either complete you              
    // calculations entirely before reading another block or append the      
    // data to another array.                                                
    fmt.Printf("1st column: %.2f\n", cols[0])
    fmt.Printf("3rd column: %.2f\n", cols[1])
}
```

The file will be broken into a series of blocks, each of which will be close to `MaxBlockSize` bytes
long (or the number of bytes required to get to the end of the file, whichever is smaller.) Lines will
not be broken in half, and to achieve this, all lines in the file must be smaller than `MaxLineSize`
characters.
