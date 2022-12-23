package tracerotelcol

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	// "strings"
	"time"

	"github.com/google/uuid"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.9.0"
	// "go.opentelemetry.io/otel/attribute"

	// "go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type Atm struct{
	ID            int64
	Version       string
	Name          string
	StateID       string
	SerialNumber  string
	ISPNetwork    string
}

type BackendSystem struct {
	Version       string
	ProcessName   string
	OSType        string
	OSVersion     string
	CloudProvider string
	CloudRegion   string
	ServiceName   string
	Endpoint      string
}

func generateAtm() Atm{
	i := getRandomNumber(1,2)
	var newAtm Atm

	switch i {
	    case 1: 
	        newAtm = Atm{
			    ID: 111,
			    Name: "ATM-111-RAJ",
			    SerialNumber: "atmxph-2022-111",
			    Version: "v1.0",
			    ISPNetwork: "SKJ-Udaipur",
			    StateID: "RAJ",

		    }
	    case 2:
		    newAtm = Atm{
			    ID: 444,
		        Name: "ATM-444-MAH",
			    SerialNumber: "atmxph-2022-444",
			    Version: "v1.0",
			    ISPNetwork: "SKJ-Mumbai",
			    StateID: "MAH",
		    }
    }

	return newAtm
}

func generateBackendSystem() BackendSystem{
    i := getRandomNumber(1,3)

	newBackend := BackendSystem{
		ProcessName: "accounts",
		Version: "v2.5",
		OSType: "lnx",
		OSVersion: "4.16.10-300.fc28.x86_64",
		CloudProvider: "amazon",
		CloudRegion: "us-east-4",
	}

	switch i {
	    case 1:
		    newBackend.Endpoint = "api/v2.5/balance"
	    case 2:
		    newBackend.Endpoint = "api/v2.5/deposite"
	    case 3:
		    newBackend.Endpoint = "api/v2.5/withdrawn"

	}

	return newBackend
}

func getRandomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	i := (rand.Intn(max - min + 1) + min)
	return i
}

func generateTraces() ptrace.Traces{
	traces := ptrace.NewTraces()


	numeberOfTraces := 0
	for i := 0; i <= numeberOfTraces; i++{
		newAtm := generateAtm()
		newBackendSystem := generateBackendSystem()

		resourceSpan := traces.ResourceSpans().AppendEmpty()
		atmResource := resourceSpan.Resource()
		fillResourceWithAtm(&atmResource, newAtm)

		atmInstScope := appendAtmSystemInstrScopeSpans(&resourceSpan)

		resourceSpan = traces.ResourceSpans().AppendEmpty()
		backendResource := resourceSpan.Resource()
		fillResourceWithBackendSystem(&backendResource, newBackendSystem)

		backendInstScope := appendAtmSystemInstrScopeSpans(&resourceSpan)


		appendTraceSpans(&newBackendSystem, &backendInstScope, &atmInstScope)
	}

	return traces
}

func fillResourceWithAtm(resource *pcommon.Resource, atm Atm){
	atmAttrs := resource.Attributes()
	atmAttrs.PutInt("atm.id", atm.ID)
	atmAttrs.PutStr("atm.stateid", atm.StateID)
	atmAttrs.PutStr("atm.ispnetwork", atm.ISPNetwork)
	atmAttrs.PutStr("atm.serialnumber", atm.SerialNumber)
	atmAttrs.PutStr(conventions.AttributeServiceName, atm.Name)
	atmAttrs.PutStr(conventions.AttributeServiceVersion, atm.Version)
 
 }

