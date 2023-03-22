// models/models.go
package models

type Contributor struct {
	ID        int64
	Login     string
	Repos     []*Repo
}

type Repo struct {
	ID          int64
	Name        string
	FullName    string
	Contributor *Contributor
}