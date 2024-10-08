package iec61850

// #include <iec61850_server.h>
import "C"
import "unsafe"

type IedServer struct {
	server C.IedServer
}

type MmsServer struct {
	mmsServer C.MmsServer
}

// GetMmsServer  直接从已经有的iedserver里面获取
func GetMmsServer(server *IedServer) *MmsServer {
	return &MmsServer{
		mmsServer: C.IedServer_getMmsServer(server.server),
	}
}

// 设置服务器的基础文件路径，用于处理 MMS 文件服务。
func setFileStoreBasePath(is *IedServer, path string) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	C.IedServer_setFilestoreBasepath(is.server, cpath)
}

// NewServer creates a new instance of the IedServer using the provided model.
func NewServer(model *IedModel) *IedServer {
	return &IedServer{
		server: C.IedServer_create(model.model),
	}
}

// Start initiates the IedServer on the provided port.
func (is *IedServer) Start(port int) {
	C.IedServer_start(is.server, C.int(port))
	// If there's another way to detect the error, handle it here.
}

// IsRunning checks if the IedServer is currently running.
func (is *IedServer) IsRunning() bool {
	return bool(C.IedServer_isRunning(is.server))
}

// Stop terminates the IedServer.
func (is *IedServer) Stop() {
	C.IedServer_stop(is.server)
}

// Destroy frees all resources associated with the IedServer.
func (is *IedServer) Destroy() {
	C.IedServer_destroy(is.server)
}

// LockDataModel locks the data model of the IedServer.
func (is *IedServer) LockDataModel() {
	C.IedServer_lockDataModel(is.server)
}

// UnlockDataModel unlocks the data model of the IedServer.
func (is *IedServer) UnlockDataModel() {
	C.IedServer_unlockDataModel(is.server)
}

// UpdateUTCTimeAttributeValue updates a DataAttribute with a UTC time value.
func (is *IedServer) UpdateUTCTimeAttributeValue(attr *DataAttribute, value int64) {
	C.IedServer_updateUTCTimeAttributeValue(is.server, attr.attribute, C.uint64_t(value))
}

// UpdateFloatAttributeValue updates a DataAttribute with a float value.
func (is *IedServer) UpdateFloatAttributeValue(attr *DataAttribute, value float32) {
	C.IedServer_updateFloatAttributeValue(is.server, attr.attribute, C.float(value))
}

// UpdateInt32AttributeValue updates a DataAttribute with an Int32 value.
func (is *IedServer) UpdateInt32AttributeValue(attr *DataAttribute, value int32) {
	C.IedServer_updateInt32AttributeValue(is.server, attr.attribute, C.int32_t(value))
}

func (is *IedServer) UpdateInt64AttributeValue(attr *DataAttribute, value int64) {
	C.IedServer_updateInt64AttributeValue(is.server, attr.attribute, C.int64_t(value))
}

// UpdateVisibleStringAttributeValue updates a DataAttribute with a visible string value.
func (is *IedServer) UpdateVisibleStringAttributeValue(attr *DataAttribute, value string) {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	C.IedServer_updateVisibleStringAttributeValue(is.server, attr.attribute, cvalue)
}
