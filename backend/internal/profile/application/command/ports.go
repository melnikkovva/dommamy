package command

type MotherDirectory interface {
	Exists(id int) (bool, error)
}
