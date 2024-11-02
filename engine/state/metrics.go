package state

type ErrorMetrics struct {
	emptyResponses   uint
	unparsedResonses uint
	unknownActions   uint
	invalidActions   uint
	erroredActions   uint
	timedoutActions  uint
}

type Metrics struct {
	maxStep        uint
	currentStep    uint
	validResponses uint
	validActions   uint
	successActions uint
	errors         uint
}

func NewMetrics(maxStep uint) *Metrics {
	return &Metrics{
		maxStep: maxStep,
	}
}
