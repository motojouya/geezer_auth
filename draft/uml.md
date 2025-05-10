
# UML
的な情報。実際にUMLを使うかはわからない。

## model

```go
struct LinboUser {
  expose_id: string
  expose_email_id: string
  name: string
  bot_flag: bool
  email: *string
  companyRole: *CompanyRole
}

struct User {
  user_id: uint
  LinboUser
}
CreateUser(name: string, emailId: string, botFlag: bool): LinboUser
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
publishAccessToken(user: User, tokens: []AccessToken): AccessToken
getUserFromAccessToken(token: AccessToken): User
```

```go
type CompanyInviteToken = string
generateCompanyInviteToken(): CompanyInviteToken
// tokenのverifyはequalityで良い
```

```go
struct LinboCompany {
  expose_id: string
  name: string
}
struct Company {
  company_id: uint
  LinboCompany
}
CreateCompany(name: string): LinboCompany
NewCompany(expose_id: string, name: string): Company
```

```go
struct Role {
  role_id: uint
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

## query
特定tableの1レコードについての主キーでの操作は便利関数を用意して行う。
できれば、unique keyでの操作も用意したいが、gorpの機能的に、unique keyの扱いを調べてから判断という感じ。
履歴的なテーブル構造になってる場合は、主キーだけでは取得できないね。

- password
  - get password
  - bulk expire password
- email
  - get email
  - bulk expire email
- refresh_token
  - get refresh token
  - bulk expire refresh token
- access_token
  - get access tokens
- company
  - get company users
  - get company roles
- user
  - get user all

## 実装順序
この順番かな
- model
- query
- procedure
- middleware
- router

それぞれ、utilityも定義した後かな。modelは無いけど

