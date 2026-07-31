package message

import (
	"encoding/hex"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	ngapConvert "github.com/free5gc/amf/internal/ngap/convert"
	"github.com/free5gc/amf/internal/util"
	ngapAper "github.com/free5gc/ngap/aper"
	ngapIE "github.com/free5gc/ngap/ie"
	"github.com/free5gc/openapi/models"
)

func AppendPDUSessionResourceSetupListSUReq(list *ngapIE.PDUSessionResourceSetupListSUReq,
	pduSessionId int32, snssai models.Snssai, nasPDU []byte, transfer []byte,
) {
	ngapSnssai := ngapConvert.SNssaiToNgap(snssai)
	transferValue := ngapAper.OctetString(transfer)
	item := ngapIE.PDUSessionResourceSetupItemSUReq{
		PDUSessionID:                           &ngapIE.PDUSessionID{Value: int64(pduSessionId)},
		SNSSAI:                                 &ngapSnssai,
		PDUSessionResourceSetupRequestTransfer: &transferValue,
	}
	if nasPDU != nil {
		item.PDUSessionNASPDU = &ngapIE.NASPDU{Value: nasPDU}
	}
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceSetupListHOReq(list *ngapIE.PDUSessionResourceSetupListHOReq,
	pduSessionId int32, snssai models.Snssai, transfer []byte,
) {
	ngapSnssai := ngapConvert.SNssaiToNgap(snssai)
	transferValue := ngapAper.OctetString(transfer)
	item := ngapIE.PDUSessionResourceSetupItemHOReq{
		PDUSessionID:            &ngapIE.PDUSessionID{Value: int64(pduSessionId)},
		SNSSAI:                  &ngapSnssai,
		HandoverRequestTransfer: &transferValue,
	}
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceSetupListCxtReq(list *ngapIE.PDUSessionResourceSetupListCxtReq,
	pduSessionId int32, snssai models.Snssai, nasPDU []byte, transfer []byte,
) {
	ngapSnssai := ngapConvert.SNssaiToNgap(snssai)
	transferValue := ngapAper.OctetString(transfer)
	item := ngapIE.PDUSessionResourceSetupItemCxtReq{
		PDUSessionID:                           &ngapIE.PDUSessionID{Value: int64(pduSessionId)},
		SNSSAI:                                 &ngapSnssai,
		PDUSessionResourceSetupRequestTransfer: &transferValue,
	}
	if nasPDU != nil {
		item.NASPDU = &ngapIE.NASPDU{Value: nasPDU}
	}
	list.List = append(list.List, item)
}

func ConvertPDUSessionResourceSetupListCxtReqToSUReq(
	listCxtReq *ngapIE.PDUSessionResourceSetupListCxtReq,
) *ngapIE.PDUSessionResourceSetupListSUReq {
	if listCxtReq == nil {
		return nil
	}
	listSUReq := ngapIE.PDUSessionResourceSetupListSUReq{}
	for _, itemCxt := range listCxtReq.List {
		itemSU := ngapIE.PDUSessionResourceSetupItemSUReq{
			PDUSessionID:                           itemCxt.PDUSessionID,
			PDUSessionNASPDU:                       itemCxt.NASPDU,
			SNSSAI:                                 itemCxt.SNSSAI,
			PDUSessionResourceSetupRequestTransfer: itemCxt.PDUSessionResourceSetupRequestTransfer,
		}
		listSUReq.List = append(listSUReq.List, itemSU)
	}
	return &listSUReq
}

func AppendPDUSessionResourceModifyListModReq(list *ngapIE.PDUSessionResourceModifyListModReq,
	pduSessionId int32, nasPDU []byte, transfer []byte,
) {
	transferValue := ngapAper.OctetString(transfer)
	item := ngapIE.PDUSessionResourceModifyItemModReq{
		PDUSessionID:                            &ngapIE.PDUSessionID{Value: int64(pduSessionId)},
		PDUSessionResourceModifyRequestTransfer: &transferValue,
	}
	if nasPDU != nil {
		item.NASPDU = &ngapIE.NASPDU{Value: nasPDU}
	}
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceModifyListModCfm(list *ngapIE.PDUSessionResourceModifyListModCfm,
	pduSessionId int64, transfer []byte,
) {
	transferValue := ngapAper.OctetString(transfer)
	item := ngapIE.PDUSessionResourceModifyItemModCfm{
		PDUSessionID:                            &ngapIE.PDUSessionID{Value: pduSessionId},
		PDUSessionResourceModifyConfirmTransfer: &transferValue,
	}
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceFailedToModifyListModCfm(list *ngapIE.PDUSessionResourceFailedToModifyListModCfm,
	pduSessionId int64, transfer []byte,
) {
	transferValue := ngapAper.OctetString(transfer)
	item := ngapIE.PDUSessionResourceFailedToModifyItemModCfm{
		PDUSessionID: &ngapIE.PDUSessionID{Value: pduSessionId},
		PDUSessionResourceModifyIndicationUnsuccessfulTransfer: &transferValue,
	}
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceToReleaseListRelCmd(list *ngapIE.PDUSessionResourceToReleaseListRelCmd,
	pduSessionId int32, transfer []byte,
) {
	if list == nil {
		return
	}
	transferValue := ngapAper.OctetString(transfer)
	item := ngapIE.PDUSessionResourceToReleaseItemRelCmd{
		PDUSessionID:                             &ngapIE.PDUSessionID{Value: int64(pduSessionId)},
		PDUSessionResourceReleaseCommandTransfer: &transferValue,
	}
	list.List = append(list.List, item)
}

func BuildIEMobilityRestrictionList(ue *context.AmfUe) ngapIE.MobilityRestrictionList {
	servingPLMN := ngapConvert.PlmnIdToNgap(ue.PlmnId)
	mobilityRestrictionList := ngapIE.MobilityRestrictionList{ServingPLMN: &servingPLMN}

	if ue.AccessAndMobilitySubscriptionData != nil && len(ue.AccessAndMobilitySubscriptionData.RatRestrictions) > 0 {
		mobilityRestrictionList.RATRestrictions = new(ngapIE.RATRestrictions)
		ratRestrictions := mobilityRestrictionList.RATRestrictions
		for _, ratType := range ue.AccessAndMobilitySubscriptionData.RatRestrictions {
			plmn := ngapConvert.PlmnIdToNgap(ue.PlmnId)
			restriction := ngapConvert.RATRestrictionInformationToNgap(ratType)
			item := ngapIE.RATRestrictionsItem{
				PLMNIdentity:              &plmn,
				RATRestrictionInformation: &restriction,
			}
			ratRestrictions.List = append(ratRestrictions.List, item)
		}
	}

	if ue.AccessAndMobilitySubscriptionData != nil && len(ue.AccessAndMobilitySubscriptionData.ForbiddenAreas) > 0 {
		mobilityRestrictionList.ForbiddenAreaInformation = new(ngapIE.ForbiddenAreaInformation)
		forbiddenAreaInformation := mobilityRestrictionList.ForbiddenAreaInformation
		for _, info := range ue.AccessAndMobilitySubscriptionData.ForbiddenAreas {
			plmn := ngapConvert.PlmnIdToNgap(ue.PlmnId)
			item := ngapIE.ForbiddenAreaInformationItem{
				PLMNIdentity:  &plmn,
				ForbiddenTACs: &ngapIE.ForbiddenTACs{},
			}
			for _, tac := range info.Tacs {
				tacBytes, err := hex.DecodeString(tac)
				if err != nil {
					logger.NgapLog.Errorf(
						"[Error] DecodeString tac error: %+v", err)
				}
				tacNgap := ngapIE.TAC{Value: tacBytes}
				item.ForbiddenTACs.List = append(item.ForbiddenTACs.List, tacNgap)
			}
			forbiddenAreaInformation.List = append(forbiddenAreaInformation.List, item)
		}
	}

	if ue.AmPolicyAssociation != nil && ue.AmPolicyAssociation.ServAreaRes != nil {
		mobilityRestrictionList.ServiceAreaInformation = new(ngapIE.ServiceAreaInformation)
		serviceAreaInformation := mobilityRestrictionList.ServiceAreaInformation

		plmn := ngapConvert.PlmnIdToNgap(ue.PlmnId)
		item := ngapIE.ServiceAreaInformationItem{PLMNIdentity: &plmn}
		var tacList []ngapIE.TAC
		for _, area := range ue.AmPolicyAssociation.ServAreaRes.Areas {
			for _, tac := range area.Tacs {
				tacBytes, err := hex.DecodeString(tac)
				if err != nil {
					logger.NgapLog.Errorf(
						"[Error] DecodeString tac error: %+v", err)
				}
				tacNgap := ngapIE.TAC{Value: tacBytes}
				tacList = append(tacList, tacNgap)
			}
		}
		if ue.AmPolicyAssociation.ServAreaRes.RestrictionType == models.RestrictionType_ALLOWED_AREAS {
			item.AllowedTACs = new(ngapIE.AllowedTACs)
			item.AllowedTACs.List = append(item.AllowedTACs.List, tacList...)
		} else {
			item.NotAllowedTACs = new(ngapIE.NotAllowedTACs)
			item.NotAllowedTACs.List = append(item.NotAllowedTACs.List, tacList...)
		}
		serviceAreaInformation.List = append(serviceAreaInformation.List, item)
	}
	return mobilityRestrictionList
}

func BuildUnavailableGUAMIList(guamiList []models.Guami) (unavailableGUAMIList ngapIE.UnavailableGUAMIList) {
	for _, guami := range guamiList {
		plmn := ngapConvert.PlmnIdToNgap(util.PlmnIdNidToModelsPlmnId(*guami.PlmnId))
		regionId, setId, ptrId := ngapConvert.AmfIdToNgap(guami.AmfId)
		item := ngapIE.UnavailableGUAMIItem{
			GUAMI: &ngapIE.GUAMI{
				PLMNIdentity: &plmn,
				AMFRegionID:  &ngapIE.AMFRegionID{Value: regionId},
				AMFSetID:     &ngapIE.AMFSetID{Value: setId},
				AMFPointer:   &ngapIE.AMFPointer{Value: ptrId},
			},
		}
		// TODO: item.TimerApproachForGUAMIRemoval and item.BackupAMFName not support yet
		unavailableGUAMIList.List = append(unavailableGUAMIList.List, item)
	}
	return
}
