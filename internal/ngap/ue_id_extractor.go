package ngap

import (
	"reflect"

	"github.com/free5gc/amf/internal/logger"
	ngapMessage "github.com/free5gc/ngap/message"
)

// ExtractUEID performs lightweight NGAP message decoding to extract an AMF or
// RAN UE identifier. PR #32 models each procedure as a distinct message type;
// the identifiers are therefore direct optional fields rather than entries in
// a shared NGAPPDU ProtocolIE container. AMF-UE-NGAP-ID takes precedence.
func ExtractUEID(encoded []byte) (uint64, bool) {
	msg, err := ngapMessage.Parse(encoded)
	if err != nil {
		logger.NgapLog.Warnf("Failed to parse NGAP message for UE ID extraction: %v", err)
		return 0, false
	}
	if msg == nil {
		logger.NgapLog.Trace("NGAP message is nil")
		return 0, false
	}
	for _, fieldName := range []string{"AMFUENGAPID", "SourceAMFUENGAPID", "RANUENGAPID"} {
		if id, ok := ngapIntegerIEValue(msg, fieldName); ok {
			logger.NgapLog.Tracef("Extracted UE ID: %d", id)
			return id, true
		}
	}

	logger.NgapLog.Tracef("No UE ID in NGAP procedure code: %d", msg.ProcedureCode())
	return 0, false
}

// ngapIntegerIEValue reads one of the generated *ie.AMFUENGAPID or
// *ie.RANUENGAPID fields. Reflection is intentionally limited to this common
// generated shape: it avoids a fragile procedure-by-procedure type switch
// while preserving a nil check and requiring an integer Value field.
func ngapIntegerIEValue(msg ngapMessage.Message, fieldName string) (uint64, bool) {
	value := reflect.ValueOf(msg)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return 0, false
	}

	ie := value.Elem().FieldByName(fieldName)
	if !ie.IsValid() || ie.Kind() != reflect.Pointer || ie.IsNil() {
		return 0, false
	}

	id := ie.Elem().FieldByName("Value")
	if !id.IsValid() || !id.CanInt() {
		return 0, false
	}
	if id.Int() < 0 {
		return 0, false
	}
	return uint64(id.Int()), true
}
