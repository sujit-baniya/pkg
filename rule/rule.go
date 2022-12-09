package rule

import (
	"sort"
	"strings"
	"sync"

	"github.com/sujit-baniya/frame/pkg/common/xid"
	"github.com/sujit-baniya/pkg/timeutil"
)

// Creating rule to work with data of type map[string]any

type ConditionOperator string

type JoinOperator string

const (
	AND JoinOperator = "&&"
	OR  JoinOperator = "||"
	NOT JoinOperator = "!"
)

const (
	EQ          ConditionOperator = "eq"
	NEQ         ConditionOperator = "neq"
	GT          ConditionOperator = "gt"
	LT          ConditionOperator = "lt"
	GTE         ConditionOperator = "gte"
	LTE         ConditionOperator = "lte"
	BETWEEN     ConditionOperator = "between"
	IN          ConditionOperator = "in"
	NotIn       ConditionOperator = "not_in"
	CONTAINS    ConditionOperator = "contains"
	NotContains ConditionOperator = "not_contains"
	StartsWith  ConditionOperator = "starts_with"
	EndsWith    ConditionOperator = "ends_with"
)

type Data map[string]any
type Condition struct {
	Field    string            `json:"field"`
	Operator ConditionOperator `json:"operator"`
	Value    any               `json:"value"`
}

func (condition *Condition) Validate(data Data) bool {
	val, ok := data[condition.Field]
	if !ok {
		return false
	}
	switch condition.Operator {
	case EQ:
		switch val := val.(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case string:
				return strings.EqualFold(val, gtVal)
			}
			return false
		case int:
			switch gtVal := condition.Value.(type) {
			case int:
				return val == gtVal
			case uint:
				return val == int(gtVal)
			case float64:
				return float64(val) == gtVal
			}
			return false
		case float64:
			switch gtVal := condition.Value.(type) {
			case int:
				return val == float64(gtVal)
			case uint:
				return val == float64(gtVal)
			case float64:
				return val == gtVal
			}
			return false
		}
	case NEQ:
		switch val := val.(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case string:
				return !strings.EqualFold(val, gtVal)
			}
			return false
		case int:
			switch gtVal := condition.Value.(type) {
			case int:
				return val != gtVal
			case float64:
				return float64(val) != gtVal
			}
			return false
		case float64:
			switch gtVal := condition.Value.(type) {
			case int:
				return val != float64(gtVal)
			case float64:
				return val != gtVal
			}
			return false
		}

		return false
	case GT:
		switch val := data[condition.Field].(type) {
		case string:
			from, err := timeutil.ParseTime(val)
			if err != nil {
				return false
			}
			switch gtVal := condition.Value.(type) {
			case string:
				smaller, err := timeutil.ParseTime(gtVal)
				if err != nil {
					return false
				}
				return from.After(smaller)
			}
			return false
		case int:
			switch gtVal := condition.Value.(type) {
			case int:
				return val > gtVal
			case float64:
				return float64(val) > gtVal
			}
			return false
		case float64:
			switch gtVal := condition.Value.(type) {
			case int:
				return val > float64(gtVal)
			case float64:
				return val > gtVal
			}
			return false
		}

		return false
	case LT:
		switch val := data[condition.Field].(type) {
		case string:
			from, err := timeutil.ParseTime(val)
			if err != nil {
				return false
			}
			switch gtVal := condition.Value.(type) {
			case string:
				smaller, err := timeutil.ParseTime(gtVal)
				if err != nil {
					return false
				}
				return from.Before(smaller)
			}
			return false
		case int:
			switch ltVal := condition.Value.(type) {
			case int:
				return val < ltVal
			case uint:
				return val < int(ltVal)
			case float64:
				return float64(val) < ltVal
			}
			return false
		case float64:
			switch ltVal := condition.Value.(type) {
			case int:
				return val < float64(ltVal)
			case float64:
				return val < ltVal
			}
			return false
		}

		return false
	case GTE:
		switch val := data[condition.Field].(type) {
		case string:
			from, err := timeutil.ParseTime(val)
			if err != nil {
				return false
			}
			switch gtVal := condition.Value.(type) {
			case string:
				smaller, err := timeutil.ParseTime(gtVal)
				if err != nil {
					return false
				}
				return from.After(smaller) || from.Equal(smaller)
			}
			return false
		case int:
			switch gtVal := condition.Value.(type) {
			case int:
				return val >= gtVal
			case float64:
				return float64(val) >= gtVal
			}
			return false
		case float64:
			switch gtVal := condition.Value.(type) {
			case int:
				return val >= float64(gtVal)
			case float64:
				return val >= gtVal
			}
			return false
		}
		return false
	case LTE:
		switch val := data[condition.Field].(type) {
		case string:
			from, err := timeutil.ParseTime(val)
			if err != nil {
				return false
			}
			switch gtVal := condition.Value.(type) {
			case string:
				smaller, err := timeutil.ParseTime(gtVal)
				if err != nil {
					return false
				}
				return from.Before(smaller) || from.Equal(smaller)
			}
			return false
		case int:
			switch ltVal := condition.Value.(type) {
			case int:
				return val <= ltVal
			case float64:
				return float64(val) <= ltVal
			}
			return false
		case float64:
			switch ltVal := condition.Value.(type) {
			case int:
				return val <= float64(ltVal)
			case float64:
				return val <= ltVal
			}
			return false
		}

		return false
	case BETWEEN:
		switch val := data[condition.Field].(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case []string:
				from, err := timeutil.ParseTime(val)
				if err != nil {
					return false
				}
				start, err := timeutil.ParseTime(gtVal[0])
				if err != nil {
					return false
				}
				last, err := timeutil.ParseTime(gtVal[1])
				if err != nil {
					return false
				}
				return (from.After(start) || from.Equal(start)) && (from.Before(last) || from.Equal(last))
			}
			return false
		case int:
			switch ltVal := condition.Value.(type) {
			case []int:
				return val >= ltVal[0] && val <= ltVal[1]
			case []float64:
				return float64(val) >= ltVal[0] && float64(val) <= ltVal[1]
			}
			return false
		case float64:
			switch ltVal := condition.Value.(type) {
			case []int:
				return val >= float64(ltVal[0]) && val <= float64(ltVal[1])
			case []float64:
				return val >= ltVal[0] && val <= ltVal[1]
			}
			return false
		}

		return false
	case IN:
		switch val := data[condition.Field].(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case []string:
				for _, v := range gtVal {
					if strings.EqualFold(val, v) {
						return true
					}
				}
				return false
			}
			return false
		case int:
			switch ltVal := condition.Value.(type) {
			case []int:
				for _, v := range ltVal {
					if val == v {
						return true
					}
				}
				return false
			}
			return false
		case float64:
			switch ltVal := condition.Value.(type) {
			case []float64:
				for _, v := range ltVal {
					if val == v {
						return true
					}
				}
				return false
			}
			return false
		}

		return false
	case NotIn:
		switch val := data[condition.Field].(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case []string:
				for _, v := range gtVal {
					if strings.EqualFold(val, v) {
						return false
					}
				}
				return true
			}
			return false
		case int:
			switch ltVal := condition.Value.(type) {
			case []int:
				for _, v := range ltVal {
					if val == v {
						return false
					}
				}
				return true
			}
			return false
		case float64:
			switch ltVal := condition.Value.(type) {
			case []float64:
				for _, v := range ltVal {
					if val == v {
						return false
					}
				}
				return true
			}
			return false
		}

		return false
	case CONTAINS:
		switch val := data[condition.Field].(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case string:
				return strings.Contains(val, gtVal)
			}
			return false
		}

		return false
	case NotContains:
		switch val := data[condition.Field].(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case string:
				return !strings.Contains(val, gtVal)
			}
			return false
		}
		return false
	case StartsWith:
		switch val := data[condition.Field].(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case string:
				return strings.HasPrefix(val, gtVal)
			}
			return false
		}
		return false
	case EndsWith:
		switch val := data[condition.Field].(type) {
		case string:
			switch gtVal := condition.Value.(type) {
			case string:
				return strings.HasSuffix(val, gtVal)
			}
			return false
		}
		return false
	}
	return false
}

