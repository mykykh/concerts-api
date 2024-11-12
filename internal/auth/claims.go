package auth

type ResourceRoles struct {
    Roles []string `json:"roles"`
}

type Claims struct {
    ID string `json:"sub"`
    Username string `json:"prefered_username"`
    Email string `json:"email"`
    EmailVerified bool `json:"email_verified"`
    ResourceAccess map[string]ResourceRoles `json:"resource_access"`
}

func (claims Claims) HasResourceAccessRole(resource string, role string) bool {
    resourceRoles, resourceExist := claims.ResourceAccess[resource]
    if !resourceExist {
        return false
    }

    for _, r := range resourceRoles.Roles {
        if r == role {
            return true
        }
    }
    return false
}
