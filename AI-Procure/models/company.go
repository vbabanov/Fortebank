package models

type CompanyProfile struct {
	Name   string
	Adress string
	BIN    string
}
type HistoryRecord struct {
	Company      CompanyProfile
	IsSuspicious bool
}
