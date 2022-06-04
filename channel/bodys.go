package channel

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/fields"
)

const (
	ChannelTransferDirectionHacashLeftToRight  uint8 = 1
	ChannelTransferDirectionHacashRightToLeft  uint8 = 2
	ChannelTransferDirectionSatoshiLeftToRight uint8 = 3
	ChannelTransferDirectionSatoshiRightToLeft uint8 = 4
)

// Channel transfer, data body
type ChannelChainTransferProveBodyInfo struct {
	ChannelId fields.ChannelId // Channel ID

	ReuseVersion   fields.VarUint4 // Channel reuse sequence number
	BillAutoNumber fields.VarUint8 // Serial number of channel bill

	PayDirection fields.VarUint1         // Capital flow direction: HAC 1 Left = > right; 2. right = > left BTC 3 Left = > right; 4. right = > left
	PayAmount    fields.Amount           // Payment amount cannot be negative
	PaySatoshi   fields.SatoshiVariation // Pay bitcoin sat amount

	LeftBalance  fields.Amount // Real time amount on the left
	RightBalance fields.Amount // Right real time amount

	LeftSatoshi  fields.SatoshiVariation // Number of bitcoin sat on the left
	RightSatoshi fields.SatoshiVariation // Number of bitcoin sat on the right

	LeftAddress  fields.Address // Left address
	RightAddress fields.Address // Right address
}

func CreateEmptyProveBody(cid fields.ChannelId) *ChannelChainTransferProveBodyInfo {
	emptyamt1 := fields.NewEmptyAmount()
	emptyamt2 := fields.NewEmptyAmount()
	emptyamt3 := fields.NewEmptyAmount()
	sat1 := fields.NewEmptySatoshiVariation()
	sat2 := fields.NewEmptySatoshiVariation()
	sat3 := fields.NewEmptySatoshiVariation()
	return &ChannelChainTransferProveBodyInfo{
		ChannelId:      cid,
		ReuseVersion:   1,
		BillAutoNumber: 0,
		PayDirection:   1,
		PayAmount:      *emptyamt3,
		PaySatoshi:     sat3,
		LeftBalance:    *emptyamt1,
		RightBalance:   *emptyamt2,
		LeftSatoshi:    sat1,
		RightSatoshi:   sat2,
		LeftAddress:    nil,
		RightAddress:   nil,
	}
}

// interface
func (e *ChannelChainTransferProveBodyInfo) GetChannelId() fields.ChannelId {
	return e.ChannelId
}
func (e *ChannelChainTransferProveBodyInfo) GetLeftBalance() fields.Amount {
	return e.LeftBalance
}
func (e *ChannelChainTransferProveBodyInfo) GetRightBalance() fields.Amount {
	return e.RightBalance
}
func (e *ChannelChainTransferProveBodyInfo) GetLeftSatoshi() fields.Satoshi {
	return e.LeftSatoshi.GetRealSatoshi()
}
func (e *ChannelChainTransferProveBodyInfo) GetRightSatoshi() fields.Satoshi {
	return e.RightSatoshi.GetRealSatoshi()
}
func (e *ChannelChainTransferProveBodyInfo) GetLeftAddress() fields.Address {
	return e.LeftAddress
}
func (e *ChannelChainTransferProveBodyInfo) GetRightAddress() fields.Address {
	return e.RightAddress
}
func (e *ChannelChainTransferProveBodyInfo) GetReuseVersion() uint32 {
	return uint32(e.ReuseVersion)
}
func (e *ChannelChainTransferProveBodyInfo) GetAutoNumber() uint64 {
	return uint64(e.BillAutoNumber)
}

