package adapter

type MotherChecker interface {
	Exists(id int) (bool, error)
}

type MotherDirectory struct {
	mothers MotherChecker
}

func NewMotherDirectory(mothers MotherChecker) *MotherDirectory {
	return &MotherDirectory{mothers: mothers}
}

func (m *MotherDirectory) Exists(id int) (bool, error) {
	return m.mothers.Exists(id)
}
