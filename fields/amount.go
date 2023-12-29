package fields

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Amount struct {
	Unit    uint8
	Dist    int8
	Numeral []byte
}

func ParseAmount(buf []byte, seek uint32) *Amount {
	empty := NewEmptyAmount()
	empty.Parse(buf, seek)
	return empty
}

func NewEmptyAmountValue() Amount {
	return Amount{
		Unit:    0,
		Dist:    0,
		Numeral: []byte{},
	}
}

func NewEmptyAmount() *Amount {
	return &Amount{
		Unit:    0,
		Dist:    0,
		Numeral: []byte{},
	}
}

func NewAmountNumSmallCoin(num uint8) *Amount {
	return &Amount{
		Unit:    248,
		Dist:    1,
		Numeral: []byte{num},
	}
}

func NewAmountNumOneByUnit(unit uint8) *Amount {
	return &Amount{
		Unit:    unit,
		Dist:    1,
		Numeral: []byte{1},
	}
}

func NewAmountByUnit248(num int64) *Amount {
	amt, err := NewAmountByBigIntWithUnit(big.NewInt(num), 248)
	if err != nil {
		panic(err)
	}
	return amt
}

func NewAmountByBigIntWithUnit(bignum *big.Int, unit int) (*Amount, error) {
	var unitint = new(big.Int).Exp(big.NewInt(int64(10)), big.NewInt(int64(unit)), big.NewInt(int64(0)))
	//fmt.Println(bignum.String())
	//fmt.Println(unitint.String())
	return NewAmountByBigInt(bignum.Mul(bignum, unitint))
}

func NewAmountByBigInt(bignum *big.Int) (*Amount, error) {
	longnumstr := bignum.String()
	if longnumstr == "0" {
		return NewEmptyAmount(), nil
	}
	longnumstrary := []byte(longnumstr)
	strlen := len(longnumstrary)
	unit := 0
	for i := strlen - 1; i >= 0; i-- {
		if longnumstrary[i] == '0' {
			unit++
			if unit == 255 {
				break
			}
		} else {
			break
		}
	}
	numeralstr := string(longnumstrary[0 : strlen-unit])
	//fmt.Println("longnumstrary:", bignum.String())
	//fmt.Println("numeralstr:", numeralstr)
	numeralbigint, ok1 := new(big.Int).SetString(numeralstr, 10)
	if !ok1 {
		return nil, fmt.Errorf("Amount too big")
	}
	numeralbytes := numeralbigint.Bytes()
	dist := len(numeralbytes)
	if dist > 127 {
		return nil, fmt.Errorf("Amount too big")
	}
	new_amount := &Amount{
		Unit:    uint8(unit),
		Dist:    int8(dist),
		Numeral: numeralbytes,
	}
	if bignum.Sign() == -1 {
		new_amount.Dist *= -1 // amount is negative
	}
	return new_amount, nil
}

func NewAmount(unit uint8, num []byte) *Amount {
	dist := len(num)
	if dist > 127 {
		panic("Amount Numeral too long !")
	}
	return &Amount{
		Unit:    unit,
		Dist:    int8(dist),
		Numeral: num,
	}
}

func NewAmountSmallValue(num uint8, unit uint8) Amount {
	return Amount{
		Unit:    unit,
		Dist:    1,
		Numeral: []byte{num},
	}
}

func NewAmountSmall(num uint8, unit uint8) *Amount {
	v := NewAmountSmallValue(num, unit)
	return &v
}

func (bill Amount) Serialize() ([]byte, error) {
	var buffer = new(bytes.Buffer)
	buffer.Write([]byte{bill.Unit})
	buffer.Write([]byte{byte(bill.Dist)})
	buffer.Write(bill.Numeral)
	return buffer.Bytes(), nil
}

