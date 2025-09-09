
# design

## summary
### requirement
エンタープライズなアプリケーションにおいて、認証機能を単独で管理できるようにしておきたい。  
企業は人に利用許可を与えることができるが、アカウントは人のものであるべき  
人は物理的な人とアカウントは1対1で紐づけ、複数アカウント利用によるセキュリティの抜け道をなるべく減らしておきたい  
登録容易のため、使用者が認識するキーはメールアドレスとしたいが、サロゲートキーも利用できるようにしておく。  
独立性を高めるため、認可は行わないが、認可に一貫性を持たせるために、利用者の役職、役割などを示すタグは管理する  
他アプリケーションから認証を呼び出すことがなるべく少なく、スケールできるようにしておく  

### consideration
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
#### 企業作成
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

#### 企業参照
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
#### ロール一覧
- /auth/company/企業id/role
- response
  - roles
    - role
      - id
      - name
      - label
      - description
- model
  - get roles

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

#### ユーザ一覧
- /auth/company/企業id/user
- request
- response
  - users
    - user
- model
  - get users

#### ロール付与
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
#### ユーザ参照
- /auth/user/self
- request
- response
  - user
    - id
    - name
    - ...
- model
  - get user

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
#### ユーザ参照
上と同様

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
      - sqlだけでチェックできて、有効なaccess tokenだけをDBから取り出せる。
    - republish access token
    - access tokenは2つまで発行でき、有効期限が切れていないものが2つある場合は、新しい方をreturn

### その他

## Rdb Schema
id はserialで自動発行連番  

```sql
create table user (
  user_id serial
  expose_id unique
  expose_email_id
  name
  bot_flag
  registered_date
  updated_date
);

create table user_email (
  user_email_id
  user_id
  email
  verify_token
  register_date
  verify_date
  expire_date
);

create table user_password (
  user_password_id
  user_id
  password
  registered_date
  expire_date
);

create table user_refresh_token (
  user_refresh_token_id
  user_id
  refresh_token
  registered_date
  expire_date
);

create table user_access_token (
  user_access_token_id
  user_id
  access_token
  source_updated_date -- userの情報が更新されたときに、access tokenを再発行するためのフラグ。user.updated_dateの値をコピーして入れておく
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

create table company_invite (
  company_id
  token
  role_id
  registered_date
  expire_date
);

create table user_company_role (
  user_company_role_id
  user_id
  company_id
  role_id
  registered_date
  expire_date
);

create table role (
  label unique
  name
  description
  registered_date
);

create table role_permission (
  role_label
  self_edit bool
  company_access bool
  company_invite bool
  company_edit bool
  primary uint
);

insert into role (label     ,name    ,description,register_date)
          values ('EMPLOYEE','作業者','作業者'   ,now()          )
               , ('MANAGER' ,'管理者','管理者'   ,now()          )
               ;
insert into role_permission (role_label,self_edit,company_access,company_invite,company_edit,priority)
                     values ('EMPLOYEE',true     ,true          ,false         ,false       ,5      )
                          , ('MANAGER' ,true     ,true          ,true          ,true        ,9      )
                          ;
-- defaultの値は、modelが知っているので、ここでは設定しない。
```

```sql
insert into role (label ,name ,description,register_date) values ('EMPLOYEE','作業者','作業者' ,now() ) , ('MANAGER' ,'管理者','管理者' ,now() );
insert into role_permission (role_label ,self_edit,company_access,company_invite,company_edit,priority) values ('EMPLOYEE',true ,true ,false ,false ,5 ) , ('MANAGER' ,true ,true ,true ,true ,9 ) ;
```


## Application Modeling
Rdb Schemaに習う感じだが、そっちで表現しきれないものを記載していく  

expose_idはアプリケーションで発行するuuidか、uuidは長すぎるので数桁のcodeとか  
6桁のアルファベットなら3億パターンとかなのでバッティングしづらいはず。あとはdbチェックしてという感じで  
expose idは頭にuser_,company_,role_をつけるか  
そう考えると、role expose idはcompany単位で一意としたほうが良さそう  

refresh tokenはuuid  
passwordは100文字ぐらいまで  
access tokenはjwt

userはidでもemailでもログインできる。最初に登録したemailのみ。後にemail自体は変更できるが、ログイン時に使うemailは変更できずそのまま。
idは自動発行。emailの前者はemail idと呼ぶ。

