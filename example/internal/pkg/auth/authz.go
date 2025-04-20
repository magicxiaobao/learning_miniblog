package auth

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// 定义一个简单的RBAC模型
const rbacModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

// Authz 授权器，提供RBAC访问控制
type Authz struct {
	enforcer *casbin.SyncedEnforcer
}

// NewAuthz 创建一个新的授权器
func NewAuthz(db *gorm.DB) (*Authz, error) {
	// 创建基于gorm的数据库适配器
	a, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 从字符串加载模型
	m, err := model.NewModelFromString(rbacModel)
	if err != nil {
		return nil, err
	}

	// 创建强制执行器
	enforcer, err := casbin.NewSyncedEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	// 启动自动加载策略
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	return &Authz{enforcer: enforcer}, nil
}

// Authorize 检查用户是否有权限执行操作
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加策略
func (a *Authz) AddPolicy(sub, obj, act string) (bool, error) {
	return a.enforcer.AddPolicy(sub, obj, act)
}

// AddRoleForUser 为用户添加角色
func (a *Authz) AddRoleForUser(user, role string) (bool, error) {
	return a.enforcer.AddGroupingPolicy(user, role)
}

// GetRolesForUser 获取用户的所有角色
func (a *Authz) GetRolesForUser(user string) ([]string, error) {
	return a.enforcer.GetRolesForUser(user)
}

// DeleteRole 删除一个角色及其关联的策略
func (a *Authz) DeleteRole(role string) (bool, error) {
	// 删除角色的所有策略
	if _, err := a.enforcer.RemoveFilteredPolicy(0, role); err != nil {
		return false, err
	}
	// 删除所有与该角色相关的用户分配
	if _, err := a.enforcer.RemoveFilteredGroupingPolicy(1, role); err != nil {
		return false, err
	}
	return true, nil
}

// DeleteRoleForUser 删除用户的角色
func (a *Authz) DeleteRoleForUser(user, role string) (bool, error) {
	return a.enforcer.RemoveGroupingPolicy(user, role)
}

// GetAllPolicies 获取所有策略
func (a *Authz) GetAllPolicies() [][]string {
	policies, _ := a.enforcer.GetPolicy()
	return policies
}

// CheckPolicy 检查策略是否存在
func (a *Authz) CheckPolicy(sub, obj, act string) bool {
	has, _ := a.enforcer.HasPolicy(sub, obj, act)
	return has
}