func (bill *Amount) Parse(buf []byte, seek uint32) (uint32, error) {
	if uint32(len(buf)) < seek+2 {
		return 0, fmt.Errorf("buf length not less than 2.")
	}
	bill.Unit = uint8(buf[seek])
	bill.Dist = int8(buf[seek+1])
	var numCount = int(bill.Dist)
	if bill.Dist < 0 {
		numCount *= -1
	}
	var tail = seek + 2 + uint32(numCount)
	if uint32(len(buf)) < tail {
		return 0, fmt.Errorf("buf length error.")
	}
	var nnnold = buf[seek+2 : tail]
	bill.Numeral = make([]byte, len(nnnold))
	copy(bill.Numeral, nnnold)
	return tail, nil
}

func (bill Amount) Size() uint32 {
	return 1 + 1 + uint32(len(bill.Numeral))
}

//////////////////////////////////////////////////////////

func (bill Amount) Copy() *Amount {
	num := make([]byte, len(bill.Numeral))
	copy(num, bill.Numeral)
	return &Amount{
		Unit:    bill.Unit,
		Dist:    bill.Dist,
		Numeral: num,
	}
}

func (bill Amount) GetValue() *big.Int {
	var bignum = new(big.Int)
	bignum.SetBytes(bill.Numeral)
	var sign = big.NewInt(int64(big.NewInt(int64(bill.Dist)).Sign()))
	var unit = new(big.Int).Exp(big.NewInt(int64(10)), big.NewInt(int64(bill.Unit)), big.NewInt(int64(0)))
	bignum.Mul(bignum, unit)
	bignum.Mul(bignum, sign) // do sign
	return bignum
}

func (bill Amount) IsEmpty() bool {
	return bill.Unit == 0 || bill.Dist == int8(0) || len(bill.Numeral) == 0
}
func (bill Amount) IsNotEmpty() bool {
	return bill.IsEmpty() == false
}

// Judgment must be positive and cannot be zero
func (bill Amount) IsPositive() bool {
	if bill.Unit == 0 {
		return false
	}
	if bill.Dist <= 0 {
		return false
	}
	// Meet requirements
	return true
}

// Judgment must be negative and cannot be zero
func (bill Amount) IsNegative() bool {
	if bill.Unit == 0 {
		return false
	}
	if bill.Dist >= 0 {
		return false
	}
	// Meet requirements
	return true
}

func (bill Amount) ToMeiOrFinString(usemei bool) string {
	if usemei {
		return bill.ToMeiString()
	} else {
		return bill.ToFinString()
	}
}

// Convert the unit to pieces, keep 8 decimal places, and round off the excess
func (bill Amount) ToUnitString(unit_name string) string {
	unit_name = strings.ToLower(unit_name)
	setunit := -1
	if unit_name == "mei" {
		setunit = 248
	}
	if unit_name == "zhu" {
		setunit = 240
	}
	if unit_name == "shuo" {
		setunit = 232
	}
	if unit_name == "ai" {
		setunit = 224
	}
	if unit_name == "miao" {
		setunit = 216
	}
	if setunit == -1 {
		// fin string
		return bill.ToFinString()
	}
	bigunit := bill.ToUnitBigFloat(setunit)
	meistr := bigunit.Text('f', 9)
	spx := strings.Split(meistr, ".")
	if len(spx) == 2 {
		if len(spx[1]) == 9 {
			spx[1] = strings.TrimRight(spx[1], spx[1][8:])
			meistr = strings.Join(spx, ".")
		}
	}
	return strings.TrimRight(strings.TrimRight(meistr, "0"), ".")

}

// Convert the unit to pieces, keep 8 decimal places, and round off the excess
func (bill Amount) ToMeiString() string {
	return bill.ToUnitString("mei")
}

// Conversion unit: piece
func (bill Amount) ToMeiBigFloat() *big.Float {
	return bill.ToUnitBigFloat(248)
}

func (bill Amount) ToZhuBigFloat() *big.Float {
	return bill.ToUnitBigFloat(240)
}

