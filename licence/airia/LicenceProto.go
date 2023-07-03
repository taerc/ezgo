// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package airia

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type LicenceProto struct {
	_tab flatbuffers.Table
}

func GetRootAsLicenceProto(buf []byte, offset flatbuffers.UOffsetT) *LicenceProto {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &LicenceProto{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *LicenceProto) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *LicenceProto) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *LicenceProto) Version() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LicenceProto) MutateVersion(n byte) bool {
	return rcv._tab.MutateByteSlot(4, n)
}

func (rcv *LicenceProto) MagicValue() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *LicenceProto) MagicSignature() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *LicenceProto) AuthType() AuthType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return AuthType(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 1
}

func (rcv *LicenceProto) MutateAuthType(n AuthType) bool {
	return rcv._tab.MutateByteSlot(10, byte(n))
}

func (rcv *LicenceProto) DeviceDesc() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *LicenceProto) TimeInfo(obj *TimeInfo) *TimeInfo {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(TimeInfo)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *LicenceProto) LocalInfo(obj *LocalInfo) *LocalInfo {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(LocalInfo)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *LicenceProto) CentreInfo(obj *CentreInfo) *CentreInfo {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(CentreInfo)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func LicenceProtoStart(builder *flatbuffers.Builder) {
	builder.StartObject(8)
}
func LicenceProtoAddVersion(builder *flatbuffers.Builder, version byte) {
	builder.PrependByteSlot(0, version, 0)
}
func LicenceProtoAddMagicValue(builder *flatbuffers.Builder, magicValue flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(magicValue), 0)
}
func LicenceProtoAddMagicSignature(builder *flatbuffers.Builder, magicSignature flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(magicSignature), 0)
}
func LicenceProtoAddAuthType(builder *flatbuffers.Builder, authType AuthType) {
	builder.PrependByteSlot(3, byte(authType), 1)
}
func LicenceProtoAddDeviceDesc(builder *flatbuffers.Builder, deviceDesc flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(deviceDesc), 0)
}
func LicenceProtoAddTimeInfo(builder *flatbuffers.Builder, timeInfo flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(timeInfo), 0)
}
func LicenceProtoAddLocalInfo(builder *flatbuffers.Builder, localInfo flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(localInfo), 0)
}
func LicenceProtoAddCentreInfo(builder *flatbuffers.Builder, centreInfo flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(centreInfo), 0)
}
func LicenceProtoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}