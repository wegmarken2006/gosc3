package sc3

import (
	"fmt"
	. "gosc3/osc"
	"sort"
	"strconv"
)

const (
	RateIr = 0
	RateKr = 1
	RateAr = 2
	RateDr = 3
)

type UgenType interface {
	isUgen()
}

type NodeType interface {
	isNode()
}

type NodeTypeList []NodeType

// for sort
func (nn NodeTypeList) Len() int { return len(nn) }
func (s NodeTypeList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s NodeTypeList) Less(i, j int) bool {
	return i < j
}

type UgenList []UgenType

func (uu UgenList) Len() int { return len(uu) }
func (s UgenList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s UgenList) Less(i, j int) bool {
	return i < j
}

type Primitive struct {
	Rate int
	name string
	//inputs  []interface{}
	inputs  UgenList
	outputs []int
	Special int
	Index   int
}

// NewPrimitive primitive constructor
func NewPrimitive(name string, inputs UgenList, outputs []int) Primitive {
	pr := Primitive{}
	pr.Index = 0
	pr.Special = 0
	pr.Rate = RateKr
	pr.name = name
	pr.inputs = inputs
	pr.outputs = outputs
	return pr
}

type Proxy struct {
	primitive Primitive
	Index     int
}

func NewProxy(primitive Primitive) Proxy {
	px := Proxy{}
	px.Index = 0
	px.primitive = primitive
	return px
}

type Control struct {
	Rate  int
	name  string
	Index int
}

func NewControl(name string) Control {
	cc := Control{}
	cc.Index = 0
	cc.Rate = RateKr
	cc.name = name
	return cc
}

type Mce struct {
	ugens UgenList
}

func NewMce(ugens UgenList) Mce {
	mc := Mce{ugens: ugens}
	return mc
}

type Mrg struct {
	left  UgenType
	right UgenType
}

func NewMrg(left UgenType, right UgenType) Mrg {
	mg := Mrg{left: left, right: right}
	return mg
}

type IConst struct {
	value int
}

func NewIConst(value int) IConst {
	ic := IConst{}
	ic.value = value
	return ic
}

type FConst struct {
	value float64
}

func NewFConst(value float64) FConst {
	fc := FConst{}
	fc.value = value
	return fc
}

func (p Primitive) isUgen() {}
func (p Proxy) isUgen()     {}
func (p Mce) isUgen()       {}
func (p Mrg) isUgen()       {}
func (p Control) isUgen()   {}
func (p IConst) isUgen()    {}
func (p FConst) isUgen()    {}
func (p fromPortC) isUgen() {}
func (p fromPortK) isUgen() {}
func (p fromPortU) isUgen() {}

type NodeC struct {
	id    int
	value float64
}

type NodeK struct {
	id   int
	name string
	Rate int
	Def  int
}

type NodeU struct {
	id      int
	name    string
	Rate    int
	inputs  UgenList
	outputs []int
	Special int
	UgenID  int
}

func (n NodeC) isNode() {}
func (n NodeK) isNode() {}
func (n NodeU) isNode() {}

type Graph struct {
	nextID    int
	constants []NodeC
	controls  []NodeK
	ugens     []NodeU
}

type fromPortC struct {
	portNID int
}

type fromPortK struct {
	portNID int
}

type fromPortU struct {
	portNID int
	portIDX int
}

type input struct {
	u int
	p int
}

type MMap struct {
	cs []int
	ks []int
	us []int
}

var gNextID int = 0

func nextUID() int {
	gNextID = gNextID + 1
	return gNextID
}

func iota(n int, init int, step int) []int {
	if n == 0 {
		return []int{}
	}
	out := []int{init}
	out = append(out, iota(n-1, init+step, step)...)
	return out
}
func extend(ugens UgenList, newlen int) UgenList {
	var ln int
	var out UgenList
	ln = len(ugens)
	if ln > newlen {
		out = ugens[0:newlen]
	} else {
		out = append(ugens, ugens...)
		return extend(out, newlen)
	}
	return out
}

func isSink(ugen interface{}) bool {
	switch ugen.(type) {
	case Primitive:
		if len(ugen.(Primitive).inputs) == 0 {
			return true
		}
		return false

	case Mce:
		for _, elem := range ugen.(Mce).ugens {
			if isSink(elem) {
				return true
			}
		}
	case Mrg:
		if isSink((ugen.(Mrg).left)) {
			return true
		}
	}
	return false
}

func maxNum(nums []int, start int) int {
	max := start
	for _, elem := range nums {
		if elem > max {
			max = elem
		}
	}
	return max
}

func rateOf(ugen interface{}) int {
	switch ugen.(type) {
	case Primitive:
		return ugen.(Primitive).Rate
	case Control:
		return ugen.(Control).Rate
	case Proxy:
		return ugen.(Proxy).primitive.Rate
	case Mrg:
		return rateOf(ugen.(Mrg).left)
	case Mce:
		var rates []int
		for _, elem := range ugen.(Mce).ugens {
			rates = append(rates, rateOf(elem))
		}
		return maxNum(rates, RateKr)
	}
	return RateKr
}

func PrintUgen(ugen UgenType) {
	switch ugen.(type) {
	case IConst:
		fmt.Printf("C: " + string(ugen.(IConst).value))
		break
	case FConst:
		val := ugen.(FConst).value
		fmt.Printf("C: " + strconv.FormatFloat(val, 'E', -1, 64))
		break
	case Control:
		fmt.Printf("K: " + string(ugen.(Control).name))
		break
	case Primitive:
		fmt.Printf("P: " + string(ugen.(Primitive).name))
		break
	default:
		break
	}
}

func mceDegree(ugen UgenType) int {
	switch ugen.(type) {
	case Mrg:
		return mceDegree(ugen.(Mrg).left)
	case Mce:
		return len(ugen.(Mce).ugens)
	default:
		panic("mceDegree")
	}
}

func mceExtend(n int, ugen UgenType) UgenList {
	switch ugen.(type) {
	case Mce:
		return extend(ugen.(Mce).ugens, n)
	case Mrg:
		ex := mceExtend(n, ugen.(Mrg).left)
		if len(ex) > 0 {
			var out UgenList
			out = append(out, ugen)
			out = append(out, ex[1:]...)
			return out
		}
		panic(mceExtend)

	default:
		var out UgenList
		for ind := 0; ind < n; ind = ind + 1 {
			out = append(out, ugen)
		}
		return out
	}
}

func isMce(ugen UgenType) bool {
	switch ugen.(type) {
	case Mce:
		return true
	default:
		return false
	}
}

func ugenFilter(fun func(u UgenType) bool, ugens UgenList) UgenList {
	var out UgenList
	for _, elem := range ugens {
		if fun(elem) {
			out = append(out, elem)
		}
	}
	return out
}

func Transposer(ugens []UgenList) []UgenList {
	len1 := len(ugens)
	len2 := len(ugens[0])
	out := make([]UgenList, len2)
	for ind := range out {
		out[ind] = make(UgenList, len1)
	}
	for ind2 := 0; ind2 < len2; ind2 = ind2 + 1 {
		out1 := out[ind2]
		for ind1 := 0; ind1 < len1; ind1 = ind1 + 1 {
			in1 := ugens[ind1]
			in2 := in1[ind2]
			out1[ind1] = in2
		}
	}
	return out
}

func mceTransform(ugen UgenType) UgenType {
	switch ugen.(type) {
	case Primitive:
		ins := ugenFilter(isMce, ugen.(Primitive).inputs)
		var degs []int
		for _, elem := range ins {
			degs = append(degs, mceDegree(elem))
		}
		upr := maxNum(degs, 0)
		var ext []UgenList
		for _, elem := range ugen.(Primitive).inputs {
			ext = append(ext, mceExtend(upr, elem))
		}
		iet := Transposer(ext)
		var out UgenList
		p := ugen.(Primitive)
		for _, elem := range iet {
			p.inputs = elem
			out = append(out, p)

		}
		return Mce{ugens: out}

	}
	panic("mceTransform")
}

func mceExpand(ugen UgenType) UgenType {
	switch ugen.(type) {
	case Mce:
		var lst UgenList
		for _, elem := range ugen.(Mce).ugens {
			lst = append(lst, mceExpand(elem))
		}
		return Mce{ugens: lst}
	case Mrg:
		lst := mceExpand(ugen.(Mrg).left)
		return Mrg{left: lst, right: ugen.(Mrg).right}
	default:
		rec := func(ugen UgenType) bool {
			switch ugen.(type) {
			case Primitive:
				ins := ugenFilter(isMce, ugen.(Primitive).inputs)
				return len(ins) != 0
			default:
				return false
			}
		}
		if rec(ugen) {
			return mceExpand(mceTransform(ugen))
		}
		return ugen
	}
}

