
# design

## summary
### requirement
エンタープライズなアプリケーションにおいて、認証機能を単独で管理できるようにしておきたい。  
企業は人に利用許可を与えることができるが、アカウントは人のものであるべき  
人は物理的な人とアカウントは1対1で紐づけ、複数アカウント利用によるセキュリティの抜け道をなるべく減らしておきたい  
登録容易のため、使用者が認識するキーはメールアドレスとしたいが、サロゲートキーも利用できるようにしておく。  
独立性を高めるため、認可は行わないが、認可に一貫性を持たせるために、利用者の役職、役割などを示すタグは管理する  
他アプリケーションから認証を呼び出すことがなるべく少なく、スケールできるようにしておく  

### thinking
認証の機能  
認可の機能は提供しないが、認可に一貫性を持たせるための支援機能は提供する  
権限管理する企業と、使用者である人がいる。  

### physical
アプリケーションサーバとrdbを使う。  
実際に運用するわけではなく、面倒なのでフロントは書かない。  
実際に運用するわけではなく、面倒なのでin memory dbは使わない。  
単独のurlドメインで動作することを想定する  

## データ定義
ざっくり  

- ユーザ
  - id
  - name
  - email address
  - bot or human
  - credencials
    - password
    - refresh token
    - access token
  - role id
- 企業
  - id
  - user id
  - name
- 権限
  - id
  - company id
  - name

## Work Flow
work flowとして流れを記載するが、WEB API ENDPOINTも記載する  

### ユーザ登録
#### ユーザ登録
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

#### email varification
- /auth/user/verify_email
- request
  - e mail address
  - tempral token(メールで送られて来る)
- response
  - result
- model
  - email verify
    - token check

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

### 企業アサイン
#### ユーザ誘導
- /auth/company/企業id/invite
- request
  - optional
    - user id
    - ロールid
- response
  - token
- model
  - company invite user
    - publish token
    - expire old user token

#### 企業登録
- /auth/user/assign/company/企業id
- request
  - token
  - (ログイン済み)
- model
  - assign company
    - token check

### ロール作成
#### 作成
- /auth/company/企業id/role/create
- request
  - name
  - label
  - description
- response
  - role
    - id
    - name
    - label
    - description
- model
  - create role

#### 削除
- /auth/company/企業id/role/delete
- request
  - id
- response
  - role
    - id
    - name
    - label
    - description
- model
  - role delete

### ロール付与
- /auth/company/企業id/role/assign
- request
  - user id
  - role id
- response
  - user
  - role
- model
  - assign role

### ユーザ情報変更
#### name
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

#### password
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
#### メールアドレス変更
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

#### email verification
上と同様

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
    - republish access token
    - access tokenは2つまで発行でき、有効期限が切れていないものが2つある場合は、新しい方をreturn

### その他

## Rdb Schema
id はserialで自動発行連番  
expose_idはアプリケーションで発行するuuidか、uuidは長すぎるので数桁のcodeとか  
6桁のアルファベットなら3億パターンとかなのでバッティングしづらいはず。あとはdbチェックしてという感じで  
expose idは頭にuser_,company_,role_をつけるか

refresh tokenはuuid  
passwordは100文字ぐらいまで  

```sql
create table user (
  user_id serial
  expose_id unique
  expose_email_id
  name
  bot_flag
  registered_date
);

create table user_email (
  user_id
  user_email_id
  email
  verified
  registered_date
  expire_date
);

create table user_password (
  user_id
  password
  registered_date
  expire_date
);

create table user_refresh_token (
  user_id
  refresh_token
  registered_date
  expire_date
);

create table user_access_token (
  user_id
  access_token
  registered_date
  expire_date
);

create table company (
  company_id
  expose_id
  name
  master_user_id
  registered_date
);

create table user_company (
  company_id
  user_id
  registered_date
);

create table company_invite (
  company_id
  token
  user_id nullable
  registered_date
  expire_date
);

create table user_role (
  user_id
  role_id
  registered_date
  expire_date
);

create table role (
  company_id
  role_id
  name
  label unique
  description
);
```

## Application Modeling
Rdb Schemaに習う感じだが、そっちで表現しきれないものを記載していく  

- user
  - user
    - create user
      - email verify request
    - assign company
      - token check
    - update user
  - email
    - email verify
      - token check
    - update user email
      - verify email request
  - password
    - login
      - password check
      - republish refresh token
      - expire old refresh token
    - update user password
      - expire old password
  - refresh_token
    - publish access token
      - expire check
      - republish access token
      - access tokenは2つまで発行でき、有効期限が切れていないものが2つある場合は、新しい方をreturn
  - access token
    - verify access token
- company
  - company
    - create company
    - company invite user
      - publish token
      - expire old user token
  - role
    - assign role
    - create role
    - role delete

