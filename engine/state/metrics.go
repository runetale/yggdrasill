package state

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorMetrics struct {
	emptyResponses    uint
	unparsedResponses uint
	unknownActions    uint
	invalidActions    uint
	erroredActions    uint
	timedoutActions   uint
}

func (e ErrorMetrics) HasResponseErrors() bool {
	return e.emptyResponses > 0 || e.unparsedResponses > 0
}

func (e ErrorMetrics) HasActionErrors() bool {
	return e.erroredActions > 0 || e.unknownActions > 0 || e.invalidActions > 0
}

func MemoryStats() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc
}

func HumanBytes(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
}

type Metrics struct {
	maxStep        uint
	currentStep    uint
	validResponses uint
	validActions   uint
	successActions uint
	errors         ErrorMetrics
}

func NewMetrics(maxStep uint) *Metrics {
	return &Metrics{
		maxStep: maxStep,
	}
}

func (m *Metrics) Display() string {
	var sb strings.Builder

	sb.WriteString("step:")
	if m.maxStep > 0 {
		sb.WriteString(fmt.Sprintf("%d/%d ", m.currentStep, m.maxStep))
	} else {
		sb.WriteString(fmt.Sprintf("%d ", m.currentStep))
	}

	if m.errors.HasResponseErrors() {
		sb.WriteString(fmt.Sprintf(
			"responses(valid:%d empty:%d broken:%d) ",
			m.validResponses, m.errors.emptyResponses, m.errors.unparsedResponses,
		))
	} else if m.validResponses > 0 {
		sb.WriteString(fmt.Sprintf("responses:%d ", m.validResponses))
	}

	if m.errors.HasActionErrors() {
		sb.WriteString(fmt.Sprintf(
			"actions(valid:%d ok:%d errored:%d unknown:%d invalid:%d) ",
			m.validActions,
			m.successActions,
			m.errors.erroredActions,
			m.errors.unknownActions,
			m.errors.invalidActions,
		))
	} else if m.validActions > 0 {
		sb.WriteString(fmt.Sprintf("actions:%d ", m.validActions))
	}

	memUsage := MemoryStats()
	sb.WriteString(fmt.Sprintf("mem:%s", HumanBytes(memUsage)))

	return sb.String()
}
