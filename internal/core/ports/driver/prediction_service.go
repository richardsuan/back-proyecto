package driver

type PredictionService interface {
	PredictAndSave(clientName string, data map[string]float64) (bool, error)
}
