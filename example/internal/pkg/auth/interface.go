package auth

// Authorizer 定义授权接口，包含了认证授权所需的所有方法
type Authorizer interface {
	// Authorize 检查用户是否有权限执行操作
	Authorize(sub, obj, act string) (bool, error)

	// AddPolicy 添加策略
	AddPolicy(sub, obj, act string) (bool, error)

	// AddRoleForUser 为用户添加角色
	AddRoleForUser(user, role string) (bool, error)

	// GetRolesForUser 获取用户的所有角色
	GetRolesForUser(user string) ([]string, error)

	// DeleteRole 删除一个角色及其关联的策略
	DeleteRole(role string) (bool, error)

	// DeleteRoleForUser 删除用户的角色
	DeleteRoleForUser(user, role string) (bool, error)

	// GetAllPolicies 获取所有策略
	GetAllPolicies() [][]string

	// CheckPolicy 检查策略是否存在
	CheckPolicy(sub, obj, act string) bool
}
