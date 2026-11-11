package report

import (
	"fmt"
	"strconv"
)

type Repository interface {
	AidRows(from, to string) ([]AidRow, error)
	AidTotals(from, to string) ([]KindTotal, error)
	EventRows(from, to string) ([]EventRow, error)
	MotherRows(from, to string) ([]MotherRow, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Build(name string, f Filter) (*Report, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	from, to := f.Bounds()
	switch name {
	case "aid":
		return s.aid(from, to)
	case "events":
		return s.events(from, to)
	case "mothers":
		return s.mothers(from, to)
	}
	return nil, ErrorUnknown
}

func unit(kind string) string {
	if kind == "food" {
		return "₽"
	}
	return "шт."
}

func label(names map[string]string, key string) string {
	if v, ok := names[key]; ok {
		return v
	}
	return key
}

func (s *Service) aid(from, to string) (*Report, error) {
	rows, err := s.repo.AidRows(from, to)
	if err != nil {
		return nil, err
	}
	totals, err := s.repo.AidTotals(from, to)
	if err != nil {
		return nil, err
	}
	rep := &Report{
		Name: "aid",
		Title: "Оказанная помощь",
		Headers: []string{"Дата", "Мама", "Вид помощи", "Объём", "Единица", "Чек"},
		Rows: [][]string{},
		Totals: []Total{},
	}
	mothers := map[int]bool{}
	for _, r := range rows {
		mothers[r.MotherId] = true
		rep.Rows = append(rep.Rows, []string{r.Date, r.Mother, label(kindNames, r.Kind), strconv.Itoa(r.Sum), unit(r.Kind), r.CheckUrl})
	}
	for _, t := range totals {
		rep.Totals = append(rep.Totals, Total{label(kindNames, t.Kind), fmt.Sprintf("%d выдач, %d %s", t.Count, t.Sum, unit(t.Kind))})
	}
	rep.Totals = append(rep.Totals,
		Total{"Всего выдач", strconv.Itoa(len(rows))},
		Total{"Помощь получили, мам", strconv.Itoa(len(mothers))})
	return rep, nil
}

func (s *Service) events(from, to string) (*Report, error) {
	rows, err := s.repo.EventRows(from, to)
	if err != nil {
		return nil, err
	}
	rep := &Report{
		Name: "events",
		Title: "Участие в мероприятиях",
		Headers: []string{"Дата", "Мероприятие", "Тип", "Место", "Мама", "Статус"},
		Rows: [][]string{},
		Totals: []Total{},
	}
	mothers := map[int]bool{}
	counts := map[string]int{}
	for _, r := range rows {
		mothers[r.MotherId] = true
		counts[r.Status]++
		rep.Rows = append(rep.Rows, []string{r.When, r.Title, r.Type, r.Place, r.Mother, label(eventStatusNames, r.Status)})
	}
	rep.Totals = append(rep.Totals,
		Total{"Всего записей", strconv.Itoa(len(rows))},
		Total{"Пришли", strconv.Itoa(counts["came"])},
		Total{"Не пришли", strconv.Itoa(counts["absent"])},
		Total{"Ещё ожидаются", strconv.Itoa(counts["registered"])},
		Total{"Разных мам", strconv.Itoa(len(mothers))})
	return rep, nil
}

func (s *Service) mothers(from, to string) (*Report, error) {
	rows, err := s.repo.MotherRows(from, to)
	if err != nil {
		return nil, err
	}
	rep := &Report{
		Name: "mothers",
		Title: "Подопечные",
		Headers: []string{"№", "Мама", "Телефон", "Обратилась", "Детей", "Ситуация", "Источник дохода", "Доход, ₽", "Анкета", "Статус", "Выдач помощи"},
		Rows: [][]string{},
		Totals: []Total{},
	}
	children, active := 0, 0
	for _, r := range rows {
		children += r.Children
		if r.Status == "active" {
			active++
		}
		rep.Rows = append(rep.Rows, []string{
			strconv.Itoa(r.Id), r.Name, r.Phone, r.Registered, strconv.Itoa(r.Children), r.Situation,
			r.IncomeSource, strconv.Itoa(r.Income), label(profileStatusNames, r.ProfileStatus),
			label(motherStatusNames, r.Status), strconv.Itoa(r.HelpCount),
		})
	}
	rep.Totals = append(rep.Totals,
		Total{"Всего мам", strconv.Itoa(len(rows))},
		Total{"На сопровождении", strconv.Itoa(active)},
		Total{"Всего детей", strconv.Itoa(children)})
	return rep, nil
}