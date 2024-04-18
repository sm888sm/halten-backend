package roleshierarchy

import "github.com/sm888sm/halten-backend/common/constants/roles"

var RoleHierarchy = map[string]int{
	roles.ObserverRole: 1,
	roles.MemberRole:   2,
	roles.AdminRole:    3,
	roles.OwnerRole:    4,
}
