package permission

import (
	"github.com/sujit-baniya/pkg/str"
)

type OperationPermission struct {
	Operation  string
	Permission string
}

type RoutePermission struct {
	Route  string
	Method string
}

type UserRole struct {
	User string
	Role string
}

type Enforcer struct {
	*casbin.Enforcer
}

func (c *Enforcer) GetDomainsForUser(user string) ([]string, error) {
	domains, err := c.GetAllDomains()
	if err != nil {
		return nil, err
	}
	userDomains, err := c.Enforcer.GetDomainsForUser(user)
	if err != nil {
		return nil, err
	}
	if str.Contains(userDomains, "*") {
		return domains, nil
	}
	return userDomains, nil
}

func (c *Enforcer) GetRolesInDomain(domain string) []string {
	rs := c.GetFilteredPolicy(0, "", domain, "")
	var r []string
	mapRoles := map[string]struct{}{}

	for i := range rs {
		roleNames := rs[i][0]
		if _, ok := mapRoles[roleNames]; !ok {
			mapRoles[roleNames] = struct{}{}
			r = append(r, roleNames)
		}
	}
	return r
}

func (c *Enforcer) GetAllPermissionsForRole(domain string, role string) [][]string {
	ps := c.GetPermissionsForUserInDomain(role, domain)
	var p [][]string
	for i := range ps {
		p = append(p, []string{ps[i][2], ps[i][3]})
	}
	return p
}

func (c *Enforcer) GetPermissionsForRole(domain string, role string) (optPermissions []OperationPermission) {
	ps := c.GetPermissionsForUserInDomain(role, domain)
	methods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}
	for i := range ps {
		if !str.Contains(methods, ps[i][3]) {
			optPermissions = append(optPermissions, OperationPermission{
				Operation:  ps[i][2],
				Permission: ps[i][3],
			})
		}
	}
	return
}

func (c *Enforcer) GetRoutePermissionsForRole(domain string, role string) (routePermissions []RoutePermission) {
	ps := c.GetPermissionsForUserInDomain(role, domain)
	methods := []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}
	for i := range ps {
		if str.Contains(methods, ps[i][3]) {
			routePermissions = append(routePermissions, RoutePermission{
				Route:  ps[i][2],
				Method: ps[i][3],
			})
		}
	}
	return
}

func (c *Enforcer) GetAllUsersInDomainWithRole(domain string) []UserRole {
	var urs []UserRole
	users := c.GetFilteredGroupingPolicy(0, "", "", domain)
	for _, user := range users {
		ur := UserRole{
			User: user[0],
			Role: user[1],
		}
		urs = append(urs, ur)
	}
	return urs
}

func (c *Enforcer) GetUserRole(domain, user string) []string {
	var roleUser []string
	roles := c.GetFilteredGroupingPolicy(0, user, "", domain)
	for _, role := range roles {
		roleUser = append(roleUser, role[1])
	}
	return roleUser
}

func (c *Enforcer) GetModuleRelatedByRole(domain, role string) []string {
	var modules []string
	mapModules := map[string]struct{}{}
	policies := c.GetFilteredPolicy(0, role, domain, "", "")
	for _, policy := range policies {
		moduleName := policy[2]
		if _, ok := mapModules[moduleName]; !ok {
			mapModules[moduleName] = struct{}{}
			modules = append(modules, moduleName)
		}
	}
	return modules
}