func (elm *ChannelChainTransferProveBodyInfo) Size() uint32 {
	size := elm.ChannelId.Size() +
		elm.ReuseVersion.Size() +
		elm.BillAutoNumber.Size() +
		elm.PayDirection.Size() +
		elm.PayAmount.Size() +
		elm.PaySatoshi.Size() +
		elm.LeftBalance.Size() +
		elm.RightBalance.Size() +
		elm.LeftSatoshi.Size() +
		elm.RightSatoshi.Size() +
		elm.LeftAddress.Size() +
		elm.RightAddress.Size()
	// ok
	return size
}

func (elm *ChannelChainTransferProveBodyInfo) Serialize() ([]byte, error) {
	var bt []byte
	var buffer bytes.Buffer
	bt, _ = elm.ChannelId.Serialize()
	buffer.Write(bt)
	bt, _ = elm.ReuseVersion.Serialize()
	buffer.Write(bt)
	bt, _ = elm.BillAutoNumber.Serialize()
	buffer.Write(bt)
	bt, _ = elm.PayDirection.Serialize()
	buffer.Write(bt)
	bt, _ = elm.PayAmount.Serialize()
	buffer.Write(bt)
	bt, _ = elm.PaySatoshi.Serialize()
	buffer.Write(bt)
	bt, _ = elm.LeftBalance.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightBalance.Serialize()
	buffer.Write(bt)
	bt, _ = elm.LeftSatoshi.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightSatoshi.Serialize()
	buffer.Write(bt)
	bt, _ = elm.LeftAddress.Serialize()
	buffer.Write(bt)
	bt, _ = elm.RightAddress.Serialize()
	buffer.Write(bt)
	return buffer.Bytes(), nil
}

func (elm *ChannelChainTransferProveBodyInfo) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ReuseVersion.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.BillAutoNumber.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.PayDirection.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.PayAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.PaySatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftBalance.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightBalance.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightSatoshi.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.LeftAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.RightAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *ChannelChainTransferProveBodyInfo) GetSignStuff() []byte {
	var conbt, _ = elm.Serialize() // Data body
	return conbt                   // Hash
}
func (elm *ChannelChainTransferProveBodyInfo) GetSignStuffHashHalfChecker() fields.HashHalfChecker {
	var conbt = elm.GetSignStuff()                      // Data body
	return fields.CalculateHash(conbt).GetHalfChecker() // Hash detection
}

// Check signature
func (elm *ChannelChainTransferProveBodyInfo) CheckAddressAndSign(leftAddress, rightAddress fields.Address) error {
	// All checked successfully
	return nil
}

////////////////////////////////

type ChannelPayProveBodyList struct {
	Count      fields.VarUint1
	ProveBodys []*ChannelChainTransferProveBodyInfo
}

func (c ChannelPayProveBodyList) Size() uint32 {
	size := c.Count.Size()
	for i := 0; i < int(c.Count); i++ {
		size += c.ProveBodys[i].Size()
	}
	// ok
	return size
}

func (c ChannelPayProveBodyList) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = c.Count.Serialize() // Data body
	buffer.Write(bt1)
	for i := 0; i < len(c.ProveBodys); i++ {
		var bt6, _ = c.ProveBodys[i].Serialize()
		buffer.Write(bt6)
	}
	return buffer.Bytes(), nil
}

