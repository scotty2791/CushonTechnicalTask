package domain

// FundName represents the available fund types
type FundName string

const (
	// CushonEquitiesFund represents the Cushon Equities Fund
	CushonEquitiesFund FundName = "Cushon Equities Fund"
)

// IsValid checks if the fund name is valid
func (f FundName) IsValid() bool {
	switch f {
	case CushonEquitiesFund:
		return true
	default:
		return false
	}
}

// String returns the string representation of the fund name
func (f FundName) String() string {
	return string(f)
} 