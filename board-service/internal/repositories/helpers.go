package repositories

import "github.com/sm888sm/halten-backend/common/constants/roleshierarchy"

func canAssignRole(currentRole string, targetRole string) bool {

	// Check if the roles exist
	currentRoleValue, currentRoleExists := roleshierarchy.RoleHierarchy[currentRole]
	targetRoleValue, targetRoleExists := roleshierarchy.RoleHierarchy[targetRole]

	if !currentRoleExists || !targetRoleExists {
		return false
	}

	// A user can't downgrade a role which is equal or higher than its role
	return currentRoleValue > targetRoleValue
}
