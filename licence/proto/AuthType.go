// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package proto

import "strconv"

type AuthType byte

const (
	AuthTypeTimeAuth   AuthType = 1
	AuthTypeLocalAuth  AuthType = 2
	AuthTypeCentreAuth AuthType = 4
)

var EnumNamesAuthType = map[AuthType]string{
	AuthTypeTimeAuth:   "TimeAuth",
	AuthTypeLocalAuth:  "LocalAuth",
	AuthTypeCentreAuth: "CentreAuth",
}

var EnumValuesAuthType = map[string]AuthType{
	"TimeAuth":   AuthTypeTimeAuth,
	"LocalAuth":  AuthTypeLocalAuth,
	"CentreAuth": AuthTypeCentreAuth,
}

func (v AuthType) String() string {
	if s, ok := EnumNamesAuthType[v]; ok {
		return s
	}
	return "AuthType(" + strconv.FormatInt(int64(v), 10) + ")"
}
