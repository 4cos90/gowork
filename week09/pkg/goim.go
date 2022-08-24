package pkg

import (
	"encoding/binary"
	"errors"
)

var (
	packLengthFieldSize    = 4
	headerLengthFieldSize  = 2
	protocVersionFieldSize = 2
	operationCodeFieldSize = 4
	sequenceIdFieldSize    = 4
	ErrPackInComplete      = errors.New("error: package is incomplete")
)

// define package
type Pack struct {
	Length          int
	HeaderLength    int
	ProtocolVersion int
	OperationCode   int
	Seq             int
	Content         []byte
}

func NewPack(version, code, seq int, content []byte) *Pack {
	headerSize := headerSize()
	return &Pack{
		Length:          len(content) + headerSize,
		HeaderLength:    headerSize,
		ProtocolVersion: version,
		OperationCode:   code,
		Seq:             seq,
		Content:         content,
	}
}

// encode requset message
func Encoder(pack *Pack) []byte {
	res := make([]byte, pack.Length)

	// set package length
	binary.BigEndian.PutUint32(
		res[:headerLengthStart()],
		uint32(pack.Length),
	)
	// set header length
	binary.BigEndian.PutUint16(
		res[headerLengthStart():protocolVersionStart()],
		uint16(headerSize()),
	)
	// set protocol version
	binary.BigEndian.PutUint16(
		res[protocolVersionStart():operationCodeStart()],
		uint16(pack.ProtocolVersion),
	)
	// set operation code
	binary.BigEndian.PutUint32(
		res[operationCodeStart():sequenceIdStart()],
		uint32(pack.OperationCode),
	)
	// set sequence id
	binary.BigEndian.PutUint32(
		res[sequenceIdStart():sequenceIdStart()+sequenceIdFieldSize],
		uint32(pack.Seq),
	)
	// set body
	copy(res[headerSize():], pack.Content)

	return res
}

// decode request message
func Decoder(msg []byte) (*Pack, error) {
	if len(msg) < headerSize()+1 {
		return nil, ErrPackInComplete
	}
	// get package length
	packageLength := binary.BigEndian.Uint32(msg[:headerLengthStart()])
	// get header length
	headerLength := binary.BigEndian.Uint16(msg[headerLengthStart():protocolVersionStart()])
	// get protocol version
	protocolVersion := binary.BigEndian.Uint16(msg[protocolVersionStart():operationCodeStart()])
	// get operation code
	operationCode := binary.BigEndian.Uint32(msg[operationCodeStart():sequenceIdStart()])
	// get sequence id
	sequenceId := binary.BigEndian.Uint32(msg[sequenceIdStart() : sequenceIdStart()+sequenceIdFieldSize])
	// get data
	content := msg[headerSize():]
	return &Pack{
		Length:          int(packageLength),
		HeaderLength:    int(headerLength),
		ProtocolVersion: int(protocolVersion),
		OperationCode:   int(operationCode),
		Seq:             int(sequenceId),
		Content:         content,
	}, nil
}

// caculate header size
func headerSize() int {
	return packLengthFieldSize +
		headerLengthFieldSize +
		protocVersionFieldSize +
		operationCodeFieldSize +
		sequenceIdFieldSize
}

// position of header length field start
func headerLengthStart() int {
	return packLengthFieldSize
}

// position of protocol version field start
func protocolVersionStart() int {
	return headerLengthStart() + headerLengthFieldSize
}

// postion of operation code field start
func operationCodeStart() int {
	return protocolVersionStart() + protocVersionFieldSize
}

// position of sequence id field start
func sequenceIdStart() int {
	return operationCodeStart() + operationCodeFieldSize
}

func PackageLengthSize() int {
	return packLengthFieldSize
}
