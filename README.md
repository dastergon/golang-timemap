# golang-timemap
A thread-safe time-based key-value store for Go.

## Documentation
Full docs are available on [Godoc](http://godoc.org/github.com/dastergon/golang-timemap)

## Example
```go
func main() {
	l := timemap.New()
	l.Set("name", "Rustacean", time.Now())
	time.Sleep(10 * time.Millisecond)
	l.Set("name", "Gopher", time.Now())
	v, _ := l.Get("name", time.Now())
	fmt.Println(v) # Expect "Gopher" :-)
}
```
