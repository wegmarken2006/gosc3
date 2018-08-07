package sc3

import (
	. "gosc3/osc"
	"testing"
)

func Test001(t *testing.T) {
	buf1 := EncodeI16(25)
	num := DecodeI16(buf1)
	if num != 25 {
		t.Errorf("Err enc/dec16 %d not ", num)
	}
	buf1 = EncodeI8(26)
	buf1 = EncodeI32(27)
	var ugens1 []UgenType
	ugens1 = append(ugens1, NewIConst(1))
	ugens1 = append(ugens1, NewFConst(3.3))

	p1 := NewPrimitive("P1", ugens1, []int{RateKr, RateIr})
	p1.Special = 0
	p1.Index = 0
	p2 := NewPrimitive("P2", []UgenType{}, []int{})
	p2.Rate = RateAr
	mc1 := NewMce([]UgenType{p1, p2})
	ugens2 := extend(p1.inputs, 5)
	if len(ugens2) != 5 {
		t.Errorf("Err extend %d not ", len(ugens2))
	}
	md := mceDegree(mc1)
	if md != 2 {
		t.Errorf("Err mceDegree %d not ", md)
	}
	mg1 := Mrg{left: mc1, right: p1}
	ex1 := mceExtend(3, mg1)
	if len(ex1) != 3 {
		t.Errorf("Err mceExtend %d not ", len(ex1))
	}
	uu1 := []UgenType{NewIConst(1), NewIConst(2)}
	uu2 := []UgenType{NewIConst(3), NewIConst(4)}
	uu3 := []UgenType{NewIConst(5), NewIConst(6)}
	uuu1 := make([][]UgenType, 3)
	uuu1[0] = uu1
	uuu1[1] = uu2
	uuu1[2] = uu3
	uuu2 := Transposer(uuu1)
	if len(uuu2) != 2 {
		t.Errorf("Err Transposer %d not ", len(uuu2))
	}

}
func Test002(t *testing.T) {
	var ugens1 []UgenType
	ugens1 = append(ugens1, NewIConst(1))
	ugens1 = append(ugens1, NewFConst(3.3))
	p1 := NewPrimitive("P1", ugens1, []int{RateKr, RateIr})
	p1.Special = 0
	p1.Index = 0
	p2 := NewPrimitive("P2", []UgenType{}, []int{})
	p2.Rate = RateAr

	mc1 := Mce{ugens: []UgenType{p1, p1}}
	mc2 := Mce{ugens: []UgenType{p1, p2}}
	mc3 := Mce{ugens: []UgenType{p1, p2, mc1}}

	p3 := Primitive{name: "P3", Rate: RateKr, inputs: []UgenType{mc1, mc3},
		outputs: []int{RateIr}, Special: 0, Index: 0}

	mc10 := mceTransform(p3)
	pp3 := mc10.(Mce).ugens[2]
	mg3 := Mrg{left: mc1, right: p2}
	switch pp3.(type) {
	case Primitive:
		if pp3.(Primitive).name != "P3" {
			t.Errorf("Err mceTransform")
		}
	default:
		t.Errorf("Err mceTransform")
	}

	//mc11 := mceExpand(mc1)
	l22 := mceChannels(mg3)
	el10 := l22[0]
	el11 := l22[1]
	switch el10.(type) {
	case Mrg:
	default:
		t.Errorf("Err mceChannels")
	}
	switch el11.(type) {
	case Primitive:
	default:
		t.Errorf("Err mceChannels")
	}
	iota1 := iota(5, 3, 1)
	if len(iota1) != 5 || iota1[4] != 7 || iota1[2] != 5 {
		t.Error(iota1)
		t.Errorf("Err iota")
	}
	prx1 := proxify(mc2)
	l23 := prx1.(Mce).ugens
	el12 := l23[0]
	el13 := l23[1]
	switch el12.(type) {
	case Mce:
	default:
		t.Errorf("Err Proxify")
	}
	switch el13.(type) {
	case Primitive:
	default:
		t.Errorf("Err Proxify")
	}

	ndk1 := NodeK{name: "ndk1", Def: 5, id: 30}
	ndk2 := NodeK{name: "ndk1", Def: 5, id: 31}
	ndc1 := NodeC{id: 20, value: 320}
	ndc2 := NodeC{id: 21, value: 321}
	ndu1 := NodeU{id: 40, name: "ndu1", Rate: RateDr, inputs: []UgenType{}, outputs: []int{}, Special: 11, UgenID: 2}
	ndu2 := NodeU{id: 41, name: "ndu2"}
	gr1 := Graph{nextID: 11, constants: []NodeC{ndc1, ndc2}, controls: []NodeK{ndk1, ndk2}, ugens: []NodeU{ndu1, ndu2}}
	mm1 := mkMap(gr1)
	lc1 := mm1.cs
	lk1 := mm1.ks
	lu1 := mm1.us
	if lc1[0] != 20 || lk1[1] != 31 || lu1[0] != 40 {
		t.Errorf("Err mkMap")
	}
	buf4 := encodeNodeK(mm1, ndk1)
	buf4Res := []byte{4, 0x6e, 0x64, 0x6b, 0x31, 0, 0, 0, 0, 0, 0}
	for ind, elem := range buf4 {
		if buf4Res[ind] != elem {
			t.Errorf("Err encodeNodeK")
		}
	}

	if !findCP(320, ndc1) {
		t.Errorf("Err findCP")
	}

	nn10, _ := mkNodeC(IConst{value: 320}, gr1)
	if nn10.(NodeC).id != 20 {
		t.Errorf("Err mkNodeC")
	}

	if !findKP("ndk1", ndk1) {
		t.Errorf("Err findKP")
	}

	ck1 := Control{name: "ndk1", Rate: RateKr, Index: 3}
	nn11, _ := mkNodeK(ck1, gr1)
	if nn11.(NodeK).id != 30 {
		t.Errorf("Err mkNodeK")
	}

}
