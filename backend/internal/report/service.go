package report

import (
	"fmt"
	"strconv"
)

type Repository interface {
	AidRows(from, to string) ([]AidRow, error)
	AidTotals(from, to string) ([]KindTotal, error)
	EventRows(from, to string) ([]EventRow, error)
	MotherRows(from, to, problemType string) ([]MotherRow, error)
	ProblemTypeTotals(from, to, problemType string) ([]ProblemTypeTotal, error)
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
		return s.mothers(from, to, f.ProblemType)
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

func (s *Service) mothers(from, to, problemType string) (*Report, error) {
	rows, err := s.repo.MotherRows(from, to, problemType)
	if err != nil {
		return nil, err
	}
	ptTotals, err := s.repo.ProblemTypeTotals(from, to, problemType)
	if err != nil {
		return nil, err
	}
	rep := &Report{
		Name: "mothers",
		Title: "Подопечные",
		Headers: []string{"№", "Мама", "Телефон", "Обратилась", "Детей", "Тип проблемы", "Ситуация", "Источник дохода", "Доход, ₽", "Анкета", "Статус", "Выдач помощи", "Заявок", "Консультаций", "Психолог", "Юрист", "Батюшка", "Центр занятости", "Мероприятий"},
		Rows: [][]string{},
		Totals: []Total{},
	}
	children, active, totalRequests, totalConsultations, totalEvents := 0, 0, 0, 0, 0
	totalPsychologist, totalLawyer, totalPriest, totalEmployment := 0, 0, 0, 0
	for _, r := range rows {
		children += r.Children
		totalRequests += r.RequestCount
		totalConsultations += r.ConsultationCount
		totalEvents += r.EventCount
		totalPsychologist += r.PsychologistCount
		totalLawyer += r.LawyerCount
		totalPriest += r.PriestCount
		totalEmployment += r.EmploymentCount
		if r.Status == "active" {
			active++
		}
		rep.Rows = append(rep.Rows, []string{
			strconv.Itoa(r.Id), r.Name, r.Phone, r.Registered, strconv.Itoa(r.Children),
			label(problemTypeNames, r.ProblemType), r.Situation,
			r.IncomeSource, strconv.Itoa(r.Income), label(profileStatusNames, r.ProfileStatus),
			label(motherStatusNames, r.Status), strconv.Itoa(r.HelpCount),
			strconv.Itoa(r.RequestCount), strconv.Itoa(r.ConsultationCount),
			strconv.Itoa(r.PsychologistCount), strconv.Itoa(r.LawyerCount),
			strconv.Itoa(r.PriestCount), strconv.Itoa(r.EmploymentCount),
			strconv.Itoa(r.EventCount),
		})
	}
	for _, pt := range ptTotals {
		rep.Totals = append(rep.Totals, Total{label(problemTypeNames, pt.ProblemType), strconv.Itoa(pt.Count)})
	}
	rep.Totals = append(rep.Totals,
		Total{"Всего мам", strconv.Itoa(len(rows))},
		Total{"На сопровождении", strconv.Itoa(active)},
		Total{"Всего детей", strconv.Itoa(children)},
		Total{"Всего заявок", strconv.Itoa(totalRequests)},
		Total{"Всего консультаций", strconv.Itoa(totalConsultations)},
		Total{"Психолог", strconv.Itoa(totalPsychologist)},
		Total{"Юрист", strconv.Itoa(totalLawyer)},
		Total{"Батюшка", strconv.Itoa(totalPriest)},
		Total{"Центр занятости", strconv.Itoa(totalEmployment)},
		Total{"Всего мероприятий", strconv.Itoa(totalEvents)})
	return rep, nil
}