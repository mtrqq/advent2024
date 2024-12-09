package golang

const (
	EndOfLine = "\n"
	Space     = " "
	Pipe      = "|"
	Comma     = ","
	Decimal   = 10
	BitSize64 = 64
)

func ArrayWithoutItem(array []int64, index int) []int64 {
	arrayWithout := make([]int64, index, len(array)-1)
	copy(arrayWithout, array[:index])
	return append(arrayWithout, array[index+1:]...)
}
