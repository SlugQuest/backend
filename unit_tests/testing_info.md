# unit_tests

## Running tests
To run all tests across the entire repository, you can do this in at the root:
```bash
$ go test ./...
```

To run just the `unit_tests` folder from the root:
```bash
$ go test ./unit_tests
```

## Writing tests

All Go tests *must* be in this format:
* Begin with "Test"
* Have only one parameter, a `testing.T`
* Return nothing
```go
func TestName(t *testing.T) {
    ...
}
```

### Note:
1. From what I have seen, the unit tests are ran in sequential order in a file. **Keep this in mind when writing your tests.**

2. Instead of using `fmt.Print` or `log.Print`, use `t.Error` and `t.Errorf` to format messages.


