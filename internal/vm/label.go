package vm

// Label represents a named position in code
type Label struct {
	Name string
	Pos  int32
}

// NewLabel creates a new label
func NewLabel(name string, position int32) Label {
	return Label{
		Name: name,
		Pos:  position,
	}
}
