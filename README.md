# Simple calculator
### Building
Building requires the go package
```shell
go mod tidy
go run cli_calc
```
To run in debug mode, add the -d flag
```shell
go run cli_calc -d
```

### Examples

```
5+6
>11
```

```
5+6/0
>Cannot divide by zero
```

```
rm -rf ~
>Invalid expression. Please try again
```
