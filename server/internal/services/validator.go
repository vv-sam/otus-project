package services

// Вроде же Go-way названия интерфейсов (типа validate + er)
type validater interface {
	Validate() error
}

type Validator struct{}

func (vv *Validator) IsValid(vs ...validater) bool {
	for _, v := range vs {
		if err := v.Validate(); err != nil {
			return false
		}
	}

	return true
}
