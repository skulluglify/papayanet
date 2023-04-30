package cors

import (
  "errors"
  "net/http"
  "net/url"
)

type ManageConsumers struct {
  consumers []ConsumerImpl
}

type ManageConsumersImpl interface {
  Init() error
  Add(method string, origin string) error
  Get(method string, origin string) ConsumerImpl
  GetByOrigin(origin string) ConsumerImpl
  Check(method string, origin string) *http.Header
  GrantAll(origin string) error
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

  // update consumer
  if consumer := c.GetByOrigin(origin); consumer != nil {

    if !consumer.AcceptMethod(method) {

      return errors.New("unable to add method " + method + " in " + origin)
    }

    return nil
  }

  // create new consumer
  URL, err := url.Parse(origin)
  if err != nil {

    return errors.New("unable to parse url")
  }

  // allow all methods
  methods := make([]string, 0)
  headers := make([]string, 0)

  methods = append(methods, method) // append new method

  consumer, err := ConsumerNew(URL, methods, headers, false, 0)
  if err != nil {

    return errors.New("unable to create consumer")
  }

  c.consumers = append(c.consumers, consumer)

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

func (c *ManageConsumers) GetByOrigin(origin string) ConsumerImpl {

  for _, consumer := range c.consumers {

    if consumer.CheckOrigin(origin) {

      return consumer
    }
  }

  return nil
}

func (c *ManageConsumers) Check(method string, origin string) *http.Header {

  if consumer := c.Get(method, origin); consumer != nil {

    return consumer.Header()
  }

  return nil
}

func (c *ManageConsumers) GrantAll(origin string) error {

  // update consumer
  if consumer := c.GetByOrigin(origin); consumer != nil {

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

  return nil
}
