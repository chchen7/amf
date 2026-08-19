package gmm

import (
	"testing"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	nasMessage "github.com/free5gc/nas/message"
	"github.com/free5gc/openapi/models"
)

// TestAssignLadnInfoDnn verifies LADN assignment with the migrated OpenAPI DNN type.
func TestAssignLadnInfoDnn(t *testing.T) {
	tests := []struct {
		name string
		dnn  string
	}{
		{name: "valid string DNN is handled normally", dnn: "internet"},
		{name: "empty DNN is handled normally", dnn: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ue := &context.AmfUe{}
			ue.GmmLog = logger.GmmLog
			// LADNIndication left nil so control reaches the else-if
			// SmfSelectionData branch containing the reported sink.
			ue.RegistrationRequest = &nasMessage.RegReq{}
			ue.SmfSelectionData = &models.Udm_SDM_SmfSelectionSubscriptionData{
				SubscribedSnssaiInfos: map[string]models.Udm_SDM_SnssaiInfo{
					"01010203": {
						DnnInfos: []models.Udm_SDM_DnnInfo{{Dnn: tc.dnn}},
					},
				},
			}

			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("assignLadnInfo panicked on DNN %v (%T): %v", tc.dnn, tc.dnn, r)
				}
			}()

			assignLadnInfo(ue, models.AccessType_3_GPP_ACCESS)
		})
	}
}
