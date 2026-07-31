package testing

import (
	"encoding/base64"

	"github.com/free5gc/nas/ie"
	"github.com/free5gc/nas/message"
	"github.com/free5gc/openapi/models"
)

const (
	PDUSesModiReq    string = "PDU Session Modification Request"
	PDUSesModiCmp    string = "PDU Session Modification Complete"
	PDUSesModiCmdRej string = "PDU Session Modification Command Reject"
	PDUSesRelReq     string = "PDU Session Release Request"
	PDUSesRelCmp     string = "PDU Session Release Complete"
	PDUSesRelRej     string = "PDU Session Release Reject"
	PDUSesAuthCmp    string = "PDU Session Authentication Complete"
)

func marshal(msg message.Message) []byte {
	data, err := msg.MarshalBinary()
	if err != nil {
		return nil
	}
	return data
}

func psiFromBytes(data []byte) ie.Psi {
	var psi ie.Psi
	_ = psi.UnmarshalBinary(data)
	return psi
}

func GetRegistrationRequest(
	registrationType uint8,
	mobileIdentity ie.MobileId5GS,
	requestedNSSAI *ie.NSSAI,
	ueSecurityCapability *ie.UESecCapability,
	capability5GMM *ie.Capability5GMM,
	nasMessageContainer []uint8,
	uplinkDataStatus *ie.UplinkDataStatus,
) []byte {
	registrationRequest := &message.RegReq{
		RegType5GS:       &ie.RegType5GS{FOR_Pending: true, Value: registrationType},
		Ngksi:            &ie.NASKeySetId{Tsc: ie.SecCtxTypeNative, Ksi: ie.NASKeyNA},
		MobileId5GS:      &mobileIdentity,
		ReqNSSAI:         requestedNSSAI,
		UESecCapability:  ueSecurityCapability,
		Capability5GMM:   capability5GMM,
		UplinkDataStatus: uplinkDataStatus,
	}
	if nasMessageContainer != nil {
		registrationRequest.NASMsgCntr = &ie.NASMsgCntr{Contents: nasMessageContainer}
	}
	return marshal(registrationRequest)
}

func GetPduSessionEstablishmentRequest(pduSessionId uint8) []byte {
	return marshal(&message.PDUSessEstReq{
		PDUSessId: pduSessionId,
		IntegrityProtectionMaxDataRate: &ie.IntegrityProtectionMaxDataRate{
			Uplink: 0xff, Downlink: 0xff,
		},
		PDUSessType: &ie.PDUSessType{Value: ie.PDUSessType_IPv4},
	})
}

func newULNASTransport(payload []byte, pduSessionId uint8, requestType *uint8, dnn string,
	sNssai *models.Snssai,
) []byte {
	ulNasTransport := &message.ULNASTransport{
		PayloadCntrType: &ie.PayloadCntrType{Value: ie.PayloadCntrType_N1SMInfo},
		PayloadCntr:     &ie.PayloadCntr{Pct: ie.PayloadCntrType_N1SMInfo, Contents: payload},
		PDUSessID:       &ie.PDUSessId2{Value: pduSessionId},
	}
	if requestType != nil {
		ulNasTransport.ReqType = &ie.ReqType{Value: ie.ConstReqType(*requestType)}
	}
	if dnn != "" {
		ulNasTransport.DNN = &ie.DNN{Value: dnn}
	}
	if sNssai != nil {
		ulNasTransport.SNSSAI = &ie.SNSSAI{SST: uint8(sNssai.Sst), SD: sNssai.Sd}
	}
	return marshal(ulNasTransport)
}

func GetUlNasTransport_PduSessionEstablishmentRequest(pduSessionId uint8, requestType uint8, dnn string,
	sNssai *models.Snssai,
) []byte {
	return newULNASTransport(
		GetPduSessionEstablishmentRequest(pduSessionId), pduSessionId, &requestType, dnn, sNssai)
}

func GetUlNasTransport_PduSessionModificationRequest(pduSessionId uint8, requestType uint8, dnn string,
	sNssai *models.Snssai,
) []byte {
	return newULNASTransport(
		GetPduSessionModificationRequest(pduSessionId), pduSessionId, &requestType, dnn, sNssai)
}

