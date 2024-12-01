package types

// ----- User Define Mod -----

type RoleMod struct {
	Type        int
	Name        string
	Description string
	Actions     []*ActionMod
	Functions   []*Function
}

type EntityMod struct {
	Type        int
	Name        string
	Description string
	Functions   []*Function
}

type ActionMod struct {
	Type        int
	Name        string
	Description string
	Functions   []*Function
}

type EventMod struct {
	Type        int
	Name        string
	Description string
	From        *Object
	To          *Object
	Functions   []*Function
}

const (
	RoleObject = iota
	SystemObject
	EntityObject
	ActionObject
)

type Object struct {
	Type int
	Name string
}

type Function struct {
	CodeType string
	Name     string
	Code     []byte
}

func (f *Function) Do(params any) (any, error) {
	return nil, nil
}

// ---- runtime instance -----

type Role struct {
	ID         int
	Mod        *RoleMod
	Properties []*Entity
	Position   *Position
}

type Entity struct {
	ID       int
	Mod      *EntityMod
	Position *Position
	Children []*Entity
}

type Position struct {
	X int
	Y int
	Z int
}
