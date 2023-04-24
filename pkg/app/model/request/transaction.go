package requestmodel

import validation "github.com/go-ozzo/ozzo-validation/v4"

// TransactionAll ...
type TransactionAll struct {
	Keyword   string `json:"keyword" query:"keyword"`
	Status    string `json:"status" query:"status"`
	PageToken string `json:"pageToken" query:"pageToken"`
}

// Validate ...
func (m TransactionAll) Validate() error {
	return validation.ValidateStruct(
		&m,
	)
}
