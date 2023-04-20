package posix

import (
  "PapayaNet/papaya/koala"
)

type KPath struct {
  path string
}

type KPathImpl interface {
  Init(name string)
  String() string
  DirStr() string
  Dir() KPathImpl
  BaseStr() string
  Base() KPathImpl
  PopStr() string
  Pop() KPathImpl
  JoinStr(names ...string) string
  Join(paths ...KPathImpl) KPathImpl
  Copy() KPathImpl
}

func KPathNew(name string) KPathImpl {

  p := &KPath{}
  p.Init(name)

  return p
}

func (p *KPath) Init(name string) {

  // normalize

  // use `path.Split`

  // --- end

  p.path = name
}

func (p *KPath) String() string {

  return p.path
}

func (p *KPath) DirStr() string {

  // prefix
  var i, j, n int

  data := []byte(p.path)

  n = len(data)

  for i = 0; i < n; i++ {

    j = n - i - 1

    // \/ is 47
    if data[j] != 47 {

      continue
    }

    break
  }

  return string(data[:j])
}

func (p *KPath) Dir() KPathImpl {

  return KPathNew(p.DirStr())
}

func (p *KPath) BaseStr() string {

  // suffix
  var i, j, n int
  var suffix string

  data := []byte(p.path)

  n = len(data)

  for i = 0; i < n; i++ {

    j = n - i - 1

    // waste 1 time for iteration
    // but that look good for me

    // \/ is 47
    if data[j] != 47 {

      suffix = string(data[j]) + suffix
      continue
    }

    break
  }

  return suffix
}

func (p *KPath) Base() KPathImpl {

  return KPathNew(p.BaseStr())
}

func (p *KPath) PopStr() string {

  // prefix and suffix
  var i, j, n int
  var suffix string

  data := []byte(p.path)

  n = len(data)

  for i = 0; i < n; i++ {

    j = n - i - 1

    // \/ is 47
    if data[j] != 47 {

      suffix = string(data[j]) + suffix
      continue
    }

    break
  }

  // keep prefix in current session
  p.path = string(data[:j])

  // return
  return suffix
}

func (p *KPath) Pop() KPathImpl {

  return KPathNew(p.PopStr())
}

func (p *KPath) JoinStr(names ...string) string {

  var i, n int
  var name string

  n = len(names)

  for i = 0; i < n; i++ {

    name = names[i]

    // prefix `name`
    if koala.KStrHasPrefixChar(name, "/") {

      name = name[1:]
    }

    // suffix `name`
    if koala.KStrHasSuffixChar(name, "/") {

      name = name[:len(name)-1]
    }

    // suffix p.path
    if koala.KStrHasSuffixChar(p.path, "/") {

      p.path = p.path[:len(p.path)-1]
    }

    p.path += "/" + name
  }

  return p.path
}

func (p *KPath) Join(paths ...KPathImpl) KPathImpl {

  var i, n int
  var name string

  n = len(paths)

  for i = 0; i < n; i++ {

    name = paths[i].String()

    // prefix `name`
    if koala.KStrHasPrefixChar(name, "/") {

      name = name[1:]
    }

    // suffix `name`
    if koala.KStrHasSuffixChar(name, "/") {

      name = name[:len(name)-1]
    }

    // suffix p.path
    if koala.KStrHasSuffixChar(p.path, "/") {

      p.path = p.path[:len(p.path)-1]
    }

    p.path += "/" + name
  }

  return KPathNew(p.path)
}

func (p *KPath) Copy() KPathImpl {

  return KPathNew(p.path)
}
