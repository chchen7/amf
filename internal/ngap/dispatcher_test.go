package ngap

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	ngaptesting "github.com/free5gc/amf/internal/ngap/testing"
	"github.com/free5gc/ngap/aper"
	"github.com/free5gc/ngap/ie"
	ngapMessage "github.com/free5gc/ngap/message"
)

func TestHandleDecodeErrorSendsErrorIndication(t *testing.T) {
	conn := new(ngaptesting.SctpConnStub)
	ran := NewAmfRan(conn)
	decoded := &ngapMessage.UplinkNASTransport{
		AMFUENGAPID: &ie.AMFUENGAPID{Value: 11},
		RANUENGAPID: &ie.RANUENGAPID{Value: 22},
	}
	reportIE := ie.BuildAbstractSyntaxErrReportIe(ie.ProtocolIEIDNASPDU, ie.CriticalityReject)
	decodeErr := ie.BuildAbstractSyntaxErr(
		ngapMessage.ProcedureCodeUplinkNASTransport,
		aper.Enumerated(ngapMessage.MessageTypeInitiatingMessage),
		ie.CriticalityReject,
		&ie.AbstractSyntaxErrMissingIE{ReportIe: reportIE},
		errors.New("missing NAS-PDU"),
	)

	handleDecodeError(ran, decoded, decodeErr)
	require.Len(t, conn.MsgList, 1)

	response, err := ngapMessage.Parse(conn.MsgList[0])
	require.NoError(t, err)
	errorIndication, ok := response.(*ngapMessage.ErrorIndication)
	require.True(t, ok)
	require.Equal(t, int64(11), errorIndication.AMFUENGAPID.Value)
	require.Equal(t, int64(22), errorIndication.RANUENGAPID.Value)
	cause, ok := errorIndication.Cause.Choice.(*ie.CauseProtocol)
	require.True(t, ok)
	require.Equal(t, ie.CauseProtocolPresentAbstractSyntaxErrorReject, cause.Value)
	require.NotNil(t, errorIndication.CriticalityDiagnostics)
}