func (bill Amount) ToUnitBigFloat(unit int) *big.Float {
	// handle
	bigmei := new(big.Float).SetInt(new(big.Int).SetBytes(bill.Numeral))
	//fmt.Println(bigmei.String(), int(bill.Unit), int(bill.Unit) - 248)
	if bill.Dist < 0 {
		bigmei = bigmei.Neg(bigmei)
	}
	bigf10 := new(big.Float).SetFloat64(10.0)
	bigf10div := new(big.Float).SetFloat64(0.1)
	if bill.Unit > 0 {
		pz := int(bill.Unit) - unit
		if pz > 0 {
			for i := 0; i < pz; i++ {
				bigmei = bigmei.Mul(bigmei, bigf10)
			}
		} else if pz < 0 {
			pz = -pz
			for i := 0; i < pz; i++ {
				bigmei = bigmei.Mul(bigmei, bigf10div)
			}
		}
	}
	return bigmei
}

// Conversion unit: piece
func (bill Amount) ToMei() float64 {
	bigmei := bill.ToMeiBigFloat()
	mei, _ := bigmei.Float64()
	return mei
}

func (bill Amount) ToZhu() float64 {
	bigzhu := bill.ToZhuBigFloat()
	zhu, _ := bigzhu.Float64()
	return zhu
}

func (bill Amount) ToZhuOmit() uint64 {
	bigzhu := bill.ToZhuBigFloat()
	zhu, _ := bigzhu.Uint64()
	return zhu
}

// Create amount from string
func NewAmountFromString(numstr string) (*Amount, error) {
	if strings.Contains(numstr, ":") {
		return NewAmountFromFinString(numstr)
	} else {
		return NewAmountFromMeiStringUnsafe(numstr)
	}
}

// create form readble string
func NewAmountFromMeiStringUnsafe(meistr string) (*Amount, error) {
	mei, e1 := strconv.ParseFloat(meistr, 64)
	if e1 != nil {
		return nil, e1
	}
	return NewAmountFromMeiUnsafe(mei)
}

func NewAmountFromMeiUnsafe(mei float64) (*Amount, error) {
	mei1yi := uint64(mei * 10000 * 10000)
	return NewAmountByBigIntWithUnit(new(big.Int).SetUint64(mei1yi), 240)
}

func AmountToZeroFinString() string {
	return "ㄜ0:0"
}

// create form readble string
func NewAmountFromFinString(finstr string) (*Amount, error) {
	finstr = strings.ToUpper(finstr)
	finstr = strings.Replace(finstr, " ", "", -1) // Remove spaces
	finstr = strings.Replace(finstr, ",", "", -1) // Remove commas
	var sig = 1
	if strings.HasPrefix(finstr, "HCX") {
		finstr = string([]byte(finstr)[3:])
	} else if strings.HasPrefix(finstr, "HAC") {
		finstr = string([]byte(finstr)[3:])
	} else if strings.HasPrefix(finstr, "ㄜ") {
		finstr = string([]byte(finstr)[3:])
	}
	// negative
	if strings.HasPrefix(finstr, "-") {
		finstr = string([]byte(finstr)[1:])
		sig = -1 // negative
	}
	var main, dum, unit string
	var main_num, dum_num *big.Int
	var unit_num int
	var e error
	var ok bool
	part := strings.Split(finstr, ":")
	if len(part) != 2 {
		return nil, fmt.Errorf("format error")
	}
	unit = part[1]
	unit_num, e = strconv.Atoi(unit)
	if e != nil {
		return nil, fmt.Errorf("format error")
	}
	if unit_num < 0 || unit_num > 255 {
		return nil, fmt.Errorf("format error")
	}
	part2 := strings.Split(part[0], ":")
	if len(part2) < 1 || len(part2) > 2 {
		return nil, fmt.Errorf("format error")
	}

	main = part2[0]
	main_num, ok = new(big.Int).SetString(main, 10)
	if !ok {
		return nil, fmt.Errorf("format error")
	}
	if len(part2) == 2 {
		dum = part2[1]
		dum_num, ok = new(big.Int).SetString(dum, 10)
		if !ok {
			return nil, fmt.Errorf("format error")
		}
	}
	// Process decimal parts
	bigint0, _ := new(big.Int).SetString("0", 10)
	bigint1, _ := new(big.Int).SetString("1", 10)
	bigint10, _ := new(big.Int).SetString("10", 10)
	dum_wide10 := 0
	if dum_num != nil && dum_num.Cmp(bigint0) == 1 {
		mover := dum_num.Div(dum_num, bigint10).Add(dum_num, bigint1)
		dum_wide10 = int(mover.Int64())
		if unit_num-dum_wide10 < 0 {
			return nil, fmt.Errorf("format error")
		}
		main_num = main_num.Sub(main_num, mover)
		unit_num = unit_num - int(dum_wide10)
	}
	// negative
	if sig == -1 {
		main_num = main_num.Neg(main_num)
	}
	// transformation
	return NewAmountByBigIntWithUnit(main_num, unit_num)
}

