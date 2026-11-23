package donor

import (
	"time"
)

type Repository interface {
	SaveDonor(d *Donor) error
	UpdateDonor(d *Donor) error
	DonorById(id int) (*Donor, error) 
	ListDonors(q string) ([]*Summary, error)
	SaveDonation(n *Donation) error
	DonationsByDonor(donorId int) ([]*Donation, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name, kind, phone, email, comment string) (*Donor, error) {
	d := &Donor{Name: name, Kind: kind, Phone: phone, Email: email, Comment: comment, CreatedAt: time.Now()}
	if err := d.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.SaveDonor(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Update(id int, name, kind, phone, email, comment string) (*Donor, error) {
	d, err := s.repo.DonorById(id)
	if err != nil {
		return nil, err
	}
	d.Name, d.Kind, d.Phone, d.Email, d.Comment = name, kind, phone, email, comment
	if err := d.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateDonor(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Get(id int) (*Donor, []*Donation, error) {
	d, err := s.repo.DonorById(id)
	if err != nil {
		return nil, nil, err
	}
	donations, err := s.repo.DonationsByDonor(id)
	if err != nil {
		return nil, nil, err
	}
	return d, donations, nil
}

func (s *Service) List(q string) ([]*Summary, error) {
	return s.repo.ListDonors(q)
}

func (s *Service) AddDonation(donorId, amount int, date, comment string) (*Donation, error) {
	n, err := NewDonation(donorId, amount, date, comment, time.Now())
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.DonorById(donorId); err != nil {
		return nil, err 
	}
	if err := s.repo.SaveDonation(n); err != nil {
		return nil, err
	}
	return n, nil
}