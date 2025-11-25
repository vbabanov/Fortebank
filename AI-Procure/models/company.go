package models

type CompanyProfile struct {
	Name    string
	Adress  string
	BIN_IIN string
}
type HistoryRecord struct {
	Company         CompanyProfile
	IsSuspicious    bool
	ContractAmount  float64
	ContractSubject string
}
