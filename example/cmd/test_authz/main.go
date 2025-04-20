package main

import (
	"fmt"

	"example/internal/pkg/auth"
)

func main() {
	// 创建授权存根
	authz := auth.NewAuthzStub()

	// 添加策略和角色
	fmt.Println("添加策略和角色")
	fmt.Println("------------------")

	// 添加角色
	authz.AddRoleForUser("alice", "admin")
	authz.AddRoleForUser("bob", "editor")
	authz.AddRoleForUser("charlie", "reader")

	// 添加权限策略
	authz.AddPolicy("admin", "/v1/users/*", "*")
	authz.AddPolicy("admin", "/v1/posts/*", "*")
	authz.AddPolicy("admin", "/v1/system/*", "*")

	authz.AddPolicy("editor", "/v1/posts/*", "*")
	authz.AddPolicy("editor", "/v1/users/self", "GET")

	authz.AddPolicy("reader", "/v1/posts/*", "GET")
	authz.AddPolicy("reader", "/v1/users/self", "GET")

	// 用户特定权限
	authz.AddPolicy("bob", "/v1/users/bob", "*")
	authz.AddPolicy("charlie", "/v1/users/charlie", "*")

	// 打印所有策略和角色
	fmt.Println(authz)

	// 测试授权
	fmt.Println("\n授权测试")
	fmt.Println("------------------")

	testCases := []struct {
		user   string
		path   string
		method string
	}{
		{"alice", "/v1/users/bob", "GET"},
		{"alice", "/v1/posts/123", "POST"},
		{"alice", "/v1/system/config", "PUT"},

		{"bob", "/v1/users/alice", "GET"},
		{"bob", "/v1/users/bob", "PUT"},
		{"bob", "/v1/posts/123", "DELETE"},

		{"charlie", "/v1/users/bob", "GET"},
		{"charlie", "/v1/users/charlie", "PUT"},
		{"charlie", "/v1/posts/123", "GET"},
		{"charlie", "/v1/posts/123", "POST"},
	}

	for _, tc := range testCases {
		allowed, err := authz.Authorize(tc.user, tc.path, tc.method)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}

		result := "拒绝"
		if allowed {
			result = "允许"
		}

		fmt.Printf("用户: %-10s | 路径: %-20s | 方法: %-6s | 结果: %s\n",
			tc.user, tc.path, tc.method, result)
	}

	// 测试角色管理
	fmt.Println("\n角色管理测试")
	fmt.Println("------------------")

	// 获取用户角色
	userRoles, _ := authz.GetRolesForUser("alice")
	fmt.Printf("alice的角色: %v\n", userRoles)

	// 删除角色
	fmt.Println("\n删除alice的admin角色")
	authz.DeleteRoleForUser("alice", "admin")
	userRoles, _ = authz.GetRolesForUser("alice")
	fmt.Printf("alice的角色: %v\n", userRoles)

	// 重新测试alice的权限
	fmt.Println("\n重新测试alice的权限")
	allowed, _ := authz.Authorize("alice", "/v1/users/bob", "GET")
	fmt.Printf("alice访问/v1/users/bob (GET): %v\n", allowed)

	// 添加角色
	fmt.Println("\n为alice添加editor角色")
	authz.AddRoleForUser("alice", "editor")
	userRoles, _ = authz.GetRolesForUser("alice")
	fmt.Printf("alice的角色: %v\n", userRoles)

	// 重新测试alice的权限
	allowed, _ = authz.Authorize("alice", "/v1/posts/123", "POST")
	fmt.Printf("alice访问/v1/posts/123 (POST): %v\n", allowed)
}
