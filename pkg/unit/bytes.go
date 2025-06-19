package unit

import "fmt"

var (
	kb float64 = 1024
	mb         = kb * kb
	gb         = mb * kb
)

func FormatGB(d any) string {
	i, ok := d.(uint64)
	if !ok {
		return fmt.Sprintf("%v-%T", d, d)
	}
	f := float64(i)
	if f > gb {
		return fmt.Sprintf("%.0fGB", f/gb)
	}
	if f > mb {
		return fmt.Sprintf("%.0fMB", f/mb)
	}
	return fmt.Sprintf("%.0fKB", f/kb)
}
