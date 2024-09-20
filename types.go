package iec61850

import "C"

type MmsType int

type MmsValue struct {
	Type  MmsType
	Value interface{}
}

// data types
const (
	Array MmsType = iota
	Structure
	Boolean
	BitString
	Integer
	Unsigned
	Float
	OctetString
	VisibleString
	GeneralizedTime
	BinaryTime
	Bcd
	ObjId
	String
	UTCTime
	DataAccessError
	Int8
	Int16
	Int32
	Int64
	Uint8
	Uint16
	Uint32
)

type DataAttributeType int

const (
	IEC61850_UNKNOWN_TYPE       DataAttributeType = iota - 1 // -1
	IEC61850_BOOLEAN                                         // 0
	IEC61850_INT8                                            // 1
	IEC61850_INT16                                           // 2
	IEC61850_INT32                                           // 3
	IEC61850_INT64                                           // 4
	IEC61850_INT128                                          // 5
	IEC61850_INT8U                                           // 6
	IEC61850_INT16U                                          // 7
	IEC61850_INT24U                                          // 8
	IEC61850_INT32U                                          // 9
	IEC61850_FLOAT32                                         // 10
	IEC61850_FLOAT64                                         // 11
	IEC61850_ENUMERATED                                      // 12
	IEC61850_OCTET_STRING_64                                 // 13
	IEC61850_OCTET_STRING_6                                  // 14
	IEC61850_OCTET_STRING_8                                  // 15
	IEC61850_VISIBLE_STRING_32                               // 16
	IEC61850_VISIBLE_STRING_64                               // 17
	IEC61850_VISIBLE_STRING_65                               // 18
	IEC61850_VISIBLE_STRING_129                              // 19
	IEC61850_VISIBLE_STRING_255                              // 20
	IEC61850_UNICODE_STRING_255                              // 21
	IEC61850_TIMESTAMP                                       // 22
	IEC61850_QUALITY                                         // 23
	IEC61850_CHECK                                           // 24
	IEC61850_CODEDENUM                                       // 25
	IEC61850_GENERIC_BITSTRING                               // 26
	IEC61850_CONSTRUCTED                                     // 27
	IEC61850_ENTRY_TIME                                      // 28
	IEC61850_PHYCOMADDR                                      // 29
	IEC61850_CURRENCY                                        // 30
	IEC61850_OPTFLDS                                         // 31
	IEC61850_TRGOPS                                          // 32
)
