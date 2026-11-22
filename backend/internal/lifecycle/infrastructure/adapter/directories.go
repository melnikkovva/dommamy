package adapter

import "context"

type MotherChecker interface {
	Exists(ctx context.Context, id int) (bool, error)
}

type MotherDirectory struct {
	checker MotherChecker
}

func NewMotherDirectory(checker MotherChecker) *MotherDirectory {
	return &MotherDirectory{checker: checker}
}

func (m *MotherDirectory) Exists(ctx context.Context, id int) (bool, error) {
	return m.checker.Exists(ctx, id)
}
