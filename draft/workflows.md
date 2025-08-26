
# Work Flow
work flowとして流れを記載するが、WEB API ENDPOINTも記載する  

大まかに4つ
- user
- company
- user company role
- auth

共通のreturn形はcommonみたいにするか
正直、ワークフローでカテゴライズしてたほうがいいよな。
in/outで区別つけても、ちょっと分かりづらくなる気がする。

なので、transfer以下は、上記4つ+commonみたいにするか

in/outは同じファイルで問題ないと思う。そのほうがわかりやすいはず。file数も減るし。
struct名は、RequestとResponseというsuffixでいいかな。いずれにしろ、jsonのkey名もtag付けするので、基本的にはhttpでrequestとresponseという名前にしたほうがいい。

inは、methodでcoreを取り出す関数を定義。
outは、関数で、coreとかpremitiveな値を引数にとって、自身を返す。
dbとかは、引数がcore一つとかだったが、これはrelationが強い関連なので。outに関しては、割と関連の弱いものも、同じjsonに含めたりするので、引数に複数取る感じが良さそう。

ユーザ登録時にuserEmail, userPassword, userRefreshToken, userAccessTokenを全部作るが、保存前にちゃんと作りたい。
DBアクセス自体は、他データとの検証にアクセスはあるが、保存は最後にしたい。
そして、保存前に、それぞれのモデルのルール適用、つまりCreate関数内でのバリデーションをしたい。
ので、Unsavedな状態のものを作成してバリデーションを行えるようにすべき。
バリデーション自体は、あるかわからない状態だが、今後実装されたときに、事前にバリデーションできるのでは違うので

## ユーザ情報
### ユーザ登録
- /auth/user/register
- request
  - e mail address
  - password
  - name
  - bot or human
  - optional
    - 企業id
    - 企業 token
- response
  - refresh token
  - user
    - id
    - name
    - email
    - bot or not
- model
  - create user
    - email verify request
  - assign company
    - token check
  - login
    - password check
    - republish refresh token
    - expire old refresh token

in 
```go
type UserRegisterRequest struct {
    Email     string `json:"email"`
    Name      string `json:"name"`
    Bot       bool   `json:"bot"`
    Password  string `json:"password"`
}
```

out
```go
type UserRegisterResponse struct {
    RefreshToken string `json:"refresh_token"`
    User User `json:"user"`
}
```

### email varification
- /auth/user/verify_email
- request
  - e mail address
  - tempral token(メールで送られて来る)
- response
  - result
- model
  - email verify
    - token check

in
```go
type UserEmailVerifyRequest struct {
    Email string `json:"email"`
    Token string `json:"token"`
}
```

out
```go
type UserEmailVerifyResponse struct {
    User User `json:"user"`
}
```

### ユーザ参照
- /auth/user/self
- request
- response
  - user
    - id
    - name
    - ...
- model
  - get user

### name
- /auth/user/change
- request
  - id
  - name
- response
  - user
    - id
    - name
    - ...
- model
  - update user

### password
- /auth/user/change_password
- request
  - id
  - password
- response
  - user
    - id
    - name
    - ...
- model
  - update user password
    - expire old password

### メールアドレス変更
- /auth/user/change_email
- request
  - id
  - email address
- response
  - user
    - id
    - name
    - ...
- model
  - update user email
    - verify email request

### email verification
上と同様

## 企業作成
behaviorとしては
- create company
- get company
- get users of company
- get user of company
- get role
- publish token
- check token
- assign company role

### 企業作成
- /auth/company/create
- request
  - name
- response
  - company
    - name
    - id
    - master user
      - id
      - name
      - ...
- model
  - create company

### 企業参照
- /auth/company/<company_expose_id>
- response
  - company
    - name
    - id
    - master user
      - id
      - name
      - ...
- model
  - get company

### 企業ユーザ一覧
- /auth/company/企業id/user
- request
- response
  - users
    - user
- model
  - get users

## 企業アサイン
### ユーザ誘導
- /auth/company/企業id/invite
- request
  - ロールid
- response
  - token
- model
  - company invite user
    - publish token
    - expire old user token

単純なcreateな処理

### 企業登録
- /auth/company/企業id/accept
- request
  - token
  - (ログイン済み)
- model
  - assign company
    - token check

behaviorとしては、user,company,roleを受け取ってuser_company_roleを作成する感じ
その前段として、tokenのチェックを行うんだけど、これは別のbehaviorとして切り出す感じ
assign company自体が、複数のcompanyをassignできないようにしないとなので、そのルール自体は実装する

### ロール付与
- /auth/company/企業id/assign
- request
  - user id
  - role id
- response
  - user
  - role
- model
  - assign role

これもbehaviorがそのまま流用できそう
いや、これはそのcompanyにすでにassignされているかのチェックが必要
これは、別のbehaviorか？
その企業にassign済みかのチェックは別のbehaviorとして切り出す感じにして、assign自体はできそう。

以下になる？やな。
- getCompany
- getRole
- getUser
- checkUserAssignedToCompany
  - getUserでuser,company指定で取ってなければエラーとかでもいいかも
- assignRoleToUserInCompany

## ロール
### ロール一覧
- /auth/role
- response
  - roles
    - role
      - id
      - name
      - label
      - description
- model
  - get roles

## 認証
### login
- /auth/login
- request
  - id
  - password
- response
  - refresh token
  - access token
- model
  - login
    - check password
    - republish refresh token
    - expire old refresh token

### token refresh
- /auth/refresh
- request
  - user id
  - refresh token
- response
  - access token
- model
  - publish access token
    - expire check
      - sqlだけでチェックできて、有効なaccess tokenだけをDBから取り出せる。
    - republish access token
    - access tokenは2つまで発行でき、有効期限が切れていないものが2つある場合は、新しい方をreturn