func GetPduSessionModificationRequest(pduSessionId uint8) []byte {
	return marshal(&message.PDUSessModReq{PDUSessId: pduSessionId})
}

func GetPduSessionModificationComplete(pduSessionId uint8) []byte {
	return marshal(&message.PDUSessModComplete{PDUSessId: pduSessionId})
}

func GetPduSessionModificationCommandReject(pduSessionId uint8) []byte {
	return marshal(&message.PDUSessModCmdRej{
		PDUSessId: pduSessionId,
		Cause5GSM: &ie.Cause5GSM{},
	})
}

func GetPduSessionReleaseRequest(pduSessionId uint8) []byte {
	return marshal(&message.PDUSessRelReq{PDUSessId: pduSessionId})
}

func GetPduSessionReleaseComplete(pduSessionId uint8) []byte {
	return marshal(&message.PDUSessRelComplete{PDUSessId: pduSessionId})
}

func GetPduSessionReleaseReject(pduSessionId uint8) []byte {
	return marshal(&message.PDUSessRelRej{
		PDUSessId: pduSessionId,
		Cause5GSM: &ie.Cause5GSM{},
	})
}

func GetPduSessionAuthenticationComplete(pduSessionId uint8) []byte {
	return marshal(&message.PDUSessAuthComplete{
		PDUSessId: pduSessionId,
		EAPMsg:    &ie.EAPMsg{Eap: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}},
	})
}

func GetUlNasTransport_PduSessionCommonData(pduSessionId uint8, types string) []byte {
	var payload []byte
	switch types {
	case PDUSesModiReq:
		payload = GetPduSessionModificationRequest(pduSessionId)
	case PDUSesModiCmp:
		payload = GetPduSessionModificationComplete(pduSessionId)
	case PDUSesModiCmdRej:
		payload = GetPduSessionModificationCommandReject(pduSessionId)
	case PDUSesRelReq:
		payload = GetPduSessionReleaseRequest(pduSessionId)
	case PDUSesRelCmp:
		payload = GetPduSessionReleaseComplete(pduSessionId)
	case PDUSesRelRej:
		payload = GetPduSessionReleaseReject(pduSessionId)
	case PDUSesAuthCmp:
		payload = GetPduSessionAuthenticationComplete(pduSessionId)
	}
	return newULNASTransport(payload, pduSessionId, nil, "", nil)
}

func GetIdentityResponse(mobileIdentity ie.MobileId5GS) []byte {
	return marshal(&message.IdRsp{MobileId: &mobileIdentity})
}

func GetNotificationResponse(pDUSessionStatus []uint8) []byte {
	return marshal(&message.NotifRsp{
		PDUSessStatus: &ie.PDUSessStatus{Psi: psiFromBytes(pDUSessionStatus)},
	})
}

func GetConfigurationUpdateComplete() []byte {
	return marshal(&message.CfgUpdateComplete{})
}

func GetServiceRequest(serviceType uint8) []byte {
	serviceRequest := &message.SvcReq{
		Ngksi:   &ie.NASKeySetId{Tsc: ie.SecCtxTypeNative, Ksi: 1},
		SvcType: &ie.SvcType{Value: ie.ConstSvcType(serviceType)},
		TMSI5GS: &ie.MobileId5GS{
			TypeOfId: ie.IdType_5GS_TMSI,
			AMFSetID: uint16(0xFE) << 2,
			TMSI5G:   [4]byte{0, 0, 0, 1},
		},
	}
	switch serviceType {
	case uint8(ie.SvcType_MobileTermSvc):
		serviceRequest.AllowedPDUSessStatus = &ie.AllowedPDUSessStatus{
			Psi: psiFromBytes([]byte{0x00, 0x08}),
		}
	case uint8(ie.SvcType_Data):
		serviceRequest.UplinkDataStatus = &ie.UplinkDataStatus{
			Psi: psiFromBytes([]byte{0x00, 0x04}),
		}
	}
	return marshal(serviceRequest)
}

