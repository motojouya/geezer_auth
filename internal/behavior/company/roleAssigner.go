package company

import (
	"github.com/go-gorp/gorp"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterEssence "github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgEssence "github.com/motojouya/geezer_auth/pkg/shelter/essence"
)

type RoleAssignerDB interface {
	gorp.SqlExecutor
	userQuery.GetUserAuthenticQuery
}

type RoleAssigner interface {
	Execute(company shelterCompany.Company, userAuthentic *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error)
}

type RoleAssign struct {
	local localPkg.Localer
	db    RoleAssignerDB
}

func NewRoleAssign(local localPkg.Localer, database RoleAssignerDB) *RoleAssign {
	return &RoleAssign{
		db:    database,
		local: local,
	}
}

func (assigner RoleAssign) Execute(company shelterCompany.Company, userAuthentic *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error) {
	if userAuthentic == nil {
		return nil, pkgEssence.NewNilError("userAuthentic", "userAuthentic is nil")
	}

	now := assigner.local.GetNow()

	userCompanyRole := shelterUser.CreateUserCompanyRole(userAuthentic.GetUser(), company, role, now)
	// FIXME 特定のuserにassignできるcompany、roleの数の制限をこのUserCompanyRoleが知り得ないので、UserAuthentic側でもたせて、その集約で制限をかけるべき。
	// あるいはこれはこのままで、UserAuthenticにaddする関数を追加して、そういう表現にもできる。それがいいかも

	dhUserCompanyRole := dbUser.FromShelterUserCompanyRole(userCompanyRole)

	if err := assigner.db.Insert(dhUserCompanyRole); err != nil {
		return nil, err
	}

	updatedUser := userAuthentic.GetUser().Update(now)
	var dbUserValue = dbUser.FromShelterUser(updatedUser)

	_, err := assigner.db.Update(&dbUserValue)
	if err != nil {
		return nil, err
	}

	dbUserAuthentic, err := assigner.db.GetUserAuthentic(string(userAuthentic.Identifier), now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": string(userAuthentic.Identifier)}
		return nil, shelterEssence.NewNotFoundError("user", keys, "user not found")
	}

	return dbUserAuthentic.ToShelterUserAuthentic()
}
