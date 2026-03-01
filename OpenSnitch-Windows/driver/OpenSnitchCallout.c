#include <ntddk.h>
#include <wdf.h>
#include <fwpsk.h>
#include <fwpmk.h>

// Dummy GUIDs for callouts
DEFINE_GUID(OPENSNITCH_CALLOUT_V4,
0x12345678, 0x1234, 0x1234, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56);

// OpenSnitch WFP Callout Driver Skeleton
// This driver would register an ALE_AUTH_CONNECT callout
// to intercept outgoing network connections, suspend them,
// send a request to the user-mode Go daemon, and wait for a reply.

DRIVER_INITIALIZE DriverEntry;
EVT_WDF_DRIVER_DEVICE_ADD OpenSnitchEvtDeviceAdd;

/**
 * Initialize and create the WDF driver object for the OpenSnitch callout driver.
 *
 * @param DriverObject Pointer to the caller's driver object supplied by the system.
 * @param RegistryPath Pointer to the driver's registry path supplied by the system.
 * @returns NTSTATUS code: `STATUS_SUCCESS` on success; on failure returns the error status produced by `WdfDriverCreate`.
 */
NTSTATUS DriverEntry(_In_ PDRIVER_OBJECT DriverObject, _In_ PUNICODE_STRING RegistryPath) {
    NTSTATUS status;
    WDF_DRIVER_CONFIG config;

    WDF_DRIVER_CONFIG_INIT(&config, OpenSnitchEvtDeviceAdd);

    status = WdfDriverCreate(DriverObject, RegistryPath, WDF_NO_OBJECT_ATTRIBUTES, &config, WDF_NO_HANDLE);
    if (!NT_SUCCESS(status)) {
        KdPrint(("OpenSnitch: WdfDriverCreate failed 0x%x\n", status));
        return status;
    }

    KdPrint(("OpenSnitch Callout Driver Loaded.\n"));

    // Real implementation would:
    // 1. Create control device / named pipe for Go daemon to connect.
    // 2. FwpmEngineOpen0 to connect to BFE.
    // 3. FwpsCalloutRegister0 to register our Callout functions (Classify, Notify, FlowDelete).
    // 4. FwpmCalloutAdd0 to add the callout to the system.

    return status;
}

/**
 * Create a WDF device object for the driver.
 *
 * @param DeviceInit Pointer to a pre-initialized WDFDEVICE_INIT structure used to create the device object.
 * @returns `STATUS_SUCCESS` if the device was created successfully, otherwise the NTSTATUS error code returned by WdfDeviceCreate.
 */
NTSTATUS OpenSnitchEvtDeviceAdd(_In_ WDFDRIVER Driver, _Inout_ PWDFDEVICE_INIT DeviceInit) {
    UNREFERENCED_PARAMETER(Driver);
    NTSTATUS status;
    WDFDEVICE device;

    status = WdfDeviceCreate(&DeviceInit, WDF_NO_OBJECT_ATTRIBUTES, &device);
    if (!NT_SUCCESS(status)) {
        KdPrint(("OpenSnitch: WdfDeviceCreate failed 0x%x\n", status));
        return status;
    }

    return status;
}

/**
 * Classify an incoming network event for the callout and produce a classification decision.
 *
 * In this skeleton implementation the function sets a default permit action and prevents
 * subsequent filters from overriding it. A full implementation would inspect metadata
 * (for example the process ID in `inMetaValues`) and decide to permit, block, or pend
 * the connection while awaiting user-mode approval.
 *
 * @param inMetaValues Metadata for the incoming event; contains the originating process ID
 *        that a real implementation would use to identify the requester.
 * @param classifyOut Structure to receive the classification result; this function sets
 *        the action and updates rights. The default behavior here sets the action to
 *        allow and clears the `FWPS_RIGHT_ACTION_WRITE` bit.
 */
void OpenSnitchClassify(
    const FWPS_INCOMING_VALUES0* inFixedValues,
    const FWPS_INCOMING_METADATA_VALUES0* inMetaValues,
    void* layerData,
    const void* classifyContext,
    const FWPS_FILTER0* filter,
    UINT64 flowContext,
    FWPS_CLASSIFY_OUT0* classifyOut)
{
    // Real implementation:
    // 1. Check if we have already evaluated this connection.
    // 2. Extract PID (inMetaValues->processId).
    // 3. Send event to Go daemon via inverted call or shared memory.
    // 4. Return FWPS_ACTION_BLOCK, or set action to FWPS_ACTION_CONTINUE and pend it.

    classifyOut->actionType = FWP_ACTION_PERMIT; // allow by default in skeleton
    classifyOut->rights &= ~FWPS_RIGHT_ACTION_WRITE;
}