func GetAuthenticationResponse(authenticationResponseParam []uint8, eapMsg string) []byte {
	authenticationResponse := &message.AuthRsp{}
	if len(authenticationResponseParam) > 0 {
		authenticationResponse.AuthRspParam = &ie.AuthRspParam{Res: authenticationResponseParam}
	} else if eapMsg != "" {
		rawEapMsg, err := base64.StdEncoding.DecodeString(eapMsg)
		if err == nil {
			authenticationResponse.EAPMsg = &ie.EAPMsg{Eap: rawEapMsg}
		}
	}
	return marshal(authenticationResponse)
}

func GetAuthenticationFailure(cause5GMM uint8, authenticationFailureParam []uint8) []byte {
	authenticationFailure := &message.AuthFailure{
		Cause5GMM: &ie.Cause5GMM{Value: cause5GMM},
	}
	if cause5GMM == ie.Cause5GMM_SynchFailure {
		authenticationFailure.AuthFailureParam = &ie.AuthFailureParam{Value: authenticationFailureParam}
	}
	return marshal(authenticationFailure)
}

func GetRegistrationComplete(sorTransparentContainer []uint8) []byte {
	if sorTransparentContainer == nil {
		return marshal(&message.RegComplete{})
	}
	// SORTransparentCntr is not yet structurally decoded by the new NAS IE package.
	// Preserve the caller-provided transparent bytes in their wire representation.
	data := []byte{byte(message.Epd5GSMobilityMgmtMsg), 0, byte(message.MsgTypeRegComplete), 0x73,
		byte(len(sorTransparentContainer) >> 8), byte(len(sorTransparentContainer))}
	return append(data, sorTransparentContainer...)
}

// TS 24.501 8.2.26.
func GetSecurityModeComplete(nasMessageContainer []uint8) []byte {
	securityModeComplete := &message.SecModeComplete{
		IMEISV: &ie.MobileId5GS{
			TypeOfId:     ie.IdType_5GS_IMEISV,
			OddEvenIndic: ie.EvenNumOfIdDigit,
			IMEISV:       [16]uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	if nasMessageContainer != nil {
		securityModeComplete.NASMsgCntr = &ie.NASMsgCntr{Contents: nasMessageContainer}
	}
	return marshal(securityModeComplete)
}

func GetSecurityModeReject(cause5GMM uint8) []byte {
	return marshal(&message.SecModeRej{Cause5GMM: &ie.Cause5GMM{Value: cause5GMM}})
}

func GetDeregistrationRequest(accessType uint8, switchOff uint8, ngKsi uint8,
	mobileIdentity5GS ie.MobileId5GS,
) []byte {
	return marshal(&message.DeregReqUEOrig{
		DeregType: &ie.DeregType{
			AccessType: accessType,
			Switchoff:  switchOff != 0,
		},
		Ngksi:       &ie.NASKeySetId{Tsc: ie.SecCtxType(ngKsi), Ksi: ngKsi},
		MobileId5GS: &mobileIdentity5GS,
	})
}

func GetDeregistrationAccept() []byte {
	return marshal(&message.DeregAcceptUETerm{})
}

func GetStatus5GMM(cause uint8) []byte {
	return marshal(&message.Status5GMM{Cause5GMM: &ie.Cause5GMM{Value: cause}})
}

func GetStatus5GSM(pduSessionId uint8, cause uint8) []byte {
	return marshal(&message.Status5GSM{
		PDUSessId: pduSessionId,
		Cause5GSM: &ie.Cause5GSM{Value: cause},
	})
}

func GetUlNasTransport_Status5GSM(pduSessionId uint8, cause uint8) []byte {
	return newULNASTransport(GetStatus5GSM(pduSessionId, cause), pduSessionId, nil, "", nil)
}

func GetUlNasTransport_PduSessionReleaseRequest(pduSessionId uint8) []byte {
	return newULNASTransport(GetPduSessionReleaseRequest(pduSessionId), pduSessionId, nil, "", nil)
}

func GetUlNasTransport_PduSessionReleaseComplete(pduSessionId uint8, requestType uint8, dnn string,
	sNssai *models.Snssai,
) []byte {
	return newULNASTransport(
		GetPduSessionReleaseComplete(pduSessionId), pduSessionId, &requestType, dnn, sNssai)
}
