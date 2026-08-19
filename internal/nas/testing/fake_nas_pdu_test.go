package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/nas/ie"
	"github.com/free5gc/nas/message"
	"github.com/free5gc/openapi/models"
)

func TestFakeNASPDUsUseNewMessageAPI(t *testing.T) {
	var mobileIdentity ie.MobileId5GS
	require.NoError(t, mobileIdentity.UnmarshalBinary([]byte{
		0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0, 0, 0, 0, 0x47, 0x78,
	}))
	snssai := &models.Snssai{Sst: 1, Sd: "010203"}

	testCases := map[string]func() []byte{
		"registration request": func() []byte {
			return GetRegistrationRequest(
				ie.RegType_InitialReg,
				mobileIdentity,
				nil,
				&ie.UESecCapability{Length: 2, EA05G: true, IA2_128_5G: true},
				nil,
				nil,
				nil,
			)
		},
		"PDU session establishment request": func() []byte {
			return GetPduSessionEstablishmentRequest(1)
		},
		"UL NAS PDU session establishment": func() []byte {
			return GetUlNasTransport_PduSessionEstablishmentRequest(
				1, uint8(ie.ReqType_InitialReq), "internet", snssai)
		},
		"UL NAS PDU session modification": func() []byte {
			return GetUlNasTransport_PduSessionModificationRequest(
				1, uint8(ie.ReqType_ModReq), "internet", snssai)
		},
		"PDU session modification request": func() []byte {
			return GetPduSessionModificationRequest(1)
		},
		"PDU session modification complete": func() []byte {
			return GetPduSessionModificationComplete(1)
		},
		"PDU session modification command reject": func() []byte {
			return GetPduSessionModificationCommandReject(1)
		},
		"PDU session release request": func() []byte {
			return GetPduSessionReleaseRequest(1)
		},
		"PDU session release complete": func() []byte {
			return GetPduSessionReleaseComplete(1)
		},
		"PDU session release reject": func() []byte {
			return GetPduSessionReleaseReject(1)
		},
		"PDU session authentication complete": func() []byte {
			return GetPduSessionAuthenticationComplete(1)
		},
		"UL NAS common data": func() []byte {
			return GetUlNasTransport_PduSessionCommonData(1, PDUSesModiReq)
		},
		"identity response": func() []byte {
			return GetIdentityResponse(mobileIdentity)
		},
		"notification response": func() []byte {
			return GetNotificationResponse([]byte{0, 1})
		},
		"configuration update complete": GetConfigurationUpdateComplete,
		"service request": func() []byte {
			return GetServiceRequest(uint8(ie.SvcType_Data))
		},
		"authentication response": func() []byte {
			return GetAuthenticationResponse(make([]byte, 16), "")
		},
		"authentication failure": func() []byte {
			return GetAuthenticationFailure(ie.Cause5GMM_SynchFailure, make([]byte, 14))
		},
		"registration complete": func() []byte {
			return GetRegistrationComplete(nil)
		},
		"security mode complete": func() []byte {
			return GetSecurityModeComplete(nil)
		},
		"security mode reject": func() []byte {
			return GetSecurityModeReject(ie.Cause5GMM_ProtError)
		},
		"deregistration request": func() []byte {
			return GetDeregistrationRequest(
				ie.AccessType_3gpp, 0, ie.NASKeyNA, mobileIdentity)
		},
		"deregistration accept": GetDeregistrationAccept,
		"5GMM status": func() []byte {
			return GetStatus5GMM(ie.Cause5GMM_ProtError)
		},
		"5GSM status": func() []byte {
			return GetStatus5GSM(1, 0x1f)
		},
		"UL NAS 5GSM status": func() []byte {
			return GetUlNasTransport_Status5GSM(1, 0x1f)
		},
		"UL NAS PDU session release request": func() []byte {
			return GetUlNasTransport_PduSessionReleaseRequest(1)
		},
		"UL NAS PDU session release complete": func() []byte {
			return GetUlNasTransport_PduSessionReleaseComplete(
				1, uint8(ie.ReqType_ExistingPDUSess), "internet", snssai)
		},
	}

	for name, build := range testCases {
		t.Run(name, func(t *testing.T) {
			payload := build()
			require.NotEmpty(t, payload)
			decoded, err := message.Parse(payload, nil)
			require.NoError(t, err)
			require.NotNil(t, decoded)
		})
	}
}
