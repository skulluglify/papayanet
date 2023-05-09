# New Proposals

- Passing Indirection Value
- Inline Statement
- Coalescing Operator
- References
- Noop / Undefined

### *Inline Statement*

```js
temp = a ? b : c
```

```go
temp := pp.Lany(a, b, c)
```

```go
var temp any
if c {
    temp = b
} else {
    temp = c
}
```

### *Coalescing Operator*

```js
temp = a ?? b
```

```go
temp := pp.Qany(a, b)
```

```go
var temp any
if a != nil {
    temp = a
} else {
    temp = b
}
```