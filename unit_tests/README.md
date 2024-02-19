# unit_tests/README.md

## TL'DR: running tests
To run all test files recursively found across the entire repository, you can do this in at the root:
```bash
$ go test ./...
```

To run just the `unit_tests` folder from the root:
```bash
$ go test ./unit_tests
```

> To run individual tests, you may see a "play" button next to test functions in VSCode (or get a Go extension). While you can `go test` individual files, this will **not** start executing from `TestMain()` which can be problematic.

On the default settings, any print statements triggered from tested code are not printed out and there are only package-scoped summaries of tests done. Use the verbose `-v` flag to see console output and what tests were run:
```bash
$ go test ./unit_tests -v
...

=== RUN   TestGetUserPoints
2024/02/18 22:55:20 usertable_user_id
2024/02/18 22:55:20 1
2024/02/18 22:55:20 finding
--- PASS: TestGetUserPoints (0.00s)

...

PASS
ok      slugquest.com/backend/unit_tests        0.698s
```

## Writing tests

### Test functions
All Go tests *must* be in this format:
* Begin with "Test"
* Have only one parameter, a `testing.T`
* Return nothing, not even an `error`
* Be in a file that ends with `_test.go`
```go
// lol_test.go

func TestLol(t *testing.T) {
    naur, err := lol()
    ...
}
```

### Setup/destruct for tests

In a package, there can be a single `TestMain()` declared that is ran before any tests. In `main_test.go`, I set this up to be establishing the in-memory DB connection and loading dummy data.

### Note:
1. From what I have seen, the unit tests are ran in sequential order in a file. **Keep this in mind when writing the logic for your tests.**

2. Instead of using `fmt.Print` or `log.Print`, use `t.Error` and `t.Errorf` to format messages since these will also **quit** the test at that point, and it will be considered a FAIL.
```bash
# Output sample of t.Error("lol") at the bottom of TestGetTaskId
--- FAIL: TestGetTaskId (0.00s)
    db_tasks_test.go:66: lol
```

