
# UML
的な情報。実際にUMLを使うかはわからない。

## dir
- internal
  - model
  - procedure
  - db
   - query
   - transfer
     - in
     - out
  - http
    - middleware
    - router
    - transfer
      - in
      - out
  - utility
    - db
    - http

## model

```go
struct UnsavedUser {
  expose_id: string
  expose_email_id: string
  name: string
  bot_flag: bool
}
struct User {
  user_id: uint
  companyRole: *CompanyRole
  email: *string
  registeredDate: time.Time
  UnsavedUser
}
func CreateUser(name: string, emailId: string, botFlag: bool): UnsavedUser
func NewUser(expose_id: string, name: string, emailId: string, email: *string, botFlag: bool, companyRole: *CompanyRole): User
```

```go
struct CompanyRole {
  company: Company
  role: Role
}
func NewCompanyRole(company: Company, role: Role): CompanyRole
```

```go
struct UnsavedCompany {
  expose_id: string
  name: string
}
struct Company {
  company_id: uint
  UnsavedCompany
  registeredDate: time.Time
  roles: []Role
}
func CreateCompany(name: string): UnsavedCompany
func NewCompany(expose_id: string, name: string, roles: []Role): Company
```

```go
struct UnsavedRole {
  name: string
  label: string
  description: string
}
struct Role {
  role_id: uint
  UnsavedRole
  registeredDate: time.Time
}
func CreateRole(name: string, label: string, description: string): UnsavedRole
func NewRole(name: string, label: string, description: string): Role
```

ここまで作った

```go
type Password string
func getPassword(password: string): Password
func verifyPassword(password: Password): bool
// サンプルコードが多く、まー安全でデファクトな感じのbcrypt hashを使う
```

```go
type RefreshToken string
func generateRefreshToken(): RefreshToken
// tokenのverifyはequalityで良い
```

```go
struct AccessToken {
  token: string
  expire_at: time.Time
}
// expireの基準日がいるので、日付が必要。オプションで期間の調整ができてもいいかもしれない
func publishAccessToken(user: User, tokens: []AccessToken, date: Date): AccessToken
func getUserFromAccessToken(token: string): User
```

```go
type CompanyInviteToken string
func generateCompanyInviteToken(): CompanyInviteToken
// tokenのverifyはequalityで良い
```


実装順だが、ほとんどコンストラクタなので、company, role, user, credentialsの順で実装する。
credentialsは、ライブラリ利用もあるので、後

## query
特定tableの1レコードについての主キーでの操作は便利関数を用意して行う。
できれば、unique keyでの操作も用意したいが、gorpの機能的に、unique keyの扱いを調べてから判断という感じ。
履歴的なテーブル構造になってる場合は、主キーだけでは取得できないね。

```go
func getPassword(userId: uint): Password
func bulkExpirePassword(userId: uint): bool
```

```go
// 主キーが別なので、user_idとemailで一位に特定する感じ。用意しなくてもいけるかも
func getEmail(userId: uint, email: string): string
func getUserEmail(userId: uint): string
func bulkExpireEmail(userId: uint): bool
```

```go
func getRefreshToken(userId: uint): RefreshToken
func bulkExpireRefreshToken(userId: uint): bool
```

```go
func getAccessToken(userId: uint): []AccessToken
```

```go
struct FullCompany {
  company: Company
  roles: []Role
}
func getFullCompany(companyExposeId: string): FullCompany
func getCompanyRoles(companyId: uint): []Role
func getCompanyUsers(companyId: uint): []FullUser
```

```go
struct FullUser {
  user: User // DB tableのuser
  userEmail: UserEmail
  userRole: UserRole
  company: Company
}
func getFullUser(userExposeId: string): FullUser
```

joinするqueryについては、レコードが何らかのキーで一意になるので、関数の中で、そのチェックは行って、エラーを返すようにする。
ここではFullUserの取得する関数とかはそうだね。token取得系もそうかも。
- getCompanyUsers
- getFullUser

DBテーブルの型から、modelの方に変換する関数は必要。
FullCompany, FullUserからも同様に用意する。

テーブルとテーブルオブジェクトの関連はこんな感じになる
特定のテーブルを操作する際には、こんなルートで変換するイメージ

- user
  - FullUser
- user_email
  - FullUser  
    emailを持ってるので
- user_password
  - FullUser
  - password
- user_refresh_token
  - FullUser
  - refresh_token
- user_access_token
  - FullUser
  - refresh_token
- company
  - FullCompany
- user_company
  - FullCompany
  - FullUser
- company_invite
  - FullCompany
  - FullUser
- user_role
  - FullUser  
    CompanyもRoleも持ってるので
- role
  - Role

## 実装順序
この順番かな
- model
- DB model  
  変換ロジックも
- query
- http model  
  変換ロジックも
- procedure
- middleware
- router

それぞれ、utilityも定義した後かな。modelは無いけど
テスト書きながらボトムアップで進める。

