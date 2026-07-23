//nolint:lll // Golden byte fixtures are intentionally kept as contiguous protocol values.
package message

import (
	"encoding/hex"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/fsm"
)

// These cases deliberately use no NAS security context.  They lock the AMF
// builder's plain-NAS baseline before the nas module is upgraded.
func TestBuildIdentityRequestSUCIGolden(t *testing.T) {
	got, err := BuildIdentityRequest(
		&context.AmfUe{},
		models.AccessType__3_GPP_ACCESS,
		nasMessage.MobileIdentity5GSTypeSuci,
	)
	require.NoError(t, err)
	require.Equal(t, []byte{0x7e, 0x00, 0x5b, 0x01}, got)

	decoded := nas.NewMessage()
	require.NoError(t, decoded.PlainNasDecode(&got))
	require.Equal(t, nas.MsgTypeIdentityRequest, decoded.GmmHeader.GetMessageType())
	require.Equal(t, nasMessage.MobileIdentity5GSTypeSuci,
		decoded.GmmMessage.IdentityRequest.SpareHalfOctetAndIdentityType.GetTypeOfIdentity())
}

func TestBuildAuthenticationRequest5GAKAGolden(t *testing.T) {
	ue := &context.AmfUe{
		NgKsi: models.NgKsi{Tsc: models.ScType_NATIVE, Ksi: 1},
		ABBA:  []uint8{0x00, 0x00},
		AuthenticationCtx: &models.UeAuthenticationCtx{
			AuthType: models.AusfUeAuthenticationAuthType__5_G_AKA,
			Var5gAuthData: models.Av5gAka{
				Rand: "000102030405060708090a0b0c0d0e0f",
				Autn: "101112131415161718191a1b1c1d1e1f",
			},
		},
	}

	got, err := BuildAuthenticationRequest(ue, models.AccessType__3_GPP_ACCESS)
	require.NoError(t, err)
	require.Equal(t, []byte{
		0x7e, 0x00, 0x56, 0x01, 0x02, 0x00, 0x00, 0x21, 0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x20, 0x10, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}, got)

	decoded := nas.NewMessage()
	require.NoError(t, decoded.PlainNasDecode(&got))
	require.Equal(t, nas.MsgTypeAuthenticationRequest, decoded.GmmHeader.GetMessageType())
	require.Equal(t, uint8(1),
		decoded.GmmMessage.AuthenticationRequest.SpareHalfOctetAndNgksi.GetNasKeySetIdentifiler())
}

