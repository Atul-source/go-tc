package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaHtbUnspec = iota
	tcaHtbParms
	tcaHtbInit
	tcaHtbCtab
	tcaHtbRtab
	tcaHtbDirectQlen
	tcaHtbRate64
	tcaHtbCeil64
	tcaHtbPad
)

// Htb contains attributes of the HTB discipline
type Htb struct {
	Parms      *HtbOpt
	Init       *HtbGlob
	Ctab       []byte
	Rtab       []byte
	DirectQlen uint32
	Rate64     uint64
	Ceil64     uint64
}

// unmarshalHtb parses the Htb-encoded data and stores the result in the value pointed to by info.
func unmarshalHtb(data []byte, info *Htb) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaHtbParms:
			opt := &HtbOpt{}
			if err := unmarshalStruct(ad.Bytes(), opt); err != nil {
				return err
			}
			info.Parms = opt
		case tcaHtbInit:
			glob := &HtbGlob{}
			if err := unmarshalStruct(ad.Bytes(), glob); err != nil {
				return err
			}
			info.Init = glob
		case tcaHtbCtab:
			info.Ctab = ad.Bytes()
		case tcaHtbRtab:
			info.Rtab = ad.Bytes()
		case tcaHtbDirectQlen:
			info.DirectQlen = ad.Uint32()
		case tcaHtbRate64:
			info.Rate64 = ad.Uint64()
		case tcaHtbCeil64:
			info.Ceil64 = ad.Uint64()
		case tcaHtbPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalHtb()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalHtb returns the binary encoding of Qfq
func marshalHtb(info *Htb) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Htb options are missing")
	}
	// TODO: improve logic and check combinations
	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaHtbParms, Data: data})
	}
	if info.Init != nil {
		data, err := marshalStruct(info.Init)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaHtbInit, Data: data})
	}
	if info.DirectQlen != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaHtbDirectQlen, Data: info.DirectQlen})
	}
	if info.Rate64 != 0 {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaHtbRate64, Data: info.Rate64})
	}
	if info.Ceil64 != 0 {
		options = append(options, tcOption{Interpretation: vtUint64, Type: tcaHtbCeil64, Data: info.Ceil64})
	}
	return marshalAttributes(options)
}

// HtbGlob from include/uapi/linux/pkt_sched.h
type HtbGlob struct {
	Version      uint32
	Rate2Quantum uint32
	Defcls       uint32
	Debug        uint32
	DirectPkts   uint32
}

// HtbOpt from include/uapi/linux/pkt_sched.h
type HtbOpt struct {
	Rate    RateSpec
	Ceil    RateSpec
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
}