func mceChannel(n int, ugen UgenType) UgenType {
	switch ugen.(type) {
	case Mce:
		return ugen.(Mce).ugens[n]
	default:
		panic("mceChannel")
	}
}

func mceChannels(ugen UgenType) UgenList {
	switch ugen.(type) {
	case Mce:
		return ugen.(Mce).ugens
	case Mrg:
		lst := mceChannels(ugen.(Mrg).left)
		if len(lst) > 1 {
			mrg1 := Mrg{lst[0], ugen.(Mrg).right}
			out := UgenList{mrg1}
			out = append(out, lst[1:]...)
			return out

		}
		panic("mceChannels")

	default:
		return UgenList{ugen}
	}
}

func proxify(ugen UgenType) UgenType {
	switch ugen.(type) {
	case Mce:
		var lst UgenList
		for _, elem := range ugen.(Mce).ugens {
			lst = append(lst, proxify(elem))
		}
		return Mce{ugens: lst}
	case Mrg:
		prx := proxify(ugen.(Mrg).left)
		return Mrg{left: prx, right: ugen.(Mrg).right}
	case Primitive:
		ln := len(ugen.(Primitive).inputs)
		if ln < 2 {
			return ugen
		}
		lst1 := iota(ln, 0, 1)
		lst2 := UgenList{}
		for _, index := range lst1 {
			lst2 = append(lst2, Proxy{primitive: ugen.(Primitive), Index: index})
		}
		return Mce{ugens: lst2}

	default:
		panic("proxify")
	}

}

func mkUgen(rate int, name string, inputs UgenList, outputs []int, ind int, sp int) UgenType {
	pr1 := Primitive{name: name, Rate: rate, inputs: inputs, outputs: outputs, Special: sp, Index: ind}
	return proxify(pr1)
}

func nodeCvalue(node NodeType) float64 {
	return node.(NodeC).value
}

func nodeKdefault(node NodeType) int {
	return node.(NodeK).Def
}

func mkMap(gr Graph) MMap {
	cs := []int{}
	ks := []int{}
	us := []int{}
	for _, elem := range gr.constants {
		cs = append(cs, elem.id)
	}
	for _, elem := range gr.controls {
		ks = append(ks, elem.id)
	}
	for _, elem := range gr.ugens {
		us = append(us, elem.id)
	}
	return MMap{cs: cs, ks: ks, us: us}
}

func fetch(val int, lst []int) int {
	for ind, elem := range lst {
		if elem == val {
			return ind
		}
	}
	return -1
}

func asFromPort(node NodeType) UgenType {
	switch node.(type) {
	case NodeC:
		return fromPortC{portNID: node.(NodeC).id}
	case NodeK:
		return fromPortK{portNID: node.(NodeK).id}
	case NodeU:
		return fromPortU{portNID: node.(NodeU).id, portIDX: 0}
	default:
		panic("asFromPort")

	}
}

func findCP(val float64, node NodeType) bool {
	return val == node.(NodeC).value
}

func pushC(val float64, gr Graph) (NodeType, Graph) {
	node := NodeC{id: gr.nextID + 1, value: val}
	consts := []NodeC{node}
	consts = append(consts, gr.constants...)
	gr1 := Graph{nextID: gr.nextID + 1, constants: consts, controls: gr.controls, ugens: gr.ugens}
	return node, gr1
}

func mkNodeC(ugen UgenType, gr Graph) (NodeType, Graph) {
	var val float64
	switch ugen.(type) {
	case IConst:
		val = float64(ugen.(IConst).value)
	case FConst:
		val = ugen.(FConst).value
	default:
		panic("mkNodeC")
	}
	ln := len(gr.constants)
	for ind := 0; ind < ln; ind = ind + 1 {
		node := gr.constants[ind]
		if findCP(val, node) {
			return node, gr
		}
	}
	return pushC(val, gr)
}

func findKP(str string, node NodeType) bool {
	return node.(NodeK).name == str
}

func pushK(ugen UgenType, gr Graph) (NodeType, Graph) {
	node := NodeK{id: gr.nextID + 1, name: ugen.(Control).name,
		Def: ugen.(Control).Index, Rate: ugen.(Control).Rate}
	contrs := []NodeK{node}
	contrs = append(contrs, gr.controls...)
	gr1 := Graph{nextID: gr.nextID + 1, constants: gr.constants, controls: contrs,
		ugens: gr.ugens}
	return node, gr1
}

