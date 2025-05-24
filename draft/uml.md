
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

```go
type CompanyInviteToken string
func generateCompanyInviteToken(): CompanyInviteToken
// tokenのverifyはequalityで良い
```

```go
type RefreshToken string
func generateRefreshToken(): RefreshToken
// tokenのverifyはequalityで良い
```

```go
type Password string
func getPassword(password: string): Password
func verifyPassword(password: Password): bool
// サンプルコードが多く、まー安全でデファクトな感じのbcrypt hashを使う
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


実装順だが、ほとんどコンストラクタなので、company, role, user, credentialsの順で実装する。
credentialsは、ライブラリ利用もあるので、後

いま、refresh tokenはuuid型にしてるけど、これはtoken型にしたほうがいいね。
access tokenもtoken型のほうがいいかも。中身はわからないけど、なにかしらのルールをもって生成された文字列。みたいな定義で。
access tokenは、Authenticに変換されるので、まーそこはそういうもんという感じで。

uuid型だと、ユーザ入力を結局stringで扱わないといけなくて、DBからきたuuid型の値と照合できない。
token型も、textと同様に一万文字までという制約だけはつけておこう。現実的にはもっと短いほうがいいが、jwtにいれる文字列の長さが予想できるかな。nameが255文字までのなのでbase64するとどんぐらいの長さか

## 認可
roleは用意してて、それに対して、どんなactionができるかは、auth自体で持たないといけない。
roleは共通なので、自由に作れるが、それ自体は権限と直接関連がなく、アプリケーション内で設定しないと行けないという方針。

主に、以下の操作かな
- 自分の情報の参照、操作
- company内の他人の情報の参照
- companyの参照
- companyの操作
- company invite
  これあれだな。事前にuserの情報みれちゃだめだから、userの参照は事前にできない。inviteテーブルからuser_idは取り去ろう

以下かなー
- none
  これは自分の操作のみで、そもそもcompanyに属して無い人ができること。特に値はないが、仮想的に書いとく
- employee
  company内で情報の参照ができる
- manager
  invite、companyの操作

上記の権限としては、以下の3つかな。
- company access
- company invite
- company edit

roleとしてはこんな付与の仕方
- none
- employee
  - company access
- manager
  - company access
  - company invite
  - company edit

role_id不要だね。他のアプリケーションでもrole labelで参照できるようにするのでauthも同様であるべき

role_permissionをすべてロードした状態のモデルを構築する感じにするか。
そいつが権限を判定するような感じのイメージ。それなら、role_labelとrole_permissionのrelationが不要なので、他のサービスでも実現が簡単
こいつは、事前にロードしておいたほうがいいので、これはioを伴うサービスにしたほうがいいか。
シングルトンで、最初にロードしたら、あとはそれを使う感じ。sql自体はそっちで参照する感じにはなるな。
goでシングルトンどうなるのかのイメージはつかないが。まーそれはキャッシュ機構を自分で用意するんかな。

jwtのやつもシングルトンにしないとだな。あとは矯正リロードフラグとか用意しとくか。いや、コンテナ入ってリロードはだるくて、バイナリに対してコマンドもうてないので不要かな

## error
エラー設計どうするか。
設定のエラーで立ち上がらないとかも、トップレベルに返したうえでpanicさせたいね。
panicはとにかくありえない状況でのハンドリングという感じで

大きくは
- 入力エラー
- 設定エラー

- model argument error
- user parameter error

どう違うとかまで表現したほうがいいかな。
結局型が合わないとか、そういう話ではあるが
stringなら、その値+message
構造体の場合は、原因になったpropertyを列挙する必要はあるわな

numberならrangeがあり、stringなら、長さやregexになるな。
関連なら、その関連というエラーを表現してもいいが、2つの関連ではなく、3つとかで絡んで来るとややこしいかもな。
変数の型もまちまちなので、そこの問題もある

型はわかってたほうがいいんだよな。
modelレベルでは、エラーは気軽に作るイメージでいいか。string自体はけっこうあるので、textの名前空間でtext errorみたいなの作る感じ。

で、modelを出たら、transfer/inの世界になるので、そこはユーザに優しいmessageを持つエラーにwrapする感じのイメージ
validation errorという命名でもいい。
ただ、これは結局は、modelのエラーを解釈して返すwrapperでしかないので、model errorみたいにしてmessageを親切にするだけの役割でよさそう

json unmarshal errorとかは、たぶん組み込みのやつが使えるんじゃないかな。わからんが

あとは、db系か
- record exist error
  - table name
  - key
- record not found error
  - table name
  - key
keyはstringにならしちゃっていいかな。サロゲートキーはreturnできないので、入力の値になるはず。
keyは、map[string]stringみたいな感じで、keyの名前と値を持つ感じにしておくと、どのkeyでエラーになったかがわかりやすい。

認証認可もあるな
- authentication error
  -> これはいわゆるログインエラーかな。password modelで使う
     passwordの文脈では、passwordが一致しないというエラーになるので、string errorが妥当なんだよな。
     まーでも認証エラー自体は実装して、どちらかというとprocedureで使う感じかな
- authorization error
- token expired error
  これは認可切れなので、認可エラー何だけど、はじめたsessionであれば、特定時間内ならexpiredなtokenでも使えるようにしておきたい。でないと1時間入力にかかるとかのときに使いづらくなる。
- authentication not exist error
  anonymousなuserができることは？みたいなことなので、これは認可エラーに含めていいんじゃないかな
  authenticがnilであっても引数に入れられる形にしておいて、nilだったら、anonymousなuserとして扱う感じでいいかも
  その状態で認可チェックに投げれば、必要な場面で認可エラーとできるはず。

- token invalid error
  これは、jwt名前空間でのエラーなので、jwt errorみたいな感じのほうがいいかも。項目+messageな感じかな

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

