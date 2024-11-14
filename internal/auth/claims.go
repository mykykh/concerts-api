package auth

import (
    "github.com/gofrs/uuid"
    "github.com/coreos/go-oidc/v3/oidc"
)

type ResourceRoles struct {
    Roles []string `json:"roles"`
}

type Claims struct {
    ID uuid.UUID `json:"sub"`
    Username string `json:"prefered_username"`
    Email string `json:"email"`
    FullName string `json:"full-name"`
    EmailVerified bool `json:"email_verified"`
    ResourceAccess map[string]ResourceRoles `json:"resource_access"`
}

func TokenToClaims(token *oidc.IDToken) (*Claims, error) {
    var temp struct {
        ID string `json:"sub"`
        Username string `json:"prefered_username"`
        Email string `json:"email"`
        FullName string `json:"full-name"`
        EmailVerified bool `json:"email_verified"`
        ResourceAccess map[string]ResourceRoles `json:"resource_access"`
    }

    if err := token.Claims(&temp); err != nil {
        return nil, err
    }

    uuid_id, err := uuid.FromString(temp.ID)
    if err != nil {
        return nil, err
    }

    return &Claims{
        ID: uuid_id,
        Username: temp.Username,
        Email: temp.Email,
        FullName: temp.FullName,
        EmailVerified: temp.EmailVerified,
        ResourceAccess: temp.ResourceAccess,
    }, nil
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
