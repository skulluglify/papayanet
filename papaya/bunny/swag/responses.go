package swag

import (
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/pp"
)

func SwagResponseSchemes(responses m.KMapImpl) m.KMapImpl {

  res := &m.KMap{}

  for _, enum := range responses.Enums() {

    k, v := enum.Tuple()

    if statusCode := m.KValueToString(k); statusCode != "" {

      if mm := m.KMapCast(v); mm != nil {

        if body := m.KMapCast(mm.Get("body")); body != nil {

          for _, bEnum := range body.Enums() {

            bK, bV := bEnum.Tuple()

            if mimeTy := m.KValueToString(bK); bK != "" {

              if vM := m.KMapCast(bV); vM != nil {

                schema := vM.Get("schema")
                description := pp.QStr(m.KValueToString(vM.Get("description")), "Ok")
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
