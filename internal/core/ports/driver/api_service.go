package driver

type APIService interface {
	GetClientNames() ([]string, error)
	GetClientRawData(clientName string) ([][]string, error)
}
