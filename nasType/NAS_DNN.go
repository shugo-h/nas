package nasType

import (
	"bytes"
	"fmt"
	"strings"
)

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
type DNN struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewDNN(iei uint8) (dNN *DNN) {
	dNN = &DNN{}
	dNN.SetIei(iei)
	return dNN
}

// DNN 9.11.2.1A
// Iei Row, sBit, len = [], 8, 8
func (a *DNN) GetIei() (iei uint8) {
	return a.Iei
}

// DNN 9.11.2.1A
// Iei Row, sBit, len = [], 8, 8
func (a *DNN) SetIei(iei uint8) {
	a.Iei = iei
}

// DNN 9.11.2.1A
// Len Row, sBit, len = [], 8, 8
func (a *DNN) GetLen() (len uint8) {
	return a.Len
}

// DNN 9.11.2.1A
// Len Row, sBit, len = [], 8, 8
func (a *DNN) SetLen(len uint8) {
	a.Len = len
	a.Buffer = make([]uint8, a.Len)
}

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
func (a *DNN) GetDNN() string {
	return a.decode(a.Buffer)
}

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
func (a *DNN) SetDNN(dNN string) {
	a.Buffer, _ = a.encode(dNN)
	a.Len = uint8(len(a.Buffer))
}

// Comply with 3GPP TS 23.003 clause 9.1
func (a *DNN) encode(fqdn string) ([]uint8, error) {
	if len(fqdn) >= 100 {
		return nil, fmt.Errorf("DNN string is too long")
	}

	var dnnNI string
	labels := strings.Split(fqdn, ".")
	if labels[len(labels)-1] == "gprs" {
		// containing APN-OI
		if len(labels) < 4 {
			return nil, fmt.Errorf("the number of DNN labes are too small (%s) ", fqdn)
		}
		dnnNI = strings.Join(labels[:len(labels)-3], ".")
	} else {
		// APN-NI only
		dnnNI = fqdn
	}

	if len(dnnNI) >= 63 {
		return nil, fmt.Errorf("APN-NI of the DNN is too long")
	}

	encodedDNN := make([]uint8, 0, 100)
	for _, label := range labels {
		encodedDNN = append(encodedDNN, uint8(len(label)))
		encodedDNN = append(encodedDNN, []uint8(label)...)
	}
	return encodedDNN, nil
}

// Comply with 3GPP TS 23.003 clause 9.1
func (a *DNN) decode(dNN []uint8) string {
	var fqdn string
	buffer := bytes.NewBuffer(dNN)
	for {
		if labelLen, err := buffer.ReadByte(); err != nil {
			break
		} else {
			fqdn += string(buffer.Next(int(labelLen))) + "."
		}
	}

	if len(fqdn) == 0 {
		return ""
	} else {
		return fqdn[:len(fqdn)-1]
	}
}
