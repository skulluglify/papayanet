package swag

import (
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/pp"
  "net/http"
  "strconv"
)

func SwagResponseSchemes(responses m.KMapImpl) m.KMapImpl {

  res := &m.KMap{}

  for _, enum := range responses.Enums() {

    k, v := enum.Tuple()

    if statusCode := m.KValueToString(k); statusCode != "" {

      // --- status code ---

      n, err := strconv.Atoi(statusCode)

      if err != nil {

        n = 200

        // wrong implement status code
        panic("wrong implemented status code in responses")
      }

      statusMessage := http.StatusText(n)

      // --- status code ---

      if mm := m.KMapCast(v); mm != nil {

        if body := m.KMapCast(mm.Get("body")); body != nil {

          for _, bEnum := range body.Enums() {

            bK, bV := bEnum.Tuple()

            if mimeTy := m.KValueToString(bK); bK != "" {

              if vM := m.KMapCast(bV); vM != nil {

                schema := vM.Get("schema")
                description := pp.QStr(m.KValueToString(vM.Get("description")), statusMessage)
                res.Put(statusCode, SwagContentSchema(mimeTy, schema, description))
              }
            }
          }
        }
      }
    }
  }

  return res
}
