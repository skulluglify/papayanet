# Proposals

### SwagContext

```go

return ctx.Message("pong")
return ctx.Ok(Map{}).Message("pong")

return ctx.Ok(Map{})
return ctx.Created(Map{})
return ctx.Bad(Map{})
return ctx.Error(Map{})

res.Bytes()
```

```js
response = {
  "logs": [], // logger activities
  "data": {}, // data custom result
  "message": "", // current message
  "error": true // http.StatusCode
}
```

RPC - feature
```js
data = [
    {
        "id": 1,
        "data": null,
        "message": "queue id 1"
    },
    {
        "id": 2,
        "data": null,
        "message": "queue id 2"
    }
]
```