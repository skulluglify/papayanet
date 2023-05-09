# Kornet, Koala Repeater Network

### try parsing request body
- application/json
- application/xml
- multipart/form-data

### dependencies
- fasthttp

### Standard Kornet Result

```ts
result = {
    id: number, // RPC identify
    logs: array<string>, // logs history
    status: string, // status API
    message: string, // message received
    error: boolean, // error Bad Request Or Internal Server
    data: object | array | null, // result data
}
```