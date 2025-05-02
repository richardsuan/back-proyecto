package domain

type Client struct {
	Name string `json:"name"`
	Data map[string]string `json:"data"` // Para almacenar los datos específicos del cliente
}