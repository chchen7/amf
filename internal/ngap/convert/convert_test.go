package convert

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"
)

func TestPLMNIdentityRoundTrip(t *testing.T) {
	for _, plmn := range []models.PlmnId{
		{Mcc: "208", Mnc: "93"},
		{Mcc: "310", Mnc: "260"},
	} {
		t.Run(plmn.Mcc+"-"+plmn.Mnc, func(t *testing.T) {
			require.Equal(t, plmn, PlmnIdToModels(PlmnIdToNgap(plmn)))
		})
	}
}

func TestAMFIDRoundTrip(t *testing.T) {
	region, set, pointer := AmfIdToNgap("cafe00")
	require.Equal(t, "cafe00", AmfIdToModels(region, set, pointer))
}

func TestIPAddressRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		ipv4 string
		ipv6 string
	}{
		{name: "IPv4", ipv4: "192.0.2.1"},
		{name: "IPv6", ipv6: "2001:db8::1"},
		{name: "IPv4 and IPv6", ipv4: "192.0.2.1", ipv6: "2001:db8::1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotIPv4, gotIPv6 := IPAddressToString(IPAddressToNgap(test.ipv4, test.ipv6))
			require.Equal(t, test.ipv4, gotIPv4)
			require.Equal(t, test.ipv6, gotIPv6)
		})
	}
}

func TestGlobalGNBIDRoundTrip(t *testing.T) {
	ranID := models.GlobalRanNodeId{
		PlmnId: &models.PlmnId{Mcc: "208", Mnc: "93"},
		GNbId:  &models.GNbId{BitLength: 24, GNBValue: "000102"},
	}
	encoded, err := RanIDToNgap(ranID)
	require.NoError(t, err)
	require.Equal(t, ranID, RanIdToModels(encoded))
}

func TestTraceDataToNgap(t *testing.T) {
	trace := TraceDataToNgap(models.TraceData{
		TraceRef:                 "20893-010203",
		TraceDepth:               models.TraceDepth_MEDIUM,
		InterfaceList:            "80",
		CollectionEntityIpv4Addr: "192.0.2.1",
		CollectionEntityIpv6Addr: "",
	}, "0405")

	require.NotNil(t, trace.NGRANTraceID)
	require.Equal(t, []byte{0x02, 0xf8, 0x39, 1, 2, 3, 4, 5}, []byte(trace.NGRANTraceID.Value))
	require.NotNil(t, trace.InterfacesToTrace)
	require.NotNil(t, trace.TraceDepth)
	require.NotNil(t, trace.TraceCollectionEntityIPAddress)
}
