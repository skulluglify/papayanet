package kornet

import (
  "PapayaNet/papaya/koala/kio/leaf"
  m "PapayaNet/papaya/koala/mapping"
  "net/url"
)

// boolean, number, string

type Query map[string][]any
type QueryImpl interface {
  Get(name string) any
  GetAll(name string) []any
  Set(name string, value any)
  Del(name string)
}

func (q *Query) Get(name string) any {

  query := *q

  for k, v := range query {

    if k == name {

      if len(v) > 0 {

        return v[0]
      }

      return nil
    }
  }

  return nil
}

func (q *Query) GetAll(name string) []any {

  query := *q

  for k, v := range query {

    if k == name {

      return v
    }
  }

  return nil
}

func (q *Query) Set(name string, value any) {

  query := *q
  qq, ok := query[name]

  if ok {

    query[name] = append(qq, value)

  } else {

    query[name] = []any{value}
  }
}

func (q *Query) Del(name string) {

  query := *q

  delete(query, name)
}

type Request struct {
  Method string           `json:"method,omitempty"`
  URL    url.URL          `json:"url"`
  Header m.KMapImpl       `json:"header,omitempty"`
  Path   m.KMapImpl       `json:"header,omitempty"`
  Query  Query            `json:"query,omitempty"`
  Body   leaf.KBufferImpl `json:"body,omitempty"`
}

type Response struct {
  Header m.KMapImpl       `json:"header,omitempty"`
  Body   leaf.KBufferImpl `json:"body,omitempty"`
}
