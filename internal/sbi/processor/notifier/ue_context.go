package callback

import (
	"fmt"

	amf_context "github.com/free5gc/amf/internal/context"
	Namf_Communication "github.com/free5gc/openapi/amf/Comm"
	"github.com/free5gc/openapi/models"
)

func SendN2InfoNotifyN2Handover(ue *amf_context.AmfUe, releaseList []int32) error {
	if ue.HandoverNotifyUri == "" {
		return fmt.Errorf("N2 Info Notify N2Handover failed(uri dose not exist)")
	}
	configuration := Namf_Communication.NewConfiguration()
	client := Namf_Communication.NewAPIClient(configuration)

	n2InformationNotification := models.Amf_Comm_N2InformationNotification{
		N2NotifySubscriptionId: ue.Supi,
		ToReleaseSessionList:   releaseList,
		NotifyReason:           models.Amf_Comm_N2InfoNotifyReason_HANDOVER_COMPLETED,
	}

	n2InformationNotificationReq := Namf_Communication.N2InfoNotifyHandoverCompleteRequest{
		RequestBody: &n2InformationNotification,
	}

	ctx, pd, err := amf_context.GetSelf().GetTokenCtx(
		models.Nrf_NFMgmt_ServiceName("namf-callback"), models.Nrf_NFMgmt_NFType_AMF)
	if err != nil {
		HttpLog.Warnf("SendN2InfoNotifyN2Handover get token failed: %+v", pd)
		return err
	}

	_, err = client.IndividualUeContextDocumentApi.
		N2InfoNotifyHandoverComplete(ctx, ue.HandoverNotifyUri, &n2InformationNotificationReq)

	if err == nil {
		// TODO: handle Msg
	} else {
		HttpLog.Errorln(err.Error())
		return err
	}
	return nil
}
