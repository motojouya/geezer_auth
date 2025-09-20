
# GEEZER_AUTH

## GEEZER?
なんかいい感じにエンタープライズないい感じな経験則なあれを作ってあれするプロジェクト的ななにか。  

## GEEZER_AUTH?
エンタープライズな経験則なあれはあんま公開するとあれだが、認証ぐらいは公開してもいけるっしょ。  

## 開発マナー
- 設計する  
- Go言語やる  
- PostgreSQLする  
- AI使う  
- テスト書く  

### 設計
- ディレクトリはgo標準に習って必要なものを使うが、初期設計だけdraftディレクトリを利用する  
- 初期設計はdraftディレクトリで行う。雑多な情報が入るので、継続的に読まれることを検討しない。  
- docsには使われ方や各ドキュメントへのindexなど。モジュール単位の仕様はプログラムファイルのコメントに記載。  
- モデリングについては、大枠では説明するし図が必要ならdocsだが、テキストのみならプログラムファイル側のほうに比重を置く。  

### テスト
- 単体テストは基本参照透過な対してstate lessな感じでやる  
- RDBアクセスのテストは特定のモジュールに押し込めるのでそれも単体テストで行う  
- もう一段階上のレイヤではAPIエンドポイントに対してRDBの動き含めた結合テスト  
- 最後に特定の機能のワークフロー全体をテストする統合テスト  
- 前2つは対象モジュール内のtestファイル、後2つはtestディレクトリ  

### 開発環境
まだわからない  


## 環境構築

install golang-migrate
https://github.com/golang-migrate/migrate

```shell
make dockerup
make migrate
make postgres
```

```sql
insert into role (label     ,name    ,description,register_date)
          values ('EMPLOYEE','作業者','作業者'   ,now()        )
               , ('MANAGER' ,'管理者','管理者'   ,now()        )
               ;
insert into role_permission (role_label,self_edit,company_access,company_invite,company_edit,priority)
                     values ('EMPLOYEE',true     ,true          ,false         ,false       ,5       )
                          , ('MANAGER' ,true     ,true          ,true          ,true        ,9       )
                          ;
```

docker稼働確認
```shell
curl -X POST -H "Content-Type: application/json" -d '{"Email":"test@example.com","Name":"Test User","Bot":false,"Password":"password123"}' http://localhost:${port}/user/register
```

test
```shell
make runt port=${port}
```

## やり残し
- UserAuthenticを集約ルートとすべき関数が中途半端になっている
  - 具体的にはuser_company_roleの生成時に、UserAuthenticを使っていないので、UserAuthenticが持つビジネスルールの検証ができていない
  - roleAssigner,userGetterがbehavior/companyにあるのはよくない。これはuserの振る舞い(controlのレイヤではcompanyでいい)
- 認可ロジックで、アクセスするcompanyがどのcompanyか見てないので、roleさえあれば、どんなcompanyでもアクセスできてしまう
- user,companyの削除、あるいは無効化ロジックがなく、データモデル上も考慮されていない
- persist_keyはuuid v7にすればdbで発行しなくていいし、register_dateとかも不要になるのでそうしたい。
- CIとか。testもmakeからだとdbテストを分離できてるが、shellでするので、ciでも同様に書く必要あり
- runnの実装をhostからできるようにしてるけど、docker containerからもできるようにしたい




