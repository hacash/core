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

// 通道转账，数据体
type ChannelChainTransferProveBodyInfo struct {
	ChannelId fields.ChannelId // 通道id

	ReuseVersion   fields.VarUint4 // 通道重用序号
	BillAutoNumber fields.VarUint8 // 通道账单流水序号

	PayDirection fields.VarUint1         // 资金流动方向： HAC 1.左=>右； 2.右=>左  BTC 3.左=>右；4.右=>左
	PayAmount    fields.Amount           // 支付金额，不能为负
	PaySatoshi   fields.SatoshiVariation // 支付比特币sat金额

	LeftBalance  fields.Amount // 左侧实时金额
	RightBalance fields.Amount // 右侧实时金额

	LeftSatoshi  fields.SatoshiVariation // 左侧比特币sat数量
	RightSatoshi fields.SatoshiVariation // 右侧比特币sat数量

	LeftAddress  fields.Address // 左侧地址
	RightAddress fields.Address // 右侧地址
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
	var conbt, _ = elm.Serialize() // 数据体
	return conbt                   // 哈希
}
func (elm *ChannelChainTransferProveBodyInfo) GetSignStuffHashHalfChecker() fields.HashHalfChecker {
	var conbt = elm.GetSignStuff()                      // 数据体
	return fields.CalculateHash(conbt).GetHalfChecker() // 哈希检测
}

// 检查签名
func (elm *ChannelChainTransferProveBodyInfo) CheckAddressAndSign(leftAddress, rightAddress fields.Address) error {
	// 全部检查成功
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
	var bt1, _ = c.Count.Serialize() // 数据体
	buffer.Write(bt1)
	for i := 0; i < len(c.ProveBodys); i++ {
		var bt6, _ = c.ProveBodys[i].Serialize()
		buffer.Write(bt6)
	}
	return buffer.Bytes(), nil
}

func (c *ChannelPayProveBodyList) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	// 通道
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
	// 完成
	return seek, nil
}

////////////////////////////////

// 通道链支付

// 通道链转账交易（链下签署）可以上链仲裁
// 采用零知识证明模式
type OffChainFormPaymentChannelTransfer struct {
	Timestamp                fields.BlockTxTimestamp // 时间戳
	OrderNoteHashHalfChecker fields.HashHalfChecker  // 订单详情数据哈希  len = 16

	MustSignCount     fields.VarUint1  // 必须签名地址的数量，最大值 200
	MustSignAddresses []fields.Address // 顺序打乱/随机的通道必须签名的地址

	ChannelCount                         fields.VarUint1          // 途径通道数量，最大值 200
	ChannelTransferProveHashHalfCheckers []fields.HashHalfChecker // 通道转账证明哈希，顺序为从付款到最后收款，哈希  len = 16

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
	var bt1, _ = elm.SerializeForPrefixSignStuff() // 数据体
	buffer.Write(bt1)
	for i := 0; i < len(elm.ChannelTransferProveHashHalfCheckers); i++ {
		var bt4, _ = elm.ChannelTransferProveHashHalfCheckers[i].Serialize()
		buffer.Write(bt4)
	}
	return buffer.Bytes(), nil
}

func (elm *OffChainFormPaymentChannelTransfer) GetSignStuffHash() fields.Hash {
	var conbt, _ = elm.SerializeNoSign() // 数据体
	return fields.CalculateHash(conbt)   // 哈希
}

func (elm *OffChainFormPaymentChannelTransfer) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	var bt1, _ = elm.SerializeNoSign() // 数据体
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
	// 地址
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
	// 通道
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
	// 签名
	elm.MustSigns = make([]fields.Sign, scn)
	for i := 0; i < scn; i++ {
		elm.MustSigns[i] = fields.CreateEmptySign()
		seek, e = elm.MustSigns[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// 完成
	return seek, nil
}

// 检查数据可用性
func (elm *OffChainFormPaymentChannelTransfer) CheckValidity() error {
	return nil
}

// 验证对票据的签名
func (elm *OffChainFormPaymentChannelTransfer) VerifySignature() error {
	return elm.CheckMustAddressAndSigns()
}

// 按位置填充签名
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

// 签名并填充至指定位置
func (elm *OffChainFormPaymentChannelTransfer) DoSignFillPosition(acc *account.Account) (*fields.Sign, error) {
	// 计算哈希
	hash := elm.GetSignStuffHash()
	// 签名
	signature, e2 := acc.Private.Sign(hash)
	if e2 != nil {
		return nil, fmt.Errorf("Private Key '" + fields.Address(acc.Address).ToReadable() + "' do sign error")
	}
	sigObj := fields.Sign{
		PublicKey: acc.PublicKey,
		Signature: signature.Serialize64(),
	}
	// 填充到指定位置
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

// 检查某一个地址是否签名
func (elm *OffChainFormPaymentChannelTransfer) CheckOneAddressSign(addr fields.Address) error {
	hx := elm.GetSignStuffHash()
	for _, v := range elm.MustSigns {
		if v.GetAddress().Equal(addr) {
			ok, _ := account.CheckSignByHash32(hx, v.PublicKey, v.Signature)
			if !ok {
				return fmt.Errorf("address %s verify signature fail.", addr.ToReadable())
			} else {
				return nil // 验证成功
			}
		}
	}
	return fmt.Errorf("address %s signature not find.", addr.ToReadable())
}

// 检查所有签名
func (elm *OffChainFormPaymentChannelTransfer) CheckMustAddressAndSigns() error {
	var e error

	stuff, e := elm.SerializeNoSign()
	if e != nil {
		return e
	}
	conhx := fields.CalculateHash(stuff)

	// 检查数量
	sn := int(elm.MustSignCount)
	if sn < 2 || sn > 200 {
		return fmt.Errorf("MustSignCount error.")
	}
	if sn != len(elm.MustSignAddresses) || sn != len(elm.MustSigns) {
		return fmt.Errorf("Addresses or Signs length error.")
	}

	// 签名按地址排列，检查所有地址和签名是否匹配
	for i := 0; i < sn; i++ {
		sign := elm.MustSigns[i]
		addr := elm.MustSignAddresses[i]
		sgaddr := account.NewAddressFromPublicKeyV0(sign.PublicKey)
		// 判断地址顺序
		if addr.NotEqual(sgaddr) {
			return fmt.Errorf("Address not match, need %s nut got %s.",
				addr.ToReadable(), fields.Address(sgaddr).ToReadable())
		}
		// 检查签名
		ok, _ := account.CheckSignByHash32(conhx, sign.PublicKey, sign.Signature)
		if !ok {
			return fmt.Errorf("account %s verify signature fail.", addr.ToReadable())
		}
	}

	// 全部签名验证成功
	return nil
}