func mkNodeK(ugen UgenType, gr Graph) (NodeType, Graph) {
	ln := len(gr.controls)
	name := ugen.(Control).name
	for ind := 0; ind < ln; ind = ind + 1 {
		node := gr.controls[ind]
		if findKP(name, node) {
			return node, gr
		}
	}
	return pushK(ugen, gr)
}

func findUP(rate int, name string, id int, node NodeType) bool {
	if node.(NodeU).Rate == rate && node.(NodeU).name == name &&
		node.(NodeU).id == id {
		return true
	}
	return false
}

func pushU(ugen UgenType, gr Graph) (NodeType, Graph) {
	node := NodeU{id: gr.nextID + 1, name: ugen.(Primitive).name, Rate: ugen.(Primitive).Rate,
		inputs: ugen.(Primitive).inputs, Special: ugen.(Primitive).Special, UgenID: ugen.(Primitive).Index,
		outputs: ugen.(Primitive).outputs}
	ugens := []NodeU{node}
	ugens = append(ugens, gr.ugens...)
	gr1 := Graph{nextID: gr.nextID + 1, constants: gr.constants, controls: gr.controls, ugens: ugens}
	return node, gr1
}

func acc(ll UgenList, nn NodeTypeList, gr Graph) (NodeTypeList, Graph) {
	if len(ll) == 0 {
		//nnlen := len(nn)
		nnr := sort.Reverse(nn).(NodeTypeList)
		/*
			nnr := make(NodeTypeList, nnlen)
			for ind := 0; ind < nnlen; ind = ind + 1 {
				nnr[ind] = nn[nnlen-ind-1]
			}
		*/
		return nnr, gr
	}
	ng1, ng2 := mkNode(ll[0], gr)
	nn = append(nn, ng1)
	return acc(ll[1:len(ll)], nn, ng2)
}

func mkNodeU(ugen UgenType, gr Graph) (NodeType, Graph) {
	switch ugen.(type) {
	case Primitive:
		pr1 := ugen.(Primitive)
		ng1, gnew := acc(pr1.inputs, NodeTypeList{}, gr)
		inputs2 := UgenList{}
		for _, nd := range ng1 {
			inputs2 = append(inputs2, asFromPort(nd))
		}
		rate := pr1.Rate
		name := pr1.name
		index := pr1.Index

		for _, nd2 := range gnew.ugens {
			if findUP(rate, name, index, nd2) {
				return nd2, gnew
			}
		}
		pr := Primitive{name: name, Rate: rate, inputs: inputs2,
			outputs: pr1.outputs, Special: pr1.Special, Index: index}
		return pushU(pr, gnew)
		break
	default:
		break
	}
	panic("mknodeu")
}

func mkNode(ugen UgenType, gr Graph) (NodeType, Graph) {
	switch ugen.(type) {
	case IConst:
		return mkNodeC(ugen, gr)
	case FConst:
		return mkNodeC(ugen, gr)
	case Primitive:
		return mkNodeU(ugen, gr)
	case Control:
		return mkNodeK(ugen, gr)
	case Mrg:
		_, gr1 := mkNode(ugen.(Mrg).right, gr)
		return mkNode(ugen.(Mrg).left, gr1)
	default:
		panic("mkNode")
	}
}

func sc3Implicit(num int) NodeU {
	rates := []int{}
	for index := 1; index < num+1; index++ {
		rates = append(rates, RateKr)
	}

	node := NodeU{id: -1, name: "Control", inputs: UgenList{},
		outputs: rates, UgenID: 0, Rate: RateKr, Special: 0}
	return node
}

func mrgN(lst UgenList) UgenType {
	if len(lst) == 1 {
		return lst[0]
	} else if len(lst) == 2 {
		return Mrg{left: lst[0], right: lst[1]}
	}
	newLst := UgenList{}
	newLst = append(newLst, lst...)
	return Mrg{left: lst[0], right: mrgN(newLst)}
}

func prepareRoot(ugen UgenType) UgenType {
	switch ugen.(type) {
	case Mce:
		return mrgN(ugen.(Mce).ugens)
	case Mrg:
		return Mrg{left: prepareRoot(ugen.(Mrg).left), right: prepareRoot(ugen.(Mrg).right)}
	default:
		break
	}
	return ugen
}