func NewCondition(field string, operator ConditionOperator, value any) *Condition {
	return &Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

type CallbackFn func(data Data) any

type Node struct {
	Condition []*Condition
	Operator  JoinOperator
	Result    bool
	ID        string
}
type Join struct {
	Left     *Group
	Operator JoinOperator
	Right    *Group
	Result   bool
	ID       string
}
type Group struct {
	Left     *Node
	Operator JoinOperator
	Right    *Node
	Result   bool
	ID       string
}
type Rule struct {
	ID      string
	Handler CallbackFn
	nodes   []*Node
	groups  []*Group
	joins   []*Join
}

func New(id ...string) *Rule {
	rule := &Rule{}
	if len(id) > 0 {
		rule.ID = id[0]
	} else {
		rule.ID = xid.New().String()
	}
	return rule
}
func (r *Rule) addNode(operator JoinOperator, condition ...*Condition) *Node {
	node := &Node{
		Condition: condition,
		Operator:  operator,
	}
	r.nodes = append(r.nodes, node)
	return node
}
func (r *Rule) And(condition ...*Condition) *Node {
	return r.addNode(AND, condition...)
}
func (r *Rule) Or(condition ...*Condition) *Node {
	return r.addNode(OR, condition...)
}
func (r *Rule) Not(condition ...*Condition) *Node {
	return r.addNode(NOT, condition...)
}
func (r *Rule) Group(left *Node, operator JoinOperator, right *Node) *Group {
	group := &Group{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
	r.groups = append(r.groups, group)
	return group
}
func (r *Rule) Join(left *Group, operator JoinOperator, right *Group) *Join {
	join := &Join{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
	r.joins = append(r.joins, join)
	return join
}
func (r *Rule) Apply(d Data, callback ...CallbackFn) any {
	var result, n, g bool
	var defaultCallbackFn = func(data Data) any {
		if data == nil {
			return nil
		}
		return data
	}

	if len(callback) > 0 {
		defaultCallbackFn = callback[0]
	}
	for i, node := range r.nodes {
		if len(node.Condition) == 0 {
			continue
		}
		if i == 0 && node.Operator == AND {
			n = true
		} else if i == 0 && node.Operator == OR {
			n = false
		}
		var nodeResult bool
		switch node.Operator {
		case AND:
			nodeResult = true
			for _, condition := range node.Condition {
				nodeResult = nodeResult && condition.Validate(d)
			}
			n = n && nodeResult
			break
		case OR:
			nodeResult = false
			for _, condition := range node.Condition {
				nodeResult = nodeResult || condition.Validate(d)
			}
			n = n || nodeResult
			break
		}
		node.Result = nodeResult
	}
	if len(r.groups) == 0 {
		result = n
	}
	for i, group := range r.groups {
		if i == 0 && group.Operator == AND {
			g = true
		} else if i == 0 && group.Operator == OR {
			g = false
		}
		var groupResult bool
		switch group.Operator {
		case AND:
			groupResult = group.Left.Result && group.Right.Result
			g = g && groupResult
			break
		case OR:
			groupResult = group.Left.Result || group.Right.Result
			g = g || groupResult
			break
		}
		group.Result = groupResult
	}
	if len(r.groups) > 0 && len(r.joins) == 0 {
		result = g
	}
	for _, join := range r.joins {
		var joinResult bool
		switch join.Operator {
		case AND:
			joinResult = join.Left.Result && join.Right.Result
			break
		case OR:
			joinResult = join.Left.Result || join.Right.Result
			break
		}
		join.Result = joinResult
		result = joinResult
	}
	if !result {
		return defaultCallbackFn(nil)
	}
	if r.Handler != nil {
		return r.Handler(d)
	}
	return defaultCallbackFn(d)
}

type Priority int

const (
	HighestPriority Priority = 1
	LowestPriority  Priority = 0
)

type PriorityRule struct {
	Rule     *Rule
	Priority int
}

type Config struct {
	Rules    []*PriorityRule
	Priority Priority
}

type GroupRule struct {
	Key    string
	Rules  []*PriorityRule
	config Config
	mu     *sync.RWMutex
}

func NewRuleGroup(config ...Config) *GroupRule {
	cfg := Config{}
	if len(config) > 0 {
		cfg = config[0]
	}
	return &GroupRule{
		Rules:  cfg.Rules,
		config: cfg,
		mu:     &sync.RWMutex{},
	}
}

func (r *GroupRule) AddRule(rule *Rule, priority int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Rules = append(r.Rules, &PriorityRule{
		Rule:     rule,
		Priority: priority,
	})
}

func (r *GroupRule) ApplyHighestPriority(data Data, fn ...CallbackFn) any {
	return r.apply(r.sortByPriority("DESC"), data, fn...)
}

func (r *GroupRule) ApplyLowestPriority(data Data, fn ...CallbackFn) any {
	return r.apply(r.sortByPriority(), data, fn...)
}

func (r *GroupRule) Apply(data Data, fn ...CallbackFn) any {
	if r.config.Priority == HighestPriority {
		return r.ApplyHighestPriority(data, fn...)
	}
	return r.ApplyLowestPriority(data, fn...)
}

func (r *GroupRule) apply(sortedRules []*Rule, data Data, fn ...CallbackFn) any {
	for _, rule := range sortedRules {
		response := rule.Apply(data, fn...)
		if response != nil {
			return response
		}
	}
	return nil
}

func (r *GroupRule) SortByPriority(direction ...string) []*Rule {
	return r.sortByPriority(direction...)
}

func (r *GroupRule) sortByPriority(direction ...string) []*Rule {
	dir := "ASC"
	if len(direction) > 0 {
		dir = direction[0]
	}
	if dir == "DESC" {
		sort.Sort(sort.Reverse(byPriority(r.Rules)))
	} else {
		sort.Sort(byPriority(r.Rules))
	}
	res := make([]*Rule, 0, len(r.Rules))
	for _, q := range r.Rules {
		res = append(res, q.Rule)
	}
	return res
}

type byPriority []*PriorityRule

func (x byPriority) Len() int           { return len(x) }
func (x byPriority) Less(i, j int) bool { return x[i].Priority < x[j].Priority }
func (x byPriority) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
