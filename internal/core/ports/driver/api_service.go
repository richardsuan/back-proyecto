package driver

type APIService interface {
    GetClientNames() ([]string, error)
    GetClientData(clientName string) ([]map[string]interface{}, error) // Agregado
}
