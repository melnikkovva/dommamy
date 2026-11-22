package adapter

type MotherService interface {
	Exists(id int) (bool, error)
}

type MotherDirectory struct {
	mothers MotherService
}

func NewMotherDirectory(mothers MotherService) *MotherDirectory {
	return &MotherDirectory{mothers: mothers}
}

func (m *MotherDirectory) Exists(id int) (bool, error) {
	return m.mothers.Exists(id)
}
