package kornet

import (
  "bytes"
  "io"
  "net/http"
  "net/url"
  "skfw/papaya/koala/kio/leaf"
  "time"
)

func KHttpClient(URL url.URL, request *Request, timeout time.Duration) (*Response, error) {

  body := request.Body
  defer body.Close()

  buff := bytes.NewReader(body.ReadAll())

  req, err := http.NewRequest(request.Method, URL.String(), buff)

  if err != nil {

    return nil, err
  }

  client := &http.Client{

    Timeout: timeout,
  }

  resp, err := client.Do(req)

  if err != nil {

    return nil, err
  }

  data, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  //resp.Header cannot convert into Keys, [][]string not []string

  return &Response{
    Header: KSafeSimpleHeaders(url.Values(resp.Header)),
    Body:   leaf.KMakeBuffer(data),
  }, nil
}
