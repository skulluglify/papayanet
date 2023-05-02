package bpack

import (
  "errors"
  "github.com/gofiber/fiber/v2"
  "net/http"
  "strconv"
)

var PacketIsNULL = errors.New("packet is NULL")
var ContentTypeIsEmpty = errors.New("content-type is empty")
var PathIsEmpty = errors.New("path is empty")

func HttpExposePacket(app *fiber.App, path string, cTy string, packet *Packet) error {

  if path != "" {

    if cTy != "" {

      if packet != nil {

        app.Get(path, func(ctx *fiber.Ctx) error {

          res := ctx.Response()

          if packet.Charset != "" {

            cTy += "; charset:" + packet.Charset
          }

          res.Header.Set("Content-Type", cTy)
          res.Header.Set("Content-Length", strconv.FormatUint(packet.Size, 10))

          return ctx.Status(http.StatusOK).Send(packet.Data)
        })

        return nil
      }

      return PacketIsNULL
    }

    return ContentTypeIsEmpty
  }

  return PathIsEmpty
}
