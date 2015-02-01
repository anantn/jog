package jog

type Value interface {
	Type(path ...string) Type
	Stringify(path ...string) (string, error)

	Get(path ...string) (Value, error)

	GetInt(path ...string) (int, error)
	GetUInt(path ...string) (uint, error)
	GetFloat(path ...string) (float64, error)

	GetBool(path ...string) (bool, error)
	GetString(path ...string) (string, error)

	GetArray(path ...string) ([]Value, error)
	GetObject(path ...string) (map[string]Value, error)
}

type Type int

const (
	TypeBool Type = iota
	TypeNull
	TypeArray
	TypeNumber
	TypeString
	TypeObject
	TypeUnknown
)

func New(val string) (Value, error) {
	return newRapidValue(val)
}
