package stores

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"strconv"
	"strings"
)

const SatoshiGenesisLogStorePageLimit int = 50

type SatoshiGenesis struct {
	TransferNo               fields.VarUint4 // 转账流水编号
	BitcoinBlockHeight       fields.VarUint4 // 转账的比特币区块高度
	BitcoinBlockTimestamp    fields.VarUint5 // 转账的比特币区块时间戳
	BitcoinEffectiveGenesis  fields.VarUint4 // 在这笔之前已经成功转移的比特币数量
	BitcoinQuantity          fields.VarUint4 // 本笔转账的比特币数量（单位：枚）
	AdditionalTotalHacAmount fields.VarUint4 // 本次转账[总共]应该增发的 hac 数量 （单位：枚）
	OriginAddress            fields.Address  // 转出的比特币来源地址
	BitcoinTransferHash      fields.Hash     // 比特币转账交易哈希
}

///////////////////////////////////////

func (s *SatoshiGenesis) Size() uint32 {
	return s.TransferNo.Size() +
		s.BitcoinBlockHeight.Size() +
		s.BitcoinBlockTimestamp.Size() +
		s.BitcoinEffectiveGenesis.Size() +
		s.BitcoinQuantity.Size() +
		s.AdditionalTotalHacAmount.Size() +
		s.OriginAddress.Size() +
		s.BitcoinTransferHash.Size()
}

func (s *SatoshiGenesis) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	b1, _ := s.TransferNo.Serialize()
	b2, _ := s.BitcoinBlockHeight.Serialize()
	b3, _ := s.BitcoinBlockTimestamp.Serialize()
	b4, _ := s.BitcoinEffectiveGenesis.Serialize()
	b5, _ := s.BitcoinQuantity.Serialize()
	b6, _ := s.AdditionalTotalHacAmount.Serialize()
	b7, _ := s.OriginAddress.Serialize()
	b8, _ := s.BitcoinTransferHash.Serialize()
	buffer.Write(b1)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(b5)
	buffer.Write(b6)
	buffer.Write(b7)
	buffer.Write(b8)
	return buffer.Bytes(), nil
}

func (s *SatoshiGenesis) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = s.TransferNo.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = s.BitcoinBlockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = s.BitcoinBlockTimestamp.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = s.BitcoinEffectiveGenesis.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = s.BitcoinQuantity.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = s.AdditionalTotalHacAmount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = s.OriginAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = s.BitcoinTransferHash.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

///////////////////////////////////////////////

func SatoshiGenesisPageSerialize(page []*SatoshiGenesis) []byte {
	var pagedata = bytes.NewBuffer([]byte{})
	for _, v := range page {
		bts, e := v.Serialize()
		if e != nil {
			return pagedata.Bytes()
		}
		pagedata.Write(bts)
	}
	return pagedata.Bytes()
}

func SatoshiGenesisPageSerializeForShow(page []*SatoshiGenesis) []string {
	var pagedata = make([]string, 0)
	for _, v := range page {
		pagedata = append(pagedata, fmt.Sprintf(
			"%d,%d,%d,%d,%d,%d,%s,%s",
			v.TransferNo,
			v.BitcoinBlockHeight,
			v.BitcoinBlockTimestamp,
			v.BitcoinEffectiveGenesis,
			v.BitcoinQuantity,
			v.AdditionalTotalHacAmount,
			v.OriginAddress.ToReadable(),
			v.BitcoinTransferHash.ToHex(),
		))
	}
	return pagedata
}

func SatoshiGenesisPageParse(buf []byte, seek uint32) []*SatoshiGenesis {
	var e error = nil
	var list = make([]*SatoshiGenesis, 0)
	for {
		var item = SatoshiGenesis{}
		seek, e = item.Parse(buf, seek)
		if e != nil {
			//fmt.Println(e)
			return list
		}
		list = append(list, &item)
	}
}

func SatoshiGenesisPageParseForShow(logitemstrlist []string) []*SatoshiGenesis {
	// 开始解析
	var allgenesis = make([]*SatoshiGenesis, 0)
	for _, logitemstr := range logitemstrlist {
		logitemstr = strings.Replace(logitemstr, " ", "", -1)
		dts := strings.Split(logitemstr, ",")
		if len(dts) != 8 {
			return nil
		}
		nums := make([]int64, 6)
		for i := 0; i < 6; i++ {
			n, e := strconv.ParseInt(dts[i], 10, 0)
			if e != nil {
				return nil
			}
			nums[i] = n
		}
		// 检查地址 和 txhx
		addr, ae := fields.CheckReadableAddress(dts[6])
		if ae != nil {
			return nil
		}
		trshx, te := hex.DecodeString(dts[7])
		if te != nil {
			return nil
		}
		if len(trshx) != 32 {
			return nil
		}
		// 生成
		genesis := SatoshiGenesis{
			TransferNo:               fields.VarUint4(nums[0]), // 转账流水编号
			BitcoinBlockHeight:       fields.VarUint4(nums[1]), // 转账的比特币区块高度
			BitcoinBlockTimestamp:    fields.VarUint5(nums[2]), // 转账的比特币区块时间戳
			BitcoinEffectiveGenesis:  fields.VarUint4(nums[3]), // 在这笔之前已经成功转移的比特币数量
			BitcoinQuantity:          fields.VarUint4(nums[4]), // 本笔转账的比特币数量（单位：枚）
			AdditionalTotalHacAmount: fields.VarUint4(nums[5]), // 本次转账[总共]应该增发的 hac 数量 （单位：枚）
			OriginAddress:            *addr,
			BitcoinTransferHash:      trshx,
		}
		allgenesis = append(allgenesis, &genesis)
	}
	// 返回
	return allgenesis
}
