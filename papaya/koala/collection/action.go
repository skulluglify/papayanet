package collection

const (
  MoveLeft = iota
  MoveRight
  MoveUp
  MoveDown
)

const ()

type ActionMove struct {
  p int  // position movement, left or right
  r uint // range for next middle node
}

type ActionMoveImpl interface {
  Init(pos int, ran uint)
  Position() int
  Range() uint
}

func ActionMoveNew(pos int, ran uint) ActionMoveImpl {

  var actionMove ActionMoveImpl
  actionMove = &ActionMove{}
  actionMove.Init(pos, ran)
  return actionMove
}

func (a *ActionMove) Init(pos int, ran uint) {

  a.p = pos // position
  a.r = ran // range
}

func (a *ActionMove) Position() int {

  return a.p
}

func (a *ActionMove) Range() uint {

  return a.r
}