func emptyGraph() Graph {
	return Graph{nextID: 0, constants: []NodeC{}, controls: []NodeK{},
		ugens: []NodeU{}}
}

func synth(ugen UgenType) Graph {
	root := prepareRoot(ugen)
	_, gr := mkNode(root, emptyGraph())
	cs := gr.constants
	ks := gr.controls
	us := gr.ugens
	//reverse us
	us1 := []NodeU{}
	for ind := len(us) - 1; ind >= 0; ind-- {
		us1 = append(us1, us[ind])
	}
	if len(ks) != 0 {
		node := sc3Implicit(len(ks))
		us1 = append([]NodeU{node}, us1...)
	}
	grout := Graph{nextID: -1, constants: cs, controls: ks, ugens: us1}
	return grout
}

func encodeNodeK(mp MMap, node NodeType) []byte {
	out := StrPstr(node.(NodeK).name)
	id1 := fetch(node.(NodeK).id, mp.ks)
	out = append(out, EncodeI16(id1)...)
	return out
}

func encodeInput(inp input) []byte {
	out := EncodeI16(inp.u)
	out = append(out, EncodeI16(inp.p)...)
	return out
}

func mkInput(mm MMap, fp UgenType) input {
	switch fp.(type) {
	case fromPortC:
		p := fetch(fp.(fromPortC).portNID, mm.cs)
		return input{u: -1, p: p}
	case fromPortK:
		p := fetch(fp.(fromPortK).portNID, mm.ks)
		return input{u: 0, p: p}
	case fromPortU:
		u := fetch(fp.(fromPortU).portNID, mm.us)
		return input{u: u, p: fp.(fromPortU).portIDX}
	default:
		panic("mkInput")
	}
}

func encodeNodeU(mm MMap, node NodeType) []byte {
	len1 := len(node.(NodeU).inputs)
	len2 := len(node.(NodeU).outputs)
	out := StrPstr(node.(NodeU).name)
	out = append(out, EncodeI8(node.(NodeU).Rate)...)
	out = append(out, EncodeI16(len1)...)
	out = append(out, EncodeI16(len2)...)
	out = append(out, EncodeI16(node.(NodeU).Special)...)
	for ind := 0; ind < len1; ind = ind + 1 {
		out = append(out, encodeInput(mkInput(mm, node.(NodeU).inputs[ind]))...)
	}
	for ind := 0; ind < len2; ind = ind + 1 {
		out = append(out, EncodeI8(node.(NodeU).outputs[ind])...)
	}
	return out
}

func encodeGraphDef(name string, graph Graph) []byte {
	mm := mkMap(graph)
	out := []byte{}
	out = append(out, EncodeStr("SCgf")...)
	out = append(out, EncodeI32(0)...)
	out = append(out, EncodeI16(1)...)
	out = append(out, StrPstr(name)...)
	out = append(out, EncodeI16(len(graph.constants))...)
	l1 := []float64{}
	for _, elem := range graph.constants {
		l1 = append(l1, nodeCvalue(elem))
	}
	a5 := []byte{}
	for _, elem := range l1 {
		a5 = append(a5, EncodeF32(float32(elem))...)
	}
	out = append(out, a5...)
	out = append(out, EncodeI16(len(graph.controls))...)
	l2 := []int{}
	for _, elem := range graph.controls {
		l2 = append(l2, nodeKdefault(elem))
	}
	a7 := []byte{}
	for _, elem := range l2 {
		a7 = append(a7, EncodeF32((float32)(elem))...)
	}
	out = append(out, a7...)
	out = append(out, EncodeI16(len(graph.controls))...)
	a9 := []byte{}
	for _, elem := range graph.controls {
		a9 = append(a9, encodeNodeK(mm, elem)...)
	}
	out = append(out, a9...)
	out = append(out, EncodeI16(len(graph.ugens))...)
	a10 := []byte{}
	for _, elem := range graph.ugens {
		a10 = append(a10, encodeNodeU(mm, elem)...)
	}

	out = append(out, a10...)

	return out
}

func MkOscMce(rate int, name string, inputs []UgenType, ugen UgenType, ou int) UgenType {
	rl := []int{}
	for ind := 0; ind < ou; ind++ {
		rl = append(rl, rate)
	}
	inps := append(inputs, mceChannels(ugen)...)
	return mkUgen(rate, name, inps, rl, 0, 0)
}