func (c *ChannelPayProveBodyList) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	// passageway
	seek, e = c.Count.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	ccn := int(c.Count)
	c.ProveBodys = make([]*ChannelChainTransferProveBodyInfo, ccn)
	for i := 0; i < ccn; i++ {
		c.ProveBodys[i] = &ChannelChainTransferProveBodyInfo{}
		seek, e = c.ProveBodys[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// complete
	return seek, nil
}

////////////////////////////////

// Channel chain payment

// Channel chain transfer transactions (signed off the chain) can be arbitrated on the chain
// Adopt zero knowledge proof mode
type OffChainFormPaymentChannelTransfer struct {
	Timestamp                fields.BlockTxTimestamp // time stamp
	OrderNoteHashHalfChecker fields.HashHalfChecker  // Order detail data hash len = 16

	MustSignCount     fields.VarUint1  // Number of addresses that must be signed, max. 200
	MustSignAddresses []fields.Address // 顺序打乱/随机的通道必须签名的地址

	ChannelCount                         fields.VarUint1          // Number of access channels, max. 200
	ChannelTransferProveHashHalfCheckers []fields.HashHalfChecker // Hash of channel transfer certificate, from payment to final collection, hash len = 16

	MustSigns []fields.Sign // 顺序打乱/随机的签名，顺序与地址相同
}

func (elm *OffChainFormPaymentChannelTransfer) Size() uint32 {
	size := elm.Timestamp.Size() +
		elm.OrderNoteHashHalfChecker.Size() +
		elm.ChannelCount.Size() +
		elm.MustSignCount.Size()

	size += uint32(len(elm.ChannelTransferProveHashHalfCheckers)) * (fields.HashHalfCheckerSize)
	size += uint32(len(elm.MustSignAddresses)) * (fields.AddressSize)
	size += uint32(len(elm.MustSigns)) * fields.SignSize
	return size
}

func (elm *OffChainFormPaymentChannelTransfer) SerializeForPrefixSignStuff() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.Timestamp.Serialize()
	var bt2, _ = elm.OrderNoteHashHalfChecker.Serialize()
	buffer.Write(bt1)
	buffer.Write(bt2)
	var bt3, _ = elm.MustSignCount.Serialize()
	buffer.Write(bt3)
	for i := 0; i < len(elm.MustSignAddresses); i++ {
		var bt4, _ = elm.MustSignAddresses[i].Serialize()
		buffer.Write(bt4)
	}
	var bt4, _ = elm.ChannelCount.Serialize()
	buffer.Write(bt4)
	return buffer.Bytes(), nil

}

func (elm *OffChainFormPaymentChannelTransfer) SerializeNoSign() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.SerializeForPrefixSignStuff() // Data body
	buffer.Write(bt1)
	for i := 0; i < len(elm.ChannelTransferProveHashHalfCheckers); i++ {
		var bt4, _ = elm.ChannelTransferProveHashHalfCheckers[i].Serialize()
		buffer.Write(bt4)
	}
	return buffer.Bytes(), nil
}

func (elm *OffChainFormPaymentChannelTransfer) GetSignStuffHash() fields.Hash {
	var conbt, _ = elm.SerializeNoSign() // Data body
	return fields.CalculateHash(conbt)   // Hash
}

func (elm *OffChainFormPaymentChannelTransfer) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.SerializeNoSign() // Data body
	buffer.Write(bt1)
	for i := 0; i < len(elm.MustSigns); i++ {
		var bt6, _ = elm.MustSigns[i].Serialize()
		buffer.Write(bt6)
	}
	return buffer.Bytes(), nil
}

