# Simple calculator
### Building
Building requires the go package
```shell
go mod tidy
go run cli_calc
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