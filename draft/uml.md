
# UML
的な情報。実際にUMLを使うかはわからない。
サロゲートキーはオブジェクトに持たない。expose_idでできるとは思うので。
必要になった場合は、変更する。

```go
struct User {
  expose_id: string
  expose_email_id: string
  name: string
  bot_flag: bool
  email: *string
  companyRole: *CompanyRole
}
CreateUser(name: string, emailId: string, botFlag: bool): User
NewUser(expose_id: string, name: string, emailId: string, email: *string, botFlag: bool, companyRole: *CompanyRole): User
```

```go
type Password = string
getPassword(password: string): Password
verifyPassword(password: Password): bool
// サンプルコードが多く、まー安全でデファクトな感じのbcrypt hashを使う
```

```go
type RefreshToken = string
generateRefreshToken(): RefreshToken
// tokenのverifyはequalityで良い
```

```go
type AccessToken = string
generateAccessToken(user: User): AccessToken
getUserFromAccessToken(token: AccessToken): User
```

```go
type CompanyInviteToken = string
generateCompanyInviteToken(): CompanyInviteToken
// tokenのverifyはequalityで良い
```

```go
struct Company {
  expose_id: string
  name: string
}
CreateCompany(name: string): Company
NewCompany(expose_id: string, name: string): Company
```

```go
struct Role {
  name: string
  label: string
  description: string
}
NewRole(name: string, label: string, description: string): Role
```

```go
struct CompanyRole {
  company: Company
  role: Role
}
NewCompanyRole(company: Company, role: Role): CompanyRole
```

実装順だが、ほとんどコンストラクタなので、company, role, user, credentialsの順で実装する。
credentialsは、ライブラリ利用もあるので、後

