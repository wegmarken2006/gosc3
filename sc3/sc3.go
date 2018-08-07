package sc3

import (
	. "gosc3/osc"
)

const (
	RateKr = 0
	RateIr = 1
	RateAr = 2
	RateDr = 3
)

type UgenType interface {
	isUgen()
}

type nodeType interface {
	isNode()
}

type Primitive struct {
	Rate int
	name string
	//inputs  []interface{}
	inputs  []UgenType
	outputs []int
	Special int
	Index   int
}

// NewPrimitive primitive constructor
func NewPrimitive(name string, inputs []UgenType, outputs []int) Primitive {
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
	ugens []UgenType
}

func NewMce(ugens []UgenType) Mce {
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
	value float32
}

func NewFConst(value float32) FConst {
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
	value float32
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
	inputs  []UgenType
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
func extend(ugens []UgenType, newlen int) []UgenType {
	var ln int
	var out []UgenType
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

func mceExtend(n int, ugen UgenType) []UgenType {
	switch ugen.(type) {
	case Mce:
		return extend(ugen.(Mce).ugens, n)
	case Mrg:
		ex := mceExtend(n, ugen.(Mrg).left)
		if len(ex) > 0 {
			var out []UgenType
			out = append(out, ugen)
			out = append(out, ex[1:]...)
			return out
		}
		panic(mceExtend)

	default:
		var out []UgenType
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

func ugenFilter(fun func(u UgenType) bool, ugens []UgenType) []UgenType {
	var out []UgenType
	for _, elem := range ugens {
		if fun(elem) {
			out = append(out, elem)
		}
	}
	return out
}

func Transposer(ugens [][]UgenType) [][]UgenType {
	len1 := len(ugens)
	len2 := len(ugens[0])
	out := make([][]UgenType, len2)
	for ind := range out {
		out[ind] = make([]UgenType, len1)
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
		var ext [][]UgenType
		for _, elem := range ugen.(Primitive).inputs {
			ext = append(ext, mceExtend(upr, elem))
		}
		iet := Transposer(ext)
		var out []UgenType
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
		var lst []UgenType
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

func mceChannels(ugen UgenType) []UgenType {
	switch ugen.(type) {
	case Mce:
		return ugen.(Mce).ugens
	case Mrg:
		lst := mceChannels(ugen.(Mrg).left)
		if len(lst) > 1 {
			mrg1 := Mrg{lst[0], ugen.(Mrg).right}
			out := []UgenType{mrg1}
			out = append(out, lst[1:]...)
			return out

		}
		panic("mceChannels")

	default:
		return []UgenType{ugen}
	}
}

func proxify(ugen UgenType) UgenType {
	switch ugen.(type) {
	case Mce:
		var lst []UgenType
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
		lst2 := []UgenType{}
		for _, index := range lst1 {
			lst2 = append(lst2, Proxy{primitive: ugen.(Primitive), Index: index})
		}
		return Mce{ugens: lst2}

	default:
		panic("proxify")
	}

}

func mkUgen(rate int, name string, inputs []UgenType, outputs []int, ind int, sp int) UgenType {
	pr1 := Primitive{name: name, Rate: rate, inputs: inputs, outputs: outputs, Special: sp, Index: ind}
	return proxify(pr1)
}

func nodeCvalue(node nodeType) float32 {
	return node.(NodeC).value
}

func nodeKdefault(node nodeType) int {
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

func encodeNodeK(mp MMap, node nodeType) []byte {
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
		return input{u: u, p: fp.(fromPortK).portNID}
	default:
		panic("mkInput")
	}
}

func asFromPort(node nodeType) UgenType {
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

func encodeNodeU(mm MMap, node nodeType) []byte {
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

func findCP(val float32, node nodeType) bool {
	return val == node.(NodeC).value
}

func pushC(val float32, gr Graph) (nodeType, Graph) {
	node := NodeC{id: gr.nextID + 1, value: val}
	consts := []NodeC{node}
	consts = append(consts, gr.constants...)
	gr1 := Graph{nextID: gr.nextID + 1, constants: consts, controls: gr.controls, ugens: gr.ugens}
	return node, gr1
}

func mkNodeC(ugen UgenType, gr Graph) (nodeType, Graph) {
	var val float32
	switch ugen.(type) {
	case IConst:
		val = float32(ugen.(IConst).value)
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

func findKP(str string, node nodeType) bool {
	return node.(NodeK).name == str
}

func pushK(ugen UgenType, gr Graph) (nodeType, Graph) {
	node := NodeK{id: gr.nextID + 1, name: ugen.(Control).name,
		Def: ugen.(Control).Index, Rate: ugen.(Control).Rate}
	contrs := []NodeK{node}
	contrs = append(contrs, gr.controls...)
	gr1 := Graph{nextID: gr.nextID + 1, constants: gr.constants, controls: contrs,
		ugens: gr.ugens}
	return node, gr1
}

func mkNodeK(ugen UgenType, gr Graph) (nodeType, Graph) {
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

func findUP(rate int, name string, id int, node nodeType) bool {
	if node.(NodeU).Rate == rate && node.(NodeU).name == name &&
		node.(NodeU).id == id {
		return true
	}
	return false
}

func pushU(ugen UgenType, gr Graph) (nodeType, Graph) {
	intrates := []int{}
	for _, elem := range ugen.(Primitive).outputs {
		intrates = append(intrates, rateOf(elem))
	}
	node := NodeU{id: gr.nextID + 1, name: ugen.(Primitive).name, Rate: ugen.(Primitive).Rate,
		inputs: ugen.(Primitive).inputs, Special: ugen.(Primitive).Special, UgenID: ugen.(Primitive).Index,
		outputs: intrates}
	ugens := []NodeU{node}
	ugens = append(ugens, gr.ugens...)
	gr1 := Graph{nextID: gr.nextID + 1, constants: gr.constants, controls: gr.controls, ugens: ugens}
	return node, gr1
}

/*
func acc(ll []UgenType, nn []nodeType, gr graph) ([]nodeType, graph) {
	if len(ll) == 0 {
		nnlen := len(nn)
		nnr := make([]nodeType, nnlen)
		for ind := 0; ind < nnlen; ind = ind + 1 {
			nnr[ind] = nn[nnlen-ind-1]
		}
		return nnr, gr
	}
	ng1, ng2 := mkNode(ll[0], gr)

	//TODO
}

func mkNode(ugen UgenType, gr graph) (nodeType, graph) {
	switch ugen.(type) {
	case iConst:
		return mkNodeC(ugen, gr)
	case fConst:
		return mkNodeC(ugen, gr)
	case primitive:
		return mkNodeU(ugen, gr)
	case mrg:
		_, gr1 := mkNode(ugen.(mrg).right, gr)
		return mkNode(ugen.(mrg).left, gr1)
	default:
		panic("mkNode")
	}
}

*/