func fillResourceWithBackendSystem(resource *pcommon.Resource, backend BackendSystem){
	backendAttrs := resource.Attributes()
	var osType, cloudProvider string

	switch {
		case backend.CloudProvider == "amzn":
			cloudProvider = conventions.AttributeCloudProviderAWS
		case backend.OSType == "mcrsft":
			cloudProvider = conventions.AttributeCloudProviderAzure
		case backend.OSType == "gogl":
			cloudProvider = conventions.AttributeCloudProviderGCP
	}

	backendAttrs.PutStr(conventions.AttributeCloudProvider, cloudProvider)
	backendAttrs.PutStr(conventions.AttributeCloudRegion, backend.CloudRegion)

	switch {
		case backend.OSType == "lnx":
			osType = conventions.AttributeOSTypeLinux
		case backend.OSType == "wndws":
			osType = conventions.AttributeOSTypeWindows
		case backend.OSType == "slrs":
			osType = conventions.AttributeOSTypeSolaris
	}

	backendAttrs.PutStr(conventions.AttributeOSType, osType)
	backendAttrs.PutStr(conventions.AttributeOSVersion, backend.OSVersion)
 
	backendAttrs.PutStr(conventions.AttributeServiceName, backend.ProcessName)
	backendAttrs.PutStr(conventions.AttributeServiceVersion, backend.Version)

}

 func appendAtmSystemInstrScopeSpans(resourceSpans *ptrace.ResourceSpans) (ptrace.ScopeSpans){
	scopeSpans := resourceSpans.ScopeSpans().AppendEmpty()
	scopeSpans.Scope().SetName("atm-system")
	scopeSpans.Scope().SetVersion("v1.0")
	return scopeSpans
}

func NewTraceID() pcommon.TraceID{
	return pcommon.TraceID(uuid.New())
}

func NewSpanID() pcommon.SpanID {
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xff, 0x01, 0x02, 0x03, 0xbe, 0xef}
	r := bytes.NewReader(b)
	var rngSeed int64
	_ = binary.Read(r, binary.LittleEndian, &rngSeed)
	randSource := rand.New(rand.NewSource(rngSeed))

	var sid [8]byte
	randSource.Read(sid[:])
    spanID := pcommon.SpanID(sid)

	return spanID
}

func appendTraceSpans(backend *BackendSystem, backendScopeSpans *ptrace.ScopeSpans, atmScopeSpans *ptrace.ScopeSpans){
	// traceId := NewTraceID()

// 	var atmOperationName string

// 	switch {
// 	case strings.Contains(backend.Endpoint, "balance"):
//         atmOperationName = "Check Balance"
// 	case strings.Contains(backend.Endpoint, "deposit"):
// 		atmOperationName = "Make Deposit"
// 	case strings.Contains(backend.Endpoint, "withdraw"):
// 		atmOperationName = "Fast Cash"
// 	}

// 	atmSpanId := NewSpanID()
//     atmSpanStartTime := time.Now()
//     atmDuration, _ := time.ParseDuration("4s")
//     atmSpanFinishTime := atmSpanStartTime.Add(atmDuration)


// 	atmSpan := atmScopeSpans.Spans().AppendEmpty()
// 	atmSpan.SetTraceID(traceId)
// 	atmSpan.SetSpanID(atmSpanId)
// 	atmSpan.SetName(atmOperationName)
// 	atmSpan.SetKind(ptrace.SpanKindClient)
// 	atmSpan.Status().SetCode(ptrace.StatusCodeOk)
// 	atmSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(atmSpanStartTime))
// 	atmSpan.SetEndTimestamp(pcommon.NewTimestampFromTime(atmSpanFinishTime))


// 	backendSpanId := NewSpanID()

// 	backendDuration, _ := time.ParseDuration("2s")
//     backendSpanStartTime := atmSpanStartTime.Add(backendDuration)


// 	backendSpan := backendScopeSpans.Spans().AppendEmpty()
// 	backendSpan.SetTraceID(atmSpan.TraceID())
// 	backendSpan.SetSpanID(backendSpanId)
// 	backendSpan.SetParentSpanID(atmSpan.SpanID())
// 	backendSpan.SetName(backend.Endpoint)
// 	backendSpan.SetKind(ptrace.SpanKindServer)
// 	backendSpan.Status().SetCode(ptrace.StatusCodeOk)
// 	backendSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(backendSpanStartTime))
// 	backendSpan.SetEndTimestamp(atmSpan.EndTimestamp())

// }


    traceId := NewTraceID()
	backendSpanId := NewSpanID()

	backendDuration, _ := time.ParseDuration("1s")
    backendSpanStartTime := time.Now()
    backendSpanFinishTime := backendSpanStartTime.Add(backendDuration)


	backendSpan := backendScopeSpans.Spans().AppendEmpty()
	backendSpan.SetTraceID(traceId)
	backendSpan.SetSpanID(backendSpanId)
	backendSpan.SetName(backend.Endpoint)
	backendSpan.SetKind(ptrace.SpanKindServer)
	backendSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(backendSpanStartTime))
	backendSpan.SetEndTimestamp(pcommon.NewTimestampFromTime(backendSpanFinishTime))

}