func (elm *OffChainFormPaymentChannelTransfer) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.Timestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.OrderNoteHashHalfChecker.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// address
	seek, e = elm.MustSignCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	scn := int(elm.MustSignCount)
	elm.MustSignAddresses = make([]fields.Address, scn)
	for i := 0; i < scn; i++ {
		elm.MustSignAddresses[i] = fields.Address{}
		seek, e = elm.MustSignAddresses[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// passageway
	seek, e = elm.ChannelCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	ccn := int(elm.ChannelCount)
	elm.ChannelTransferProveHashHalfCheckers = make([]fields.HashHalfChecker, ccn)
	for i := 0; i < ccn; i++ {
		elm.ChannelTransferProveHashHalfCheckers[i] = fields.HashHalfChecker{}
		seek, e = elm.ChannelTransferProveHashHalfCheckers[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// autograph
	elm.MustSigns = make([]fields.Sign, scn)
	for i := 0; i < scn; i++ {
		elm.MustSigns[i] = fields.CreateEmptySign()
		seek, e = elm.MustSigns[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// complete
	return seek, nil
}

// Check data availability
func (elm *OffChainFormPaymentChannelTransfer) CheckValidity() error {
	return nil
}

// Verify signature on ticket
func (elm *OffChainFormPaymentChannelTransfer) VerifySignature() error {
	return elm.CheckMustAddressAndSigns()
}

// Populate signatures by location
func (elm *OffChainFormPaymentChannelTransfer) FillSignByPosition(sign fields.Sign) error {
	sgaddr := sign.GetAddress()
	sn := int(elm.MustSignCount)
	var istok = false
	for i := 0; i < sn; i++ {
		addr := elm.MustSignAddresses[i]
		if addr.Equal(sgaddr) {
			istok = true
			elm.MustSigns[i] = sign
		}
	}
	if istok == false {
		return fmt.Errorf(" sign address %s not find in must list.", sgaddr.ToReadable())
	}
	return nil
}

// Sign and fill to the specified location
func (elm *OffChainFormPaymentChannelTransfer) DoSignFillPosition(acc *account.Account) (*fields.Sign, error) {
	// Compute hash
	hash := elm.GetSignStuffHash()
	// autograph
	signature, e2 := acc.Private.Sign(hash)
	if e2 != nil {
		return nil, fmt.Errorf("Private Key '" + fields.Address(acc.Address).ToReadable() + "' do sign error")
	}
	sigObj := fields.Sign{
		PublicKey: acc.PublicKey,
		Signature: signature.Serialize64(),
	}
	// Fill to specified position
	sgaddr := sigObj.GetAddress()
	sn := int(elm.MustSignCount)
	var istok = false
	for i := 0; i < sn; i++ {
		addr := elm.MustSignAddresses[i]
		if addr.Equal(sgaddr) {
			istok = true
			elm.MustSigns[i] = sigObj
		}
	}
	if istok == false {
		return nil, fmt.Errorf(" sign address %s not find in must list.", sgaddr.ToReadable())
	}
	return &sigObj, nil
}

// Check whether an address is signed
func (elm *OffChainFormPaymentChannelTransfer) CheckOneAddressSign(addr fields.Address) error {
	hx := elm.GetSignStuffHash()
	for _, v := range elm.MustSigns {
		if v.GetAddress().Equal(addr) {
			ok, _ := account.CheckSignByHash32(hx, v.PublicKey, v.Signature)
			if !ok {
				return fmt.Errorf("address %s verify signature fail.", addr.ToReadable())
			} else {
				return nil // Validation successful
			}
		}
	}
	return fmt.Errorf("address %s signature not find.", addr.ToReadable())
}

// Check all signatures
func (elm *OffChainFormPaymentChannelTransfer) CheckMustAddressAndSigns() error {
	var e error

	stuff, e := elm.SerializeNoSign()
	if e != nil {
		return e
	}
	conhx := fields.CalculateHash(stuff)

	// Inspection quantity
	sn := int(elm.MustSignCount)
	if sn < 2 || sn > 200 {
		return fmt.Errorf("MustSignCount error.")
	}
	if sn != len(elm.MustSignAddresses) || sn != len(elm.MustSigns) {
		return fmt.Errorf("Addresses or Signs length error.")
	}

	// Signatures are arranged by address. Check whether all addresses and signatures match
	for i := 0; i < sn; i++ {
		sign := elm.MustSigns[i]
		addr := elm.MustSignAddresses[i]
		sgaddr := account.NewAddressFromPublicKeyV0(sign.PublicKey)
		// Judge address order
		if addr.NotEqual(sgaddr) {
			return fmt.Errorf("Address not match, need %s nut got %s.",
				addr.ToReadable(), fields.Address(sgaddr).ToReadable())
		}
		// Check signature
		ok, _ := account.CheckSignByHash32(conhx, sign.PublicKey, sign.Signature)
		if !ok {
			return fmt.Errorf("account %s verify signature fail.", addr.ToReadable())
		}
	}

	// All signatures verified successfully
	return nil
}
