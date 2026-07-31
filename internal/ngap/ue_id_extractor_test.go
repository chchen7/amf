package ngap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/ngap/ie"
	ngapMessage "github.com/free5gc/ngap/message"
)

func TestExtractUEID_InitialContextSetupResponse(t *testing.T) {
	const amfUeNgapID int64 = 99999

	encoded, err := (&ngapMessage.InitialContextSetupResponse{
		AMFUENGAPID: &ie.AMFUENGAPID{Value: amfUeNgapID},
		RANUENGAPID: &ie.RANUENGAPID{Value: 12345},
	}).MarshalBinary()
	require.NoError(t, err)

	ueID, found := ExtractUEID(encoded)
	assert.True(t, found)
	assert.Equal(t, uint64(amfUeNgapID), ueID)
}

func TestExtractUEID_InitialUEMessage(t *testing.T) {
	const ranUeNgapID int64 = 12345
	message := BuildInitialUEMessage(t, ranUeNgapID, []byte{0x7e, 0, 0x41}, "")
	encoded, err := message.MarshalBinary()
	require.NoError(t, err)

	ueID, found := ExtractUEID(encoded)
	assert.True(t, found)
	assert.Equal(t, uint64(ranUeNgapID), ueID)
}

func TestExtractUEID_UplinkNASTransportPrefersAMFID(t *testing.T) {
	const amfUeNgapID int64 = 67890
	initial := BuildInitialUEMessage(t, 12345, []byte{0x7e, 0, 0x41}, "")
	encoded, err := (&ngapMessage.UplinkNASTransport{
		AMFUENGAPID:             &ie.AMFUENGAPID{Value: amfUeNgapID},
		RANUENGAPID:             initial.RANUENGAPID,
		NASPDU:                  initial.NASPDU,
		UserLocationInformation: initial.UserLocationInformation,
	}).MarshalBinary()
	require.NoError(t, err)

	ueID, found := ExtractUEID(encoded)
	assert.True(t, found)
	assert.Equal(t, uint64(amfUeNgapID), ueID)
}

func TestExtractUEID_PDUSessionResourceSetupResponse(t *testing.T) {
	const amfUeNgapID int64 = 22222
	encoded, err := (&ngapMessage.PDUSessionResourceSetupResponse{
		AMFUENGAPID: &ie.AMFUENGAPID{Value: amfUeNgapID},
		RANUENGAPID: &ie.RANUENGAPID{Value: 12345},
	}).MarshalBinary()
	require.NoError(t, err)

	ueID, found := ExtractUEID(encoded)
	assert.True(t, found)
	assert.Equal(t, uint64(amfUeNgapID), ueID)
}

func TestExtractUEID_UEContextReleaseRequest(t *testing.T) {
	const amfUeNgapID int64 = 33333
	encoded, err := (&ngapMessage.UEContextReleaseRequest{
		AMFUENGAPID: &ie.AMFUENGAPID{Value: amfUeNgapID},
		RANUENGAPID: &ie.RANUENGAPID{Value: 12345},
		Cause: &ie.Cause{
			Choice: &ie.CauseRadioNetwork{Value: ie.CauseRadioNetworkPresentUserInactivity},
		},
	}).MarshalBinary()
	require.NoError(t, err)

	ueID, found := ExtractUEID(encoded)
	assert.True(t, found)
	assert.Equal(t, uint64(amfUeNgapID), ueID)
}

func TestNGAPIntegerIEValue_PriorityAndFallback(t *testing.T) {
	t.Run("source AMF UE ID", func(t *testing.T) {
		id, found := ngapIntegerIEValue(&ngapMessage.PathSwitchRequest{
			SourceAMFUENGAPID: &ie.AMFUENGAPID{Value: 11111},
		}, "SourceAMFUENGAPID")
		assert.True(t, found)
		assert.Equal(t, uint64(11111), id)
	})

	t.Run("RAN UE fallback", func(t *testing.T) {
		id, found := ngapIntegerIEValue(&ngapMessage.InitialContextSetupResponse{
			RANUENGAPID: &ie.RANUENGAPID{Value: 12345},
		}, "RANUENGAPID")
		assert.True(t, found)
		assert.Equal(t, uint64(12345), id)
	})
}

func TestExtractUEID_MessageWithoutUEID(t *testing.T) {
	encoded, err := (&ngapMessage.OverloadStop{}).MarshalBinary()
	require.NoError(t, err)

	ueID, found := ExtractUEID(encoded)
	assert.False(t, found)
	assert.Zero(t, ueID)
}

func TestExtractUEID_InvalidMessage(t *testing.T) {
	ueID, found := ExtractUEID([]byte{0xff, 0xff, 0xff, 0xff})
	assert.False(t, found)
	assert.Zero(t, ueID)
}
