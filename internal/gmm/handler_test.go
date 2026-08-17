package gmm

import (
	"testing"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi/models"
)

// TestAssignLadnInfoNonStringDnn is a regression test for free5gc/free5gc#1071.
// A subscriber whose DNN was provisioned as a JSON number is unmarshalled into
// models.DnnInfo.Dnn (typed interface{}) as a float64. The former bare
// dnnInfo.Dnn.(string) assertions in assignLadnInfo panicked on such a value,
// which the NGAP worker's recover turned into a permanent worker death (DoS).
// assignLadnInfo must now skip non-string DNN entries instead of panicking.
func TestAssignLadnInfoNonStringDnn(t *testing.T) {
	tests := []struct {
		name string
		dnn  interface{}
	}{
		{name: "numeric DNN (float64) does not panic", dnn: float64(123456)},
		{name: "valid string DNN is handled normally", dnn: "internet"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ue := &context.AmfUe{}
			ue.GmmLog = logger.GmmLog
			// LADNIndication left nil so control reaches the else-if
			// SmfSelectionData branch containing the reported sink.
			ue.RegistrationRequest = &nasMessage.RegistrationRequest{}
			ue.SmfSelectionData = &models.SmfSelectionSubscriptionData{
				SubscribedSnssaiInfos: map[string]models.SnssaiInfo{
					"01010203": {
						DnnInfos: []models.DnnInfo{{Dnn: tc.dnn}},
					},
				},
			}

			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("assignLadnInfo panicked on DNN %v (%T): %v", tc.dnn, tc.dnn, r)
				}
			}()

			assignLadnInfo(ue, models.AccessType__3_GPP_ACCESS)
		})
	}
}
