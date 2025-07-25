
package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"testing"
)

func TestGetUserAuthenticOfCompany(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	var roleRecords = []role.Role{
		//           label              , name           , description                , register_date
		role.NewRole("LABEL_ADMIN" /* */, "Administrator", "administrator description", now),
		role.NewRole("LABEL_MEMBER" /**/, "Member" /*  */, "member description" /*  */, now),
		role.NewRole("LABEL_STAFF" /* */, "Staff" /*   */, "staff description" /*   */, now),
	}
	testUtility.Ready(t, orp, roleRecords)

	var companyRecords = []company.Company{
		//                 persist_key, identifier , name          , register_date
		company.NewCompany(0 /*     */, "CP-TASTAS", "tast company", now),
		company.NewCompany(0 /*     */, "CP-TESTES", "test company", now),
		company.NewCompany(0 /*     */, "CP-TOSTOS", "tost company", now),
	}
	var savedCompanyRecords = testUtility.Ready(t, orp, companyRecords)

	var userRecords = []user.User{
		//           persist_key, identifier , email_idetifier     , name       , bot_flag  , register_date        , update_date
		user.NewUser(0 /*     */, "US-TASTAS", "test01@example.com", "tast name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TESTES", "test02@example.com", "test name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TOSTOS", "test03@example.com", "tost name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TUSTUS", "test04@example.com", "tust name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
	}
	var savedUserRecords = testUtility.Ready(t, orp, userRecords)

	var verifyDate = now.AddDate(0, 0, -3)
	var userEmailRecords = []*user.UserEmail{
		//                persist_key, user_persist_key              , email               , verify_token01  , register_date        , verify_date, expire_date
		user.NewUserEmail(0 /*     */, savedUserRecords[0].PersistKey, "test01@example.com", "verify_token01", now.AddDate(0, -1, 0), &verifyDate, nil),  // x user不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token02", now.AddDate(0, -1, 0), &verifyDate, nil),  // o 対象
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test03@example.com", "verify_token03", now.AddDate(0, -1, 0), nil /*   */, nil),  // x not verified
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test04@example.com", "verify_token04", now.AddDate(0, -1, 0), &verifyDate, &now), // x expired
	}
	testUtility.ReadyPointer(t, orp, userEmailRecords)

	var futureExpireDate = now.AddDate(0, 0, 7)
	var pastExpireDate = now.AddDate(0, 0, -3)
	var records = []*user.UserCompanyRole{
		//                      persist_key, user_persist_key              , company_persist_key              , role_label   , register_date        , expire_date
		user.NewUserCompanyRole(0 /*     */, savedUserRecords[2].PersistKey, savedCompanyRecords[2].PersistKey, "LABEL_MEMBER" /**/, now.AddDate(0, 0, -3), nil),               // x user03 指定されてない
		user.NewUserCompanyRole(0 /*     */, savedUserRecords[1].PersistKey, savedCompanyRecords[1].PersistKey, "LABEL_ADMIN" /* */, now.AddDate(0, 0, -3), nil),               // o user02 expire null
		user.NewUserCompanyRole(0 /*     */, savedUserRecords[1].PersistKey, savedCompanyRecords[1].PersistKey, "LABEL_MEMBER" /**/, now.AddDate(0, 0, -3), &futureExpireDate), // o user02 expire 未来
		user.NewUserCompanyRole(0 /*     */, savedUserRecords[0].PersistKey, savedCompanyRecords[0].PersistKey, "LABEL_MEMBER" /**/, now.AddDate(0, 0, -3), nil),               // x user01
		user.NewUserCompanyRole(0 /*     */, savedUserRecords[1].PersistKey, savedCompanyRecords[1].PersistKey, "LABEL_STAFF" /* */, now.AddDate(0, 0, -3), &pastExpireDate),   // x user02 expire 過去
		user.NewUserCompanyRole(0 /*     */, savedUserRecords[3].PersistKey, savedCompanyRecords[1].PersistKey, "LABEL_MEMBER" /**/, now.AddDate(0, 0, -3), nil),               // o user02 expire null
	}
	testUtility.ReadyPointer(t, orp, records)

	var result, err = orp.GetUserAuthenticOfCompany("CP-TESTES", now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expectEmail = "test02@example.com"
	var expects = []user.UserAuthentic{
		user.UserAuthentic{
			UserPersistKey:     savedUserRecords[1].PersistKey,
			UserIdentifier:     "US-TESTES",
			UserExposeEmailId:  "test02@example.com",
			UserName:           "test name",
			UserBotFlag:        false,
			UserRegisteredDate: now.AddDate(0, -1, 0),
			UserUpdateDate:     now.AddDate(0, 0, -3),
			Email:              &expectEmail,
			UserCompanyRole: []user.UserCompanyRoleFull{
				user.UserCompanyRoleFull{
					UserCompanyRole: user.UserCompanyRole{
						PersistKey:        1,
						UserPersistKey:    savedUserRecords[1].PersistKey,
						CompanyPersistKey: savedCompanyRecords[1].PersistKey,
						RoleLabel:         "LABEL_ADMIN",
						RegisterDate:      now.AddDate(0, 0, -3),
						ExpireDate:        nil,
					},
					UserIdentifier:        "US-TESTES",
					UserExposeEmailId:     "test02@example.com",
					UserName:              "test name",
					UserBotFlag:           false,
					UserRegisteredDate:    now.AddDate(0, -1, 0),
					UserUpdateDate:        now.AddDate(0, 0, -3),
					CompanyIdentifier:     "CP-TESTES",
					CompanyName:           "test company",
					CompanyRegisteredDate: now,
					RoleName:              "Administrator",
					RoleDescription:       "administrator description",
					RoleRegisteredDate:    now,
				},
				user.UserCompanyRoleFull{
					UserCompanyRole: user.UserCompanyRole{
						PersistKey:        2,
						UserPersistKey:    savedUserRecords[1].PersistKey,
						CompanyPersistKey: savedCompanyRecords[1].PersistKey,
						RoleLabel:         "LABEL_MEMBER",
						RegisterDate:      now.AddDate(0, 0, -3),
						ExpireDate:        &futureExpireDate,
					},
					UserIdentifier:        "US-TESTES",
					UserExposeEmailId:     "test02@example.com",
					UserName:              "test name",
					UserBotFlag:           false,
					UserRegisteredDate:    now.AddDate(0, -1, 0),
					UserUpdateDate:        now.AddDate(0, 0, -3),
					CompanyIdentifier:     "CP-TESTES",
					CompanyName:           "test company",
					CompanyRegisteredDate: now,
					RoleName:              "Member",
					RoleDescription:       "member description",
					RoleRegisteredDate:    now,
				},
			},
		},
		user.UserAuthentic{
			UserPersistKey:     savedUserRecords[3].PersistKey,
			UserIdentifier:     "US-TUSTUS",
			UserExposeEmailId:  "test04@example.com",
			UserName:           "tust name",
			UserBotFlag:        false,
			UserRegisteredDate: now.AddDate(0, -1, 0),
			UserUpdateDate:     now.AddDate(0, 0, -3),
			Email:              nil,
			UserCompanyRole: []user.UserCompanyRoleFull{
				user.UserCompanyRoleFull{
					UserCompanyRole: user.UserCompanyRole{
						PersistKey:        2,
						UserPersistKey:    savedUserRecords[3].PersistKey,
						CompanyPersistKey: savedCompanyRecords[1].PersistKey,
						RoleLabel:         "LABEL_MEMBER",
						RegisterDate:      now.AddDate(0, 0, -3),
						ExpireDate:        nil,
					},
					UserIdentifier:        "US-TUSTUS",
					UserExposeEmailId:     "test04@example.com",
					UserName:              "tust name",
					UserBotFlag:           false,
					UserRegisteredDate:    now.AddDate(0, -1, 0),
					UserUpdateDate:        now.AddDate(0, 0, -3),
					CompanyIdentifier:     "CP-TESTES",
					CompanyName:           "test company",
					CompanyRegisteredDate: now,
					RoleName:              "Member",
					RoleDescription:       "member description",
					RoleRegisteredDate:    now,
				},
			},
		},
	}

	testUtility.AssertRecords(t, expects, result, assertSameUserAuthentic)
}
