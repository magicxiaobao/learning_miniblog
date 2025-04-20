package auth

import (
	"fmt"
	"sync"
)

// 简单策略存储，用于示例
type policyStorage struct {
	policies map[string]map[string]map[string]bool // sub -> obj -> act
	roles    map[string]map[string]bool            // user -> role
	mutex    sync.RWMutex
}

// AuthzStub 是一个不依赖Casbin的简单授权实现，用于演示目的
type AuthzStub struct {
	storage *policyStorage
}

// NewAuthzStub 创建一个新的授权存根
func NewAuthzStub() *AuthzStub {
	return &AuthzStub{
		storage: &policyStorage{
			policies: make(map[string]map[string]map[string]bool),
			roles:    make(map[string]map[string]bool),
		},
	}
}

// Authorize 检查用户是否有权限执行操作
func (a *AuthzStub) Authorize(sub, obj, act string) (bool, error) {
	a.storage.mutex.RLock()
	defer a.storage.mutex.RUnlock()

	// 检查用户是否有直接权限
	if subPolicies, ok := a.storage.policies[sub]; ok {
		if objPolicies, ok := subPolicies[obj]; ok {
			if allowed, ok := objPolicies[act]; ok && allowed {
				return true, nil
			}
			if allowed, ok := objPolicies["*"]; ok && allowed {
				return true, nil
			}
		}
	}

	// 检查用户角色是否有权限
	if roles, ok := a.storage.roles[sub]; ok {
		for role := range roles {
			if allowed, _ := a.Authorize(role, obj, act); allowed {
				return true, nil
			}
		}
	}

	return false, nil
}

// AddPolicy 添加策略
func (a *AuthzStub) AddPolicy(sub, obj, act string) (bool, error) {
	a.storage.mutex.Lock()
	defer a.storage.mutex.Unlock()

	if _, ok := a.storage.policies[sub]; !ok {
		a.storage.policies[sub] = make(map[string]map[string]bool)
	}
	if _, ok := a.storage.policies[sub][obj]; !ok {
		a.storage.policies[sub][obj] = make(map[string]bool)
	}

	if a.storage.policies[sub][obj][act] {
		return false, nil // 策略已存在
	}

	a.storage.policies[sub][obj][act] = true
	return true, nil
}

// AddRoleForUser 为用户添加角色
func (a *AuthzStub) AddRoleForUser(user, role string) (bool, error) {
	a.storage.mutex.Lock()
	defer a.storage.mutex.Unlock()

	if _, ok := a.storage.roles[user]; !ok {
		a.storage.roles[user] = make(map[string]bool)
	}

	if a.storage.roles[user][role] {
		return false, nil // 角色已存在
	}

	a.storage.roles[user][role] = true
	return true, nil
}

// GetRolesForUser 获取用户的所有角色
func (a *AuthzStub) GetRolesForUser(user string) ([]string, error) {
	a.storage.mutex.RLock()
	defer a.storage.mutex.RUnlock()

	if roles, ok := a.storage.roles[user]; ok {
		result := make([]string, 0, len(roles))
		for role := range roles {
			result = append(result, role)
		}
		return result, nil
	}

	return []string{}, nil
}

// DeleteRole 删除一个角色及其关联的策略
func (a *AuthzStub) DeleteRole(role string) (bool, error) {
	a.storage.mutex.Lock()
	defer a.storage.mutex.Unlock()

	// 删除角色策略
	delete(a.storage.policies, role)

	// 删除用户角色关联
	for user, roles := range a.storage.roles {
		delete(roles, role)
		if len(roles) == 0 {
			delete(a.storage.roles, user)
		}
	}

	return true, nil
}

// DeleteRoleForUser 删除用户的角色
func (a *AuthzStub) DeleteRoleForUser(user, role string) (bool, error) {
	a.storage.mutex.Lock()
	defer a.storage.mutex.Unlock()

	if roles, ok := a.storage.roles[user]; ok {
		if _, ok := roles[role]; ok {
			delete(roles, role)
			if len(roles) == 0 {
				delete(a.storage.roles, user)
			}
			return true, nil
		}
	}

	return false, nil
}

// GetAllPolicies 获取所有策略
func (a *AuthzStub) GetAllPolicies() [][]string {
	a.storage.mutex.RLock()
	defer a.storage.mutex.RUnlock()

	result := [][]string{}

	for sub, objPolicies := range a.storage.policies {
		for obj, actPolicies := range objPolicies {
			for act, allowed := range actPolicies {
				if allowed {
					result = append(result, []string{sub, obj, act})
				}
			}
		}
	}

	return result
}

// CheckPolicy 检查策略是否存在
func (a *AuthzStub) CheckPolicy(sub, obj, act string) bool {
	a.storage.mutex.RLock()
	defer a.storage.mutex.RUnlock()

	if subPolicies, ok := a.storage.policies[sub]; ok {
		if objPolicies, ok := subPolicies[obj]; ok {
			return objPolicies[act]
		}
	}

	return false
}

// String 打印策略和角色信息
func (a *AuthzStub) String() string {
	a.storage.mutex.RLock()
	defer a.storage.mutex.RUnlock()

	result := "Policies:\n"
	for sub, objPolicies := range a.storage.policies {
		for obj, actPolicies := range objPolicies {
			for act, allowed := range actPolicies {
				if allowed {
					result += fmt.Sprintf("  %s, %s, %s\n", sub, obj, act)
				}
			}
		}
	}

	result += "\nRoles:\n"
	for user, roles := range a.storage.roles {
		result += fmt.Sprintf("  %s: ", user)
		for role := range roles {
			result += fmt.Sprintf("%s ", role)
		}
		result += "\n"
	}

	return result
}