func TestAMFNASBuilderGoldenBaseline(t *testing.T) {
	previousConfig := factory.AmfConfig
	factory.AmfConfig = &factory.Config{}
	t.Cleanup(func() { factory.AmfConfig = previousConfig })

	builds := []struct {
		name  string
		build func(*context.AmfUe) ([]byte, error)
	}{
		{"dl_nas_transport", func(ue *context.AmfUe) ([]byte, error) {
			cause, timerUnit := uint8(0x07), uint8(0x05)
			return BuildDLNASTransport(ue, models.AccessType__3_GPP_ACCESS, 1, []byte{1, 2, 3}, 1, &cause, &timerUnit, 2)
		}},
		{"notification", func(ue *context.AmfUe) ([]byte, error) {
			return BuildNotification(ue, models.AccessType__3_GPP_ACCESS)
		}},
		{"service_accept", func(ue *context.AmfUe) ([]byte, error) {
			return BuildServiceAccept(ue, models.AccessType__3_GPP_ACCESS, nil, nil, nil, nil)
		}},
		{"authentication_reject", func(ue *context.AmfUe) ([]byte, error) {
			return BuildAuthenticationReject(ue, models.AccessType__3_GPP_ACCESS, "")
		}},
		{"authentication_result", func(ue *context.AmfUe) ([]byte, error) {
			return BuildAuthenticationResult(ue, models.AccessType__3_GPP_ACCESS, true, "AQID")
		}},
		{"service_reject", func(ue *context.AmfUe) ([]byte, error) {
			return BuildServiceReject(ue, models.AccessType__3_GPP_ACCESS, nil, 0x07)
		}},
		{"registration_reject", func(ue *context.AmfUe) ([]byte, error) {
			return BuildRegistrationReject(ue, models.AccessType__3_GPP_ACCESS, 0x0b, "")
		}},
		{"security_mode_command", func(ue *context.AmfUe) ([]byte, error) {
			return BuildSecurityModeCommand(ue, models.AccessType__3_GPP_ACCESS, true, "AQID")
		}},
		{"deregistration_request", func(ue *context.AmfUe) ([]byte, error) {
			return BuildDeregistrationRequest(&context.RanUe{AmfUe: ue}, 1, true, 0x07)
		}},
		{"deregistration_accept", func(ue *context.AmfUe) ([]byte, error) {
			return BuildDeregistrationAccept(ue, models.AccessType__3_GPP_ACCESS)
		}},
		{"registration_accept", func(ue *context.AmfUe) ([]byte, error) {
			return BuildRegistrationAccept(ue, models.AccessType__3_GPP_ACCESS, nil, nil, nil, nil)
		}},
		{"status_5gmm", func(ue *context.AmfUe) ([]byte, error) {
			return BuildStatus5GMM(ue, models.AccessType__3_GPP_ACCESS, 0x5f)
		}},
		{"configuration_update_command", func(ue *context.AmfUe) ([]byte, error) {
			payload, err, _ := BuildConfigurationUpdateCommand(ue, models.AccessType__3_GPP_ACCESS,
				&context.ConfigurationUpdateCommandFlags{NeedNetworkSlicingIndication: true})
			return payload, err
		}},
	}

	for _, test := range builds {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.build(newGoldenUE())
			require.NoError(t, err)
			require.Equal(t, amfNASGolden[test.name], hex.EncodeToString(got))
		})
	}
}

var amfNASGolden = map[string]string{
	"dl_nas_transport":             "7e0200000000007e0068010003010203120158073701a2",
	"notification":                 "7e0200000000007e006501",
	"service_accept":               "7e0200000000007e004e",
	"authentication_reject":        "7e0200000000007e0058",
	"authentication_result":        "7e0200000000007e005a01000301020338020000",
	"service_reject":               "7e0200000000007e004d07",
	"registration_reject":          "7e0200000000007e00440b",
	"security_mode_command":        "7e0300000000007e005d0001000000e136010078000301020338020000",
	"deregistration_request":       "7e0200000000007e0047055807",
	"deregistration_accept":        "7e0200000000007e0046",
	"registration_accept":          "7e0200000000007e00420101",
	"status_5gmm":                  "7e0200000000007e00645f",
	"configuration_update_command": "7e0200000000007e0054d191",
}

func newGoldenUE() *context.AmfUe {
	return &context.AmfUe{
		SecurityContextAvailable: true,
		CipheringAlg:             0,
		IntegrityAlg:             0,
		NgKsi:                    models.NgKsi{Tsc: models.ScType_NATIVE, Ksi: 1},
		ABBA:                     []uint8{0, 0},
		KnasEnc:                  [16]uint8{},
		KnasInt:                  [16]uint8{},
		UESecurityCapability:     nasType.UESecurityCapability{Buffer: []byte{0, 0}},
		State: map[models.AccessType]*fsm.State{
			models.AccessType__3_GPP_ACCESS:    fsm.NewState(context.Deregistered),
			models.AccessType_NON_3_GPP_ACCESS: fsm.NewState(context.Deregistered),
		},
		RanUe:            map[models.AccessType]*context.RanUe{},
		RegistrationArea: map[models.AccessType][]models.Tai{},
		AllowedNssai:     map[models.AccessType][]models.AllowedSnssai{},
		NASLog:           logrus.NewEntry(logrus.New()),
	}
}
