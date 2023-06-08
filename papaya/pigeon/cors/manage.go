package cors

import (
  "errors"
  "net/url"
)

type ManageConsumers struct {

  // list data consumer
  consumers []ConsumerImpl
}

type ManageConsumersImpl interface {
  Init() error
  Add(method string, origin string) error
  Get(method string, origin string) ConsumerImpl
  Grant(origin string) error

  //Header(method string, origin string) (*http.Header, error) // Deprecated
}

func ManageConsumersNew() (ManageConsumersImpl, error) {

  manageConsumers := &ManageConsumers{}
  if err := manageConsumers.Init(); err != nil {

    return nil, err
  }
  return manageConsumers, nil
}

func (c *ManageConsumers) Init() error {

  c.consumers = make([]ConsumerImpl, 0)

  return nil
}

func (c *ManageConsumers) Add(method string, origin string) error {

  var err error
  var URL *url.URL
  var consumer ConsumerImpl

  if origin != "" {

    // update consumer
    if consumer = c.Get("*", origin); consumer != nil {

      if !consumer.AcceptMethod(method) {

        return errors.New("unable to add method " + method + " in " + origin)
      }

      return nil
    }

    URL = nil

    // if origin asterisk, URL set NULL
    if origin != "*" {

      // create new consumer
      URL, err = url.Parse(origin)
      if err != nil {

        return errors.New("unable to parse url")
      }
    }

    // allow all methods
    methods := make([]string, 0)
    headers := make([]string, 0)

    methods = append(methods, method) // append new method
    headers = append(headers, "Authorization")

    consumer, err = ConsumerNew(URL, methods, headers, false, 0)
    if err != nil {

      return errors.New("unable to create consumer")
    }

    c.consumers = append(c.consumers, consumer)
  }

  return nil
}

func (c *ManageConsumers) Get(method string, origin string) ConsumerImpl {

  for _, consumer := range c.consumers {

    if consumer.Check(method, origin) {

      return consumer
    }
  }

  return nil
}

// deprecated method

//func (c *ManageConsumers) Header(method string, origin string) (*http.Header, error) {
//
//  if method != "" && origin != "" {
//
//    var headers []string
//    headers = make([]string, 0)
//
//    if consumer := c.Get(method, origin); consumer != nil {
//
//      return consumer.Header(origin, method, headers)
//    }
//
//    return nil, errors.New("header not found")
//  }
//
//  return nil, errors.New("undefined method or origin")
//}

func (c *ManageConsumers) Grant(origin string) error {

  if origin != "" {

    // update consumer
    if consumer := c.Get("*", origin); consumer != nil {

      // set all methods
      for _, method := range Methods {

        if !consumer.AcceptMethod(method) {

          return errors.New("unable to add method " + method + " in " + origin)
        }
      }

      return nil
    }

    // create new consumer and update consumer
    for _, method := range Methods {

      if err := c.Add(method, origin); err != nil {

        return err
      }
    }
  }

  return nil
}
