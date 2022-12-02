package permission

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
	"github.com/sujit-baniya/pkg/str"
	"gorm.io/gorm"
)

var Instance *Engine

type Config struct {
	TableName      string
	SubjectKey     string
	DomainKey      string
	DisableMigrate bool
	Lookup         func(*frame.Context) string
	DomainLookup   func(*frame.Context) string
	Unauthorized   frame.HandlerFunc
	Forbidden      frame.HandlerFunc
	DB             *gorm.DB
	Model          interface{}
	Policy         interface{}
	Adapter        persist.Adapter
}

// Engine holds the configuration for the middleware
type Engine struct {
	*Enforcer
	PolicyAdapter persist.Adapter
	config        Config
}

func Default(cfg Config) (*Engine, error) {
	engine, err := New(cfg)
	Instance = engine
	return engine, err
}

func New(cfg Config) (*Engine, error) {
	var params []any
	if cfg.Model == nil {
		cfg.Model = roleModel()
	}
	params = append(params, cfg.Model)
	if cfg.Policy != nil {
		params = append(params, cfg.Policy)
	}
	if cfg.TableName == "" {
		cfg.TableName = "permissions"
	}
	if cfg.Lookup == nil {
		cfg.Lookup = func(ctx *frame.Context) string {
			if cfg.SubjectKey != "" {
				subjectKey := ctx.Value(cfg.SubjectKey)
				if subjectKey != nil {
					return fmt.Sprintf("%v", subjectKey)
				}
			}
			return ""
		}
	}
	if cfg.DomainLookup == nil {
		cfg.DomainLookup = func(ctx *frame.Context) string {
			if cfg.DomainKey != "" {
				domain := ctx.Value(cfg.DomainKey)
				if domain != nil {
					return domain.(string)
				}
			}

			return ""
		}
	}
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(cc context.Context, c *frame.Context) {
			c.AbortWithJSON(consts.StatusUnauthorized, consts.StatusUnauthorized)
			return
		}
	}
	if cfg.Forbidden == nil {
		cfg.Forbidden = func(cc context.Context, c *frame.Context) {
			c.AbortWithJSON(consts.StatusForbidden, consts.StatusForbidden)
			return
		}
	}

	enforcer, err := casbin.NewEnforcer(params...)
	if err != nil {
		return nil, err
	}
	engine := &Engine{
		Enforcer: &Enforcer{Enforcer: enforcer},
		config:   cfg,
	}
	if cfg.DB != nil {
		if cfg.DisableMigrate {
			gormadapter.TurnOffAutoMigrate(cfg.DB)
		}
		adapter, err := gormadapter.NewAdapterByDBUseTableName(cfg.DB, "", cfg.TableName)
		if err != nil {
			return nil, err
		}
		engine.PolicyAdapter = adapter
		enforcer.SetAdapter(adapter)
	}
	if cfg.Adapter != nil && engine.PolicyAdapter == nil {
		engine.PolicyAdapter = cfg.Adapter
		enforcer.SetAdapter(cfg.Adapter)
	}
	enforcer.EnableAutoSave(true)
	err = enforcer.LoadPolicy()
	return engine, err
}

// RequirePermissions tries to find the current subject and determine if the
// subject has the required permissions according to predefined Casbin policies.
func (cm *Engine) RequirePermissions(permissions []string, opts ...func(o *Options)) frame.HandlerFunc {
	options := &Options{
		ValidationRule:   matchAll,
		PermissionParser: permissionParserWithSeparator(":"),
	}

	for _, o := range opts {
		o(options)
	}

	return func(cc context.Context, c *frame.Context) {
		if len(permissions) == 0 {
			c.Next(cc)
			return
		}

		sub := cm.config.Lookup(c)
		if sub == "" {
			cm.config.Unauthorized(cc, c)
		}

		dom := cm.config.DomainLookup(c)
		if dom == "" {
			dom = "*"
		}
		switch options.ValidationRule {
		case matchAll:
			for _, permission := range permissions {
				vals := append([]string{sub, dom}, options.PermissionParser(permission)...)
				if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
					c.AbortWithJSON(consts.StatusInternalServerError, err.Error())
					return
				} else if !ok {
					cm.config.Forbidden(cc, c)
					return
				}
			}
			c.Next(cc)
			return
		case atLeastOne:
			for _, permission := range permissions {
				vals := append([]string{sub, dom}, options.PermissionParser(permission)...)
				if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
					c.AbortWithJSON(consts.StatusInternalServerError, err.Error())
					return
				} else if ok {
					c.Next(cc)
					return
				}
			}
			cm.config.Forbidden(cc, c)
			return
		}
		c.Next(cc)
		return
	}
}