func (bill Amount) ToFinString() string {
	return bill.ToFinStringWithMark("ㄜ")
}

func (bill Amount) ToFinStringWithMark(mark string) string {
	unitStr := strconv.Itoa(int(bill.Unit)) // string(bytes.Repeat([]byte{48}, int(bill.Unit)))
	numStr := new(big.Int).SetBytes(bill.Numeral).String()
	sig := ""
	if bill.Dist < 0 {
		sig = "-"
	}
	var numStrX string
	if len(numStr) > 3 {
		i := 1
		x := len(numStr) - 1
		for x >= 0 {
			numStrX = string([]byte(numStr)[x]) + numStrX
			if i%3 == 0 {
				numStrX = "," + numStrX
			}
			x--
			i++
		}
	} else {
		numStrX = numStr
	}
	numStrX = strings.TrimLeft(numStrX, ",")
	return mark + sig + numStrX + ":" + unitStr
}

func (bill Amount) ToFinStringWithMarkBySegmentSplit(mark string) string {
	// 248 256
	var unitnum = int(bill.Unit)
	unitStr := strconv.Itoa(unitnum) // string(bytes.Repeat([]byte{48}, int(bill.Unit)))
	numStr := new(big.Int).SetBytes(bill.Numeral).String()
	sig := ""
	if bill.Dist < 0 {
		sig = "-"
	}
	var numStrX = make([]string, 0)
	var nl = len(numStr)
	var mx int = unitnum + len(numStr)
	for i := 248 + 32; i >= 0 && i >= unitnum; i -= 8 {
		x1 := 0
		x2 := nl
		if i >= mx+8 {
			continue
		}
		if i > mx {
			x2 = 8 - (i - mx)
		} else {
			x1 = mx - i
			x2 = x1 + 8
		}
		if x2 > nl {
			x2 = nl
		}
		if x1 == x2 {
			continue
		}
		fmt.Println(i, unitnum, mx, ", ", x1, x2)
		numStrX = append(numStrX, numStr[x1:x2])
	}
	return mark + sig + strings.Join(numStrX, ",") + ":" + unitStr
}

// Omit the decimal part in order to save it in the 4-digit space
// Enlargeindicates whether to enlarge (round up)
func (bill *Amount) CompressForMainNumLen(numlen int, enlarge bool) (*Amount, bool, error) {
	var bignum = new(big.Int)
	bignum.SetBytes(bill.Numeral)
	if len(bignum.String()) <= numlen {
		return bill, false, nil // Data unchanged
	}
	var unitaddnum int = 0
	var bn1 = big.NewInt(int64(1))
	var bn10 = big.NewInt(int64(10))
	for {
		unitaddnum++
		bignum = bignum.Div(bignum, bn10)
		if enlarge {
			bignum = bignum.Add(bignum, bn1)
		}
		if len(bignum.String()) <= numlen {
			break
		}
	}
	newunit := unitaddnum + int(bill.Unit)
	if newunit > 255 {
		return nil, false, fmt.Errorf("bill.Unit too much long.")
	}
	newamt, e1 := NewAmountByBigIntWithUnit(bignum, newunit)
	if e1 != nil {
		return nil, false, nil
	}

	// success
	return newamt, true, nil
}