func MkOscID(rate int, name string, inputs []UgenType, ou int) UgenType {
	rl := []int{}
	for ind := 0; ind < ou; ind++ {
		rl = append(rl, rate)
	}

	return mkUgen(rate, name, inputs, rl, nextUID(), 0)
}
func MkOscillator(rate int, name string, inputs []UgenType, ou int) UgenType {
	rl := []int{}
	for ind := 0; ind < ou; ind++ {
		rl = append(rl, rate)
	}

	return mkUgen(rate, name, inputs, rl, 0, 0)
}

func MkFilter(name string, inputs []UgenType, ou int, sp int) UgenType {
	rates := []int{}
	for _, elem := range inputs {
		rates = append(rates, rateOf(elem))
	}
	maxrate := maxNum(rates, RateKr) //check python
	ouList := []int{}
	for ind := 0; ind < ou; ind++ {
		ouList = append(ouList, maxrate)
	}
	return mkUgen(maxrate, name, inputs, ouList, 0, sp)
}
func MkFilterID(name string, inputs []UgenType, ou int, sp int) UgenType {
	rates := []int{}
	for _, elem := range inputs {
		rates = append(rates, rateOf(elem))
	}
	maxrate := maxNum(rates, RateIr)
	ouList := []int{}
	for ind := 0; ind < ou; ind++ {
		ouList = append(ouList, maxrate)
	}
	return mkUgen(maxrate, name, inputs, ouList, nextUID(), sp)
}

func MkFilterMce(name string, inputs []UgenType, ugen UgenType, ou int) UgenType {
	inps := append(inputs, mceChannels(ugen)...)
	return MkFilter(name, inps, ou, 0)
}

func MkOperator(name string, inputs []UgenType, sp int) UgenType {
	rates := []int{}
	for _, elem := range inputs {
		rates = append(rates, rateOf(elem))
	}
	maxrate := maxNum(rates, RateIr)
	outs := []int{maxrate}
	return mkUgen(maxrate, name, inputs, outs, 0, sp)
}

func MkUnaryOperator(sp int, fun interface{}, op interface{}) UgenType {
	switch op.(type) {
	case IConst:
		val := fun.(func(float64) float64)(float64(op.(IConst).value))
		return FConst{value: val}
	case FConst:
		val := fun.(func(float64) float64)(op.(FConst).value)
		return FConst{value: val}
	default:
		break
	}
	ops := []UgenType{}
	switch op.(type) {
	case int:
		ops = append(ops, NewIConst(op.(int)))
		break
	case float64:
		ops = append(ops, NewFConst(op.(float64)))
		break
	}

	return MkOperator("UnaryOpUGen", ops, sp)
}

func MkBinaryOperator(sp int, fun interface{}, op1 interface{}, op2 interface{}) UgenType {
	switch op1.(type) {
	case IConst:
		switch op2.(type) {
		case IConst:
			opp1 := float64(op1.(IConst).value)
			opp2 := float64(op2.(IConst).value)
			val := fun.(func(float64, float64) float64)(opp1, opp2)
			return FConst{value: val}
		case FConst:
			opp1 := float64(op1.(IConst).value)
			opp2 := op2.(FConst).value
			val := fun.(func(float64, float64) float64)(opp1, opp2)
			return FConst{value: val}
		}
	case FConst:
		switch op2.(type) {
		case IConst:
			opp1 := op1.(FConst).value
			opp2 := float64(op2.(IConst).value)
			val := fun.(func(float64, float64) float64)(opp1, opp2)
			return FConst{value: val}
		case FConst:
			opp1 := op1.(FConst).value
			opp2 := op2.(FConst).value
			val := fun.(func(float64, float64) float64)(opp1, opp2)
			return FConst{value: val}
		}

	}

	ops := []UgenType{}
	switch op1.(type) {
	case int:
		ops = append(ops, IConst{value: op1.(int)})
	case float64:
		ops = append(ops, FConst{value: op1.(float64)})
	}
	switch op2.(type) {
	case int:
		ops = append(ops, IConst{value: op2.(int)})
	case float64:
		ops = append(ops, FConst{value: op2.(float64)})
	}
	return MkOperator("BinaryOpUGen", ops, sp)
}
