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
temp := KISAny(a, b, c)
```

```go
var temp any
if c {
    temp = b
} else {
    temp = c
}
```