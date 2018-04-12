package main

import "fmt"

type baseInfo struct {
	workspace  string
	versionSet string
}

func (b baseInfo) String() string {
	return fmt.Sprintf("baseInfo{workspace: \"%s\", versionSet: \"%s\"}", b.workspace, b.versionSet)
}

type packageDecl struct {
	name string
	loc  string
}

func (p packageDecl) String() string {
	return fmt.Sprintf("packageDecl{name: \"%s\", loc: \"%s\"}", p.name, p.loc)
}

type packageInfo struct {
	base     baseInfo
	packages []packageDecl
}

func (p packageInfo) String() string {
	var ps string
	if len(p.packages) == 0 {
		ps = "[]"
	} else {
		pss := p.packages[0].String()
		for _, pd := range p.packages[1:] {
			pss += "," + pd.String()
		}
		ps = fmt.Sprintf("[%s]", pss)
	}
	return fmt.Sprintf("packageInfo{base: %s, packages: %s}", p.base.String(), ps)
}

type result interface {
	isError() bool
}

type marshalErr struct {
	err string
}

func (m marshalErr) Error() string {
	return m.err
}

func (m marshalErr) isError() bool {
	return true
}

func (p packageInfo) isError() bool {
	return false
}

// Node Semantic Analysis Functions
type semantErr struct {
	err string
}

func (e semantErr) Error() string {
	return e.err
}

func (n nerror) semant() error {
	return semantErr{n.Error()}
}

func (n ndecl) semant() error {
	if err := n.ident.semant(); err != nil {
		return err
	}

	if err := n.val.semant(); err != nil {
		return err
	}

	return nil
}

func (n nobject) semant() error {
	for _, d := range n.decls {
		if err := d.semant(); err != nil {
			return err
		}
	}
	return nil
}

func (n nterm) semant() error {
	switch n.Type() {
	case nodeString:
		return nil
	case nodeIdent, nodeNumber, nodePath:
		if len(n.val) == 0 {
			return semantErr{fmt.Sprintf("invalid %s: '%s'", n.typ.String(), n.val)}
		}
		return nil
	default:
		return nil
	}
}

func (n nfile) semant() error {
	for _, d := range n.decls {
		if err := d.semant(); err != nil {
			return err
		}
	}
	return nil
}

// Marshaller Functions
func marshalPackages(m map[string]interface{}) []packageDecl {
	var ps []packageDecl

	for k, v := range m {
		p := packageDecl{
			name: k,
			loc:  v.(string),
		}
		ps = append(ps, p)
	}

	return ps
}

func marshalBaseInfo(m map[string]interface{}) baseInfo {
	var ws string
	var vs string

	if wss, ok := m["workspace"]; ok {
		ws = wss.(string)
	}

	if vss, ok := m["versionSet"]; ok {
		vs = vss.(string)
	}

	b := baseInfo{
		workspace:  ws,
		versionSet: vs,
	}
	return b
}

func marshalPackageInfo(m map[string]interface{}) packageInfo {
	var b baseInfo
	var ps []packageDecl

	if bb, ok := m["base"]; ok {
		b = marshalBaseInfo(bb.(map[string]interface{}))
	}

	if pss, ok := m["packages"]; ok {
		ps = marshalPackages(pss.(map[string]interface{}))
	}

	p := packageInfo{
		base:     b,
		packages: ps,
	}
	return p
}

func marshalObj(n node) map[string]interface{} {
	o := n.(nobject)
	m := make(map[string]interface{})

	for _, d := range o.decls {
		nd := d.(ndecl)
		k := nd.ident.(nterm).val
		v := nd.val
		if v.Type() == nodeObj {
			m[k] = marshalObj(v)
		} else {
			m[k] = v.(nterm).val
		}
	}

	return m
}

func marshalFile(n node) packageInfo {
	f := n.(nfile)
	res := make(map[string]interface{})

	for _, d := range f.decls {
		nd := d.(ndecl)
		k := nd.ident.(nterm).val
		v := nd.val
		if v.Type() == nodeObj {
			res[k] = marshalObj(v)
		} else {
			res[k] = v.(nterm).val
		}
	}

	p := marshalPackageInfo(res)
	return p
}

type marshaller struct {
	input  chan node
	output chan result
}

func (m *marshaller) run() {
	n := <-m.input
	if n.Type() == nodeErr {
		m.output <- marshalErr{n.(nerror).Error()}
		close(m.output)
		return
	}

	if err := n.semant(); err != nil {
		m.output <- marshalErr{err.Error()}
		close(m.output)
		return
	}

	// Build Data Structures
	p := marshalFile(n)
	m.output <- p
	close(m.output)
}

func newMarshaller(input chan node) *marshaller {
	m := marshaller{
		input:  input,
		output: make(chan result),
	}
	return &m
}

func marshal(input chan node) chan result {
	m := newMarshaller(input)
	go m.run()
	return m.output
}