// Omit the decimal part in order to save it in the 11 bit space
func (bill *Amount) EllipsisDecimalFor11SizeStore() (*Amount, bool, error) {
	maxnumlen := 11 - 1 - 1
	if len(bill.Numeral) <= maxnumlen {
		return bill, false, nil // Data unchanged
	}
	// Omit decimal part
	new_unit := int(bill.Unit)
	new_dist := int(bill.Dist)
	new_numeral := append([]byte{}, bill.Numeral...)
	biglongnum := new(big.Int).SetBytes(bill.Numeral)
	bignum10 := big.NewInt(10)
	for {
		biglongnum = biglongnum.Div(biglongnum, bignum10)
		new_numeral = biglongnum.Bytes()
		new_unit += 1
		new_dist = len(new_numeral)
		if new_unit > 255 || new_dist > 127 {
			return nil, true, fmt.Errorf("Amount is too big.")
		}
		if new_dist <= maxnumlen {
			new_amount := &Amount{
				uint8(new_unit),
				int8(new_dist),
				new_numeral,
			}
			if bill.Dist < 0 {
				new_amount.Dist *= -1 // amount is negative
			}
			return new_amount, true, nil // amount is change
		}
		// continue next
	}
}

/*

// Omit the decimal part in order to save it in the 23 bit space
func (bill *Amount) EllipsisDecimalFor23SizeStore() (*Amount, bool) {
	maxnumlen := 23 - 1 - 1
	if len(bill.Numeral) <= maxnumlen {
		return bill, false // Data unchanged
	}
	// Omit decimal part
	longnumstr := new(big.Int).SetBytes(bill.Numeral).String()
	baselen := 0
	mvseek := len(longnumstr) / 2
	for true {
		longnumcut := string([]byte(longnumstr)[0 : baselen+mvseek])
		cutnum, ok := new(big.Int).SetString(longnumcut, 10)
		if !ok {
			panic("Amount to big !!!")
		}
		mvseek = mvseek / 2
		if mvseek == 0 {
			mvseek = 1 // Minimum movement
		}
		cutbytes := cutnum.Bytes()
		if len(cutbytes) == maxnumlen {
			sig := int8(1)
			if bill.Dist < 0 {
				sig = -1
			}
			unit := int(bill.Unit) + (len(longnumstr) - len(longnumcut))
			if unit > 255 {
				panic("Amount to big !!!")
			}
			return &Amount{
				uint8(unit),
				19 * sig,
				cutbytes,
			}, true // Changed
		} else if len(cutbytes) > maxnumlen {
			baselen -= mvseek

		} else if len(cutbytes) < maxnumlen {
			baselen += mvseek
		}
	}
	panic("Amount Ellipsis Decimal Error")
	return nil, false
}



*/

// add
func (bill Amount) Add(amt *Amount) (*Amount, error) {
	num1 := bill.GetValue()
	num2 := amt.GetValue()
	num1 = num1.Add(num1, num2)
	return NewAmountByBigInt(num1)
}

// sub
func (bill Amount) Sub(amt *Amount) (*Amount, error) {
	num1 := bill.GetValue()
	num2 := amt.GetValue()
	num1 = num1.Sub(num1, num2)
	return NewAmountByBigInt(num1)
}

// compare
func (bill Amount) LessThan(amt *Amount) bool {
	add1 := bill.GetValue()
	add2 := amt.GetValue()
	res := add1.Cmp(add2)
	if res == -1 {
		return true
	} else {
		return false
	}
}

// compare
func (bill Amount) MoreThan(amt *Amount) bool {
	add1 := bill.GetValue()
	add2 := amt.GetValue()
	res := add1.Cmp(add2)
	if res == 1 {
		return true
	} else {
		return false
	}
}

// compare
func (bill Amount) Equal(amt *Amount) bool {
	//
	if bill.Unit != amt.Unit ||
		bill.Dist != amt.Dist ||
		bytes.Compare(bill.Numeral, amt.Numeral) != 0 {
		return false
	}
	return true
}

func (bill Amount) NotEqual(amt *Amount) bool {
	return bill.Equal(amt) == false
}
