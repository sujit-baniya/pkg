package permission

import "github.com/casbin/casbin/v2/model"

func roleModel() interface{} {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, dom, obj, act")                                                                                                                                                                               // [request_definition]
	m.AddDef("p", "p", "sub, dom, obj, act")                                                                                                                                                                               // [policy_definition]
	m.AddDef("g", "g", "_, _, _")                                                                                                                                                                                          // [role_definition]
	m.AddDef("g2", "g2", "_, _")                                                                                                                                                                                           // [role_definition]
	m.AddDef("e", "e", "some(where (p.eft == allow))")                                                                                                                                                                     // [policy_effect]
	m.AddDef("m", "m", "g(r.sub, p.sub, r.dom) && g2(r.dom, p.dom) && regexMatch(r.act, p.act) && (r.dom == p.dom || regexMatch(r.dom, p.dom)) && (r.obj == p.obj || regexMatch(r.obj, p.obj) || keyMatch(r.obj, p.obj))") // [matchers]
	return m
}
