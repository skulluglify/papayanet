package bpack

import (
  "io/fs"
  "os"
  "path/filepath"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/tools/posix"
  "strings"
)

func GetCwd() (string, error) {

  var err error
  var dir string

  if dir, err = os.Getwd(); err != nil {

    var noop string
    return noop, err
  }

  return dir, nil
}

func GetCurrentPath() (string, error) {

  var err error
  var ex string

  if ex, err = os.Executable(); err != nil {

    var noop string
    return noop, err
  }

  return filepath.Dir(ex), nil
}

func FindDataPath(paths string) string {

  var err error
  var cwd string

  var path posix.KPathImpl
  var file kio.KFileImpl

  if cwd, err = GetCwd(); err != nil {

    var noop string
    return noop
  }

  for _, p := range strings.Split(paths, ":") {

    if !strings.HasPrefix(p, "/") {

      path = posix.KPathNew(p)

    } else {

      path = posix.KPathNew(cwd)
      path.JoinStr(p)
    }

    p = path.String()

    file = kio.KFileNew(p)

    if file.IsDir() {

      return p
    }
  }

  var noop string
  return noop
}

func ReadAllDataFromPath(path string) map[string][]byte {

  var found bool
  var buff []byte

  var data map[string][]byte

  data = make(map[string][]byte)

  filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {

    if err != nil {

      return err
    }

    if !info.IsDir() {

      buff, err = os.ReadFile(p)
      if err != nil {

        return err
      }

      if p, found = strings.CutPrefix(p, path); found {

        p = posix.KPathNew("/data").JoinStr(p)
      }

      data[p] = buff
    }

    return nil
  })

  return data
}
