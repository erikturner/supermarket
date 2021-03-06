package bsapi

type Produce struct {
	ProduceCode 		*string 		`json:"produceCode,omitempty"`
	Name 				*string 		`json:"name,omitempty"`
	UnitPrice			*float64 		`json:"unitPrice"`
}
