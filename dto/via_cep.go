package dto

type Location struct {
	CEP      string `json:"cep"`
	Location string `json:"localidade"`
}
