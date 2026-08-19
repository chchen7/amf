package callback

import (
	"net/url"
	"reflect"
	"strings"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	Namf_Communication "github.com/free5gc/openapi/amf/Comm"
	"github.com/free5gc/openapi/models"
)

func callbackServiceNfType(svcName string) (models.Nrf_NFMgmt_NFType, bool) {
	switch {
	case strings.HasPrefix(svcName, "npcf"):
		return models.Nrf_NFMgmt_NFType_PCF, true
	case strings.HasPrefix(svcName, "nsmf"):
		return models.Nrf_NFMgmt_NFType_SMF, true
	case strings.HasPrefix(svcName, "nudm"):
		return models.Nrf_NFMgmt_NFType_UDM, true
	case strings.HasPrefix(svcName, "nausf"):
		return models.Nrf_NFMgmt_NFType_AUSF, true
	case strings.HasPrefix(svcName, "namf"):
		return models.Nrf_NFMgmt_NFType_AMF, true
	default:
		return "", false
	}
}

func SendAmfStatusChangeNotify(amfStatus string, guamiList []models.Guami) {
	amfSelf := amf_context.GetSelf()

	amfSelf.AMFStatusSubscriptions.Range(func(key, value interface{}) bool {
		subscriptionData := value.(models.Amf_Comm_SubscriptionData)

		configuration := Namf_Communication.NewConfiguration()
		client := Namf_Communication.NewAPIClient(configuration)
		amfStatusNotification := models.Amf_Comm_AmfStatusChangeNotification{}
		amfStatusInfo := models.Amf_Comm_AmfStatusInfo{}

		for _, guami := range guamiList {
			for _, subGumi := range subscriptionData.GuamiList {
				if reflect.DeepEqual(guami, subGumi) {
					// AMF status is available
					amfStatusInfo.GuamiList = append(amfStatusInfo.GuamiList, guami)
				}
			}
		}

		amfStatusInfo = models.Amf_Comm_AmfStatusInfo{
			StatusChange:     (models.Amf_Comm_StatusChange)(amfStatus),
			TargetAmfRemoval: "",
			TargetAmfFailure: "",
		}

		amfStatusNotification.AmfStatusInfoList = append(amfStatusNotification.AmfStatusInfoList, amfStatusInfo)
		uri := subscriptionData.AmfStatusUri

		amfStatusNotificationReq := Namf_Communication.AmfStatusChangeNotifyRequest{
			RequestBody: &amfStatusNotification,
		}

		var callbackSvcName models.Nrf_NFMgmt_ServiceName
		var targetNFType models.Nrf_NFMgmt_NFType
		if parsedURI, err := url.Parse(uri); err == nil {
			seg := strings.SplitN(strings.TrimPrefix(parsedURI.Path, "/"), "/", 2)[0]
			if nfType, ok := callbackServiceNfType(seg); ok {
				callbackSvcName = models.Nrf_NFMgmt_ServiceName(seg)
				targetNFType = nfType
			}
		}

		ctx, pd, err := amfSelf.GetTokenCtx(callbackSvcName, targetNFType)
		if err != nil {
			HttpLog.Warnf("SendAmfStatusChangeNotify get token failed: %+v", pd)
			return false
		}

		logger.ProducerLog.Infof("[AMF] Send Amf Status Change Notify to %s", uri)
		_, err = client.SubscriptionsCollectionCollectionApi.
			AmfStatusChangeNotify(ctx, uri, &amfStatusNotificationReq)
		if err != nil {
			HttpLog.Errorln(err.Error())
		}
		return true
	})
}
