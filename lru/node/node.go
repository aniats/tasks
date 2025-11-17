package node

import "testing"

type Node[V comparable] struct {
	Key  string
	Val  V
	Prev *Node[V]
	Next *Node[V]
}

func New[V comparable](key string, val V) *Node[V] {
	return &Node[V]{
		Key: key,
		Val: val,
	}
}


package vacancy

import (
"testing"

"gitlab.ozon.ru/internal-projects/mass-crm-api/internal/models/catalogue/citizenships"
)

func TestStatus_ShouldHandleImproveDataStatus(t *testing.T) {
	tests := []struct {
		name           string
		status         Status
		citizenshipID  citizenships.ID
		ticketVacancy  Vacancy
		expected       bool
	}{
		{
			name:          "should handle when status is ImproveData and security check is QuickSnils",
			status:        StatusImproveData,
			citizenshipID: citizenships.RussiaID,
			ticketVacancy: Vacancy{
				Security: Security{
					DefaultCheckType: SecurityCheckTypeShort,
					CheckTypeByCitizenship: map[citizenships.ID]SecurityCheckType{
						citizenships.RussiaID: SecurityCheckTypeQuickSnils,
					},
				},
			},
			expected: true,
		},
		{
			name:          "should not handle when status is not ImproveData",
			status:        Status("Заявка создана"),
			citizenshipID: citizenships.RussiaID,
			ticketVacancy: Vacancy{
				Security: Security{
					DefaultCheckType: SecurityCheckTypeShort,
					CheckTypeByCitizenship: map[citizenships.ID]SecurityCheckType{
						citizenships.RussiaID: SecurityCheckTypeQuickSnils,
					},
				},
			},
			expected: false,
		},
		{
			name:          "should not handle when security check is not QuickSnils",
			status:        StatusImproveData,
			citizenshipID: citizenships.RussiaID,
			ticketVacancy: Vacancy{
				Security: Security{
					DefaultCheckType: SecurityCheckTypeShort,
					CheckTypeByCitizenship: map[citizenships.ID]SecurityCheckType{
						citizenships.RussiaID: SecurityCheckTypeShort,
					},
				},
			},
			expected: false,
		},
		{
			name:          "should not handle when both conditions are false",
			status:        Status("Активен"),
			citizenshipID: citizenships.RussiaID,
			ticketVacancy: Vacancy{
				Security: Security{
					DefaultCheckType: SecurityCheckTypeShort,
					CheckTypeByCitizenship: map[citizenships.ID]SecurityCheckType{
						citizenships.RussiaID: SecurityCheckTypeFull,
					},
				},
			},
			expected: false,
		},
		{
			name:          "should handle with default check type as QuickSnils",
			status:        StatusImproveData,
			citizenshipID: citizenships.BelarusID,
			ticketVacancy: Vacancy{
				Security: Security{
					DefaultCheckType:       SecurityCheckTypeQuickSnils,
					CheckTypeByCitizenship: map[citizenships.ID]SecurityCheckType{},
				},
			},
			expected: true,
		},
		{
			name:          "should not handle when citizenship not in map and default is not QuickSnils",
			status:        StatusImproveData,
			citizenshipID: citizenships.KazakhstanID,
			ticketVacancy: Vacancy{
				Security: Security{
					DefaultCheckType:       SecurityCheckTypeShort,
					CheckTypeByCitizenship: map[citizenships.ID]SecurityCheckType{},
				},
			},
			expected: false,
		},
		{
			name:          "should handle with empty status string matching StatusImproveData",
			status:        StatusImproveData,
			citizenshipID: citizenships.RussiaID,
			ticketVacancy: Vacancy{
				Security: Security{
					CheckTypeByCitizenship: map[citizenships.ID]SecurityCheckType{
						citizenships.RussiaID: SecurityCheckTypeQuickSnils,
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.ShouldHandleImproveDataStatus(tt.citizenshipID, tt.ticketVacancy)

			if result != tt.expected {
				t.Errorf("ShouldHandleImproveDataStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}
ч