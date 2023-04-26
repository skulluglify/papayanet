package basic

import (
  "fmt"
  "net/http"
  "skfw/papaya"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
)

func MakeAuthTokenTask(pn papaya.NetImpl) *swag.SwagTask {

  return swag.MakeSwagTask("AuthToken", func(ctx *swag.SwagContext) error {

    pn.Logger().Log("Task Authorization Token ...")

    if m.KValueToBool(ctx.Event()) {

      headers := ctx.GetReqHeaders()

      if auth, ok := headers["Authorization"]; ok {

        fmt.Println(auth)

        return nil
      }
    }

    // failed authentication

    ctx.Prevent()
    return ctx.Status(http.StatusUnauthorized).JSON(kornet.MessageNew("Unauthorized", true))
  })
}
