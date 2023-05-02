package bpack

import (
  "fmt"
  "net/http"
  "os"
  "skfw/papaya/koala"
  "skfw/papaya/koala/kornet"
  "skfw/papaya/koala/tools/posix"
  "strconv"
)

func Main() {

  cwd, err := GetCwd()

  if err != nil {

    fmt.Println(err)
    return
  }

  dataPath := FindDataPath(PATH)
  dataMap := ReadAllDataFromPath(posix.KPathNew(dataPath).JoinStr("data"))

  fmt.Println("detect", dataPath)

  var data string

  data = "package bpack\nvar Pkts Packets = Packets{\n"

  var limit int

  var temp string

  limit = 20

  var nameLenAvg int

  nameLenAvg = 0

  for name := range dataMap {

    if z := len(name); z > nameLenAvg {

      nameLenAvg = z
    }

  }

  nameLenAvg += 2

  var mimetype string
  var charset string

  var size uint64

  for name, buff := range dataMap {

    temp = ""

    mimetype, charset = kornet.KSafeContentTy(http.DetectContentType(buff))

    // begin
    data += "{\"" + name + "\",\"" + mimetype + "\",\"" + charset + "\",[]byte{\n"

    size = 0

    for i, b := range buff {

      if i > 0 && i%limit == 0 {

        fmt.Print("\r", koala.KStrPadEnd(name, nameLenAvg), koala.KStrPadStart(kornet.ReprByte(uint64(i+1)), 8))
        data += temp + "\n"
        temp = ""
      }

      temp += strconv.Itoa(int(b)) + ", "
      size++
    }

    // rest
    if temp != "" {

      data += temp + "\n"
      temp = ""
    }

    data += "\n"
    data += "}," + strconv.FormatUint(size, 10) + "},\n" // end

    fmt.Println() // new line
  }

  data += "}"

  fmt.Println("writing data.go ...")

  dataOut := posix.KPathNew(cwd).JoinStr("data.go")

  file, err := os.Create(dataOut)
  if err != nil {

    fmt.Println(err)
    return
  }
  defer file.Close()

  _, err = file.WriteString(data)
  if err != nil {

    fmt.Println(err)
    return
  }

  file.Close()
}
