```go
func plus(a ...int) int {
	var sum int = 0
	for _, item := range(a) {
		sum += item
	}
	return sum
}
```