// Can try to find the current subject and determine if the
// subject has the required permissions according to predefined Casbin policies.
func (cm *Engine) Can(dom, sub, perm string, opts ...func(o *Options)) bool {
	permissions := []string{perm}
	options := &Options{
		ValidationRule:   matchAll,
		PermissionParser: permissionParserWithSeparator(":"),
	}

	for _, o := range opts {
		o(options)
	}
	if len(permissions) == 0 {
		return false
	}
	switch options.ValidationRule {
	case matchAll:
		for _, permission := range permissions {
			vals := append([]string{sub, dom}, options.PermissionParser(permission)...)
			if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
				return false
			} else if !ok {
				return false
			}
		}
		return true
	case atLeastOne:
		for _, permission := range permissions {
			vals := append([]string{sub, dom}, options.PermissionParser(permission)...)
			if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
				return false
			} else if ok {
				return true
			}
		}
		return false
	}
	return false
}

// RoutePermission tries to find the current subject and determine if the
// subject has the required permissions according to predefined Casbin policies.
// This method uses http Path and Method as object and action.
func (cm *Engine) RoutePermission(cc context.Context, c *frame.Context) {
	sub := cm.config.Lookup(c)
	if sub == "" {
		cm.config.Unauthorized(cc, c)
		return
	}
	dom := cm.config.DomainLookup(c)
	if dom == "" {
		dom = "*"
	}
	availableDomains, _ := cm.Enforcer.GetDomainsForUser(sub)
	if !str.Contains(availableDomains, dom) {
		cm.config.Forbidden(cc, c)
		return
	}
	if str.Contains(availableDomains, "*") && dom != "*" {
		dom = "*"
	}
	if ok, err := cm.Enforcer.Enforce(sub, dom, str.FromByte(c.Path()), str.FromByte(c.Method())); err != nil {
		c.AbortWithJSON(consts.StatusInternalServerError, err.Error())
		return
	} else if !ok {
		cm.config.Forbidden(cc, c)
		return
	}

	c.Next(cc)
	return
}

// RequireRoles tries to find the current subject and determine if the
// subject has the required roles according to predefined Casbin policies.
func (cm *Engine) RequireRoles(roles []string, opts ...func(o *Options)) frame.HandlerFunc {
	options := &Options{
		ValidationRule:   matchAll,
		PermissionParser: permissionParserWithSeparator(":"),
	}

	for _, o := range opts {
		o(options)
	}

	return func(cc context.Context, c *frame.Context) {
		if len(roles) == 0 {
			c.Next(cc)
			return
		}

		sub := cm.config.Lookup(c)
		if sub == "" {
			cm.config.Unauthorized(cc, c)
			return
		}
		domain := cm.config.DomainLookup(c)
		if domain == "" {
			domain = "*"
		}
		userRoles := cm.Enforcer.GetRolesForUserInDomain(sub, domain)
		if options.ValidationRule == matchAll {
			for _, role := range roles {
				if !str.Contains(userRoles, role) {
					cm.config.Forbidden(cc, c)
					return
				}
			}
			c.Next(cc)
			return
		} else if options.ValidationRule == atLeastOne {
			for _, role := range roles {
				if str.Contains(userRoles, role) {
					c.Next(cc)
					return
				}
			}
			cm.config.Forbidden(cc, c)
			return
		}

		c.Next(cc)
		return
	}
}
