package report

import (
	"time"
)

var kindNames = map[string]string{"food": "Продукты", "clothes": "Вещи", "furniture": "Мебель", "stationery": "Канцелярия", "gifts": "Подарки", "consultation": "Консультация", "other": "Другое"}
var problemTypeNames = map[string]string{"": "Не указана", "housing": "Жилищная", "unemployment": "Безработица", "single_parent": "Одинокая мама", "low_income": "Низкий доход", "violence": "Насилие", "other": "Другое"}
var eventStatusNames = map[string]string{"registered": "Записана", "came": "Пришла", "absent": "Не пришла"}
var motherStatusNames = map[string]string{"review": "На рассмотрении", "active": "На сопровождении", "inactive": "Неактивна"}
var profileStatusNames = map[string]string{"": "Нет анкеты", "draft": "Черновик", "sent": "Отправлена", "checked": "Проверена"}

type Filter struct {
	From string
	To string
	ProblemType string
}

func (f Filter) Validate() error {
	var from, to time.Time
	var err error
	if f.From != "" {
		if from, err = time.Parse("2006-01-02", f.From); err != nil {
			return ErrorPeriod
		}
	}
	if f.To != "" {
		if to, err = time.Parse("2006-01-02", f.To); err != nil {
			return ErrorPeriod
		}
	}
	if f.From != "" && f.To != "" && from.After(to) {
		return ErrorPeriod
	}
	return nil
}

func (f Filter) Bounds() (string, string) {
	from, to := f.From, f.To
	if from == "" {
		from = "1970-01-02"
	}
	if to == "" {
		to = "2037-12-30"
	}
	return from, to
}

type AidRow struct {
	Date string
	MotherId int
	Mother string
	Kind string
	Sum int
	CheckUrl string
}

type KindTotal struct {
	Kind string
	Count int
	Sum int
}

type EventRow struct {
	When string
	Title string
	Type string
	Place string
	MotherId int
	Mother string
	Status string
}

type MotherRow struct {
	Id int
	Name string
	Phone string
	Registered string
	Children int
	Situation string
	ProblemType string
	IncomeSource string
	Income int
	ProfileStatus string
	Status string
	HelpCount int
	RequestCount int
	ConsultationCount int
	PsychologistCount int
	LawyerCount int
	PriestCount int
	EmploymentCount int
	EventCount int
}

type Total struct {
	Label string
	Value string
}

type ProblemTypeTotal struct {
	ProblemType string
	Count       int
}

type Report struct {
	Name string
	Title string
	Headers []string
	Rows [][]string
	Totals []Total
}