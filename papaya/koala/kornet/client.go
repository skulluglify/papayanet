package kornet

import (
  "PapayaNet/papaya/koala/io/leaf"
  "bytes"
  "io"
  "net/http"
  "net/url"
  "time"
)

func KHttpClient(URL url.URL, request *KRequest, timeout time.Duration) (*KResponse, error) {

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

  params, err := url.ParseQuery(resp.Request.URL.RawQuery)

  if err != nil {

    return nil, err
  }

  data, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  return &KResponse{
    Headers: KHeaders(resp.Header),
    Params:  KParams(params),
    Body:    leaf.KMakeBuffer(data),
  }, nil
}
