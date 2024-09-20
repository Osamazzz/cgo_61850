package iec61850

// #include "static_model.h"
// extern IedModel iedModel;
// #include <iec61850_server.h>
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

type IedModel struct {
	model *C.IedModel
}

type LogicalDevice struct {
	device *C.LogicalDevice
}

func NewIedModel(name string) *IedModel {
	cname := C.CString(name)
	// 手动释放内存
	defer C.free(unsafe.Pointer(cname))
	return &IedModel{
		model: C.IedModel_create(cname),
	}
}

func GetModel() *IedModel {
	return &IedModel{
		model: &C.iedModel,
	}
}

// GetFirstChild 获取IedModel的第一个LogicalDevice ,子节点firstchild是逻辑设备类型
func (im *IedModel) GetFirstChild() *LogicalDevice {
	// 检查 im.model 是否为 nil
	if im.model == nil {
		fmt.Println("im.model is nil, cannot get firstChild")
		return nil
	}

	// 获取 firstChild，确保类型匹配
	// 通过unsafe.pointer来访问这些结构体中的成员，再强转为相应类型，这很重要
	firstChild := (*C.LogicalDevice)(unsafe.Pointer(im.model.firstChild))

	// 检查 firstChild 是否为 nil
	if firstChild == nil {
		fmt.Println("firstChild is nil")
		return nil
	}

	// 返回 LogicalDevice 的 Go 包装类型
	return &LogicalDevice{
		device: firstChild,
	}
}

func (ln *LogicalNode) Getsibling() *LogicalNode {
	// var ld_firstChild *LogicalNode
	//	ld_firstChild = (*LogicalNode)(unsafe.Pointer(ld.device.firstChild))
	return &LogicalNode{
		(*C.LogicalNode)(unsafe.Pointer(ln.node.sibling)),
	}
}

//func (do *DataObject) GetChild(name string) *DataAttribute {
//	cname := C.CString(name)
//	defer C.free(unsafe.Pointer(cname))
//	return &DataAttribute{
//		attribute: (*C.DataAttribute)(unsafe.Pointer(C.ModelNode_getChild())),
//	}
//}

func (m *IedModel) Destroy() {
	C.IedModel_destroy(m.model)
}

func CreateModelFromConfigFileEx(filepath string) (*IedModel, error) {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	cFilepath := C.CString(filepath)
	// 释放内存
	defer C.free(unsafe.Pointer(cFilepath))
	model := &IedModel{
		model: C.ConfigFileParser_createModelFromConfigFileEx(cFilepath),
	}
	return model, nil
}

func (m *IedModel) CreateLogicalDevice(name string) *LogicalDevice {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &LogicalDevice{
		device: C.LogicalDevice_create(cname, m.model),
	}
}

type LogicalNode struct {
	node *C.LogicalNode
}

func (d *LogicalDevice) CreateLogicalNode(name string) *LogicalNode {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &LogicalNode{
		node: C.LogicalNode_create(cname, d.device),
	}
}

type DataObject struct {
	object *C.DataObject
}

// ENS: EnumerationString
// VSS: Visible String Setting
// SAV: Sampled Value
// APC: Analogue Process Control

func (n *LogicalNode) CreateDataObjectCDC_ENS(name string) *DataObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataObject{
		object: C.CDC_ENS_create(cname, (*C.ModelNode)(n.node), 0),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_VSS(name string) *DataObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataObject{
		object: C.CDC_VSS_create(cname, (*C.ModelNode)(n.node), 0),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_SAV(name string, isInteger bool) *DataObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataObject{
		object: C.CDC_SAV_create(cname, (*C.ModelNode)(n.node), 0, C.bool(isInteger)),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_APC(name string, ctlModel int) *DataObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataObject{
		object: C.CDC_APC_create(cname, (*C.ModelNode)(n.node), 0, C.uint(ctlModel), C.bool(false)),
	}
}

type DataAttribute struct {
	attribute *C.DataAttribute
}

func (do *DataObject) GetChild(name string) *DataAttribute {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataAttribute{
		attribute: (*C.DataAttribute)(unsafe.Pointer(C.ModelNode_getChild((*C.ModelNode)(unsafe.Pointer(do.object)), cname))),
	}
}

type DataSet struct {
	dataSet *C.DataSet
}

// CreateDataSet creates a new DataSet under this LogicalNode.
func (ln *LogicalNode) CreateDataSet(name string) *DataSet {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cDataSet := C.DataSet_create(cName, ln.node)
	return &DataSet{dataSet: cDataSet}
}

// AddDataSetEntry adds a new DataSetEntry to this DataSet.
func (ds *DataSet) AddDataSetEntry(ref string) {
	cRef := C.CString(ref)
	defer C.free(unsafe.Pointer(cRef))

	C.DataSetEntry_create(ds.dataSet, cRef, -1, nil)
}