emailはverifyされて初めて使えるようになる。つまり、verifiedでnot expiredなものが有効だし、それは1レコードに保つようにロジックを組む。

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

### role
とりあえず3つ用意する
- unauthorized  
  ユーザ未登録かemailが未承認状態のもの  
  タグなので、そういったユーザに権限をつけたい場合はこちらをいじることもできるが、意味としては認証がないユーザ
- none  
  emailが承認され、会社にアサインされると付与されるもの  
  いわゆるdefault権限  
  自分個人の情報は操作できるが、それはそもそもroleを見ない仕組みとしたい  
- administrator  
  企業を作ると作ったユーザは自動的にこれになる  
  ロールの作成、付与、ユーザのinviteなどができる  

## access token
jwtだが、そのモデリングが必要
access tokenは2つまで発行できる。これは、一つ前のものがexpireしそうなときに、事前に取得しておけるようにするため。
とはいえ、アプリケーション側で、ミッションクリティカルな処理では、expireの値で一元的にと切るのではなく、sessionをexpandするなどの処理も必要。

ただし、user情報が更新された際には、access tokenは再発行できる。passwordやrefresh tokenはuser情報なので、更新されても変わらない。
このuser情報の一意性に対して、有効なaccess tokenは2つまでとする。

user
  - user expose id
  - user email for id
  - user email
  - user name
  - role label
  - company expose id
  - company name
  - bot flag
token
  - token expire date

## アーキテクチャ
[ソフトウェア設計のマイブーム](https://zenn.dev/motojouya/articles/software_design_my_boom)の記事に従う。
webだけを想定されたものではないため言及がないが、url routingは塊を一つのファイルに纏めて、コメントでworkflowなどを補って書きたい。  
基本的にinbound handlingの部分だが、入力値を、構造体に変換する機能しか持たないので、宣言的、あるいは設定的に書くことできるようにする  

## ライブラリ
### web
これは標準ライブラリが優秀なので、それを利用する。  
middlewareの実装もtransducerを用意するだけなので、grue codeとしてそういったものを用意して利用する。  
パラメータのハンドリングやresponseのjsonのハンドリングは、packageとして公開するので説明はそちらに譲る

### db
dbは、gorp+goquで行う。  
primary keyでの出し入れはモデリングを伴ったほうが楽だし、scanをそのまま使うのは大変なのでormapper likeなgorpがちょうどいい。  
queryはquery builderで組み立てたいのでgoquで行う。リレーションという概念がgorpに無いので、常にquery builderでjoin句を書くわけだが、よりsqlの構文に近くなって、わかりやすくなることを期待している。  

リレーションは、単純なテーブル同士ならワークするが、サマリテーブルのモデリングとの関連や、複数のキーでjoinパターンが違う場合等、様々考えられるので、そこは愚直にコーディングするイメージ。
なので、gorp+goquがちょうどいい。また、この構成を助けるためのhelper関数も用意して公開する。  

## package
### access token middleware
この認証サーバを利用するアプリケーションでは、access tokenからユーザ情報を取得する。  
jwt tokenなので、これを分解してユーザ情報をハンドリングするwebのmiddleware moduleを提供する

機能としては、jwt tokenを分解してuserオブジェクトを作る機能となる。jwt tokenがあってもなくても、認可自体は各種処理に記載するので、userをラップしたオブジェクトを用意し、認証があったかユーザがあったかを、各処理に渡せるようにする。  
あくまで認可は行わず、下に情報を渡すだけの機能とする。  

### webのパラメータハンドリング
goは標準ライブラリが優秀なのでそれを利用したいが、入力パラメータのハンドリング、json変換はモジュールを用意してやるほうがいい。  
jsonをparse,serializeする標準ライブラリはあるが、query string, form, file dateを宣言的に記載して自動的にハンドリングできるようにはなっていない。  
jsonのtagのような感じで、formなのかbody jsonなのかquery stringなのかを示しつつ、それらの構造体の中で、propertyを定義してparseできるようにしたい。  

returnされるjsonのserializeも同じモジュールで実現したい。  

### dbアクセス
gorp+goquを利用するつもりだが、helper関数もいくつか用意しておく。これも公開する。  

