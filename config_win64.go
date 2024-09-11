//go:build windows && amd64

package iec61850

/*
	-l会在前面加个lib
*/

// #cgo CFLAGS: -I./libiec61850/inc/hal/inc -I./libiec61850/inc/common/inc -I./libiec61850/inc/goose -I./libiec61850/inc/iec61850/inc -I./libiec61850/inc/iec61850/inc_private -I./libiec61850/inc/logging -I./libiec61850/inc/mms/inc -I./libiec61850/inc/mms/inc_private -I./libiec61850/inc/mms/iso_mms/asn1c
// #cgo CFLAGS: -I./model
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -L./libiec61850/lib/win64 -L./model -liec61850 -lstatic_model -lws2_32
import "C"
