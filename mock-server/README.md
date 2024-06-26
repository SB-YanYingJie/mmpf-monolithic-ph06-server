# mock-server
AWS環境で使用するための、モックサーバーの実装やデプロイに必要な資源を格納する。
## ディレクトリ構成
directory | 内容
---|---
`/apitest` | モックサーバの動作検証用のソースを格納
`/assets`       | モックサーバーへ送信する画像を格納
`/cert`       | crtファイル、keyファイルを格納
`/cmd`   | モックサーバーの実装に関するファイルを格納
`/deployment` | mock-serverのデプロイに関するファイルを格納
`/mock` |　モッククライアントの実装ファイルを格納

## awscliを用いたデプロイ手順

### 準備
#### 0. ghcrにloginする
あらかじめgithub上で個人用access tokenを作成し、env/devcontainer.envに記載したのち、以下のコマンドを実行する。
```sh
$ make login
```
#### 1. ghcrログイン情報をSecrets Managerに登録する
**本手順は開発者1名が実施すれば十分。**

ghcrに登録されたイメージをECSがダウンロードするための認証情報として、先の手順で作成したトークンをユーザ名と併せてAWSのSecrets Managerに登録する。
```
username=<your username>
password=<your access token>
```

#### 2. デプロイに必要なファイルを用意する
- Dockerfile
- task.json
    - タスク定義ファイル
    - コンテナの環境変数に機微情報は書いてはならず、必ずSecrets Managerに登録する
- service.json
    - タスクを実行するサービスの設定ファイル
- go.mod

### 手順
#### 1. docker imageをビルドし、リポジトリにプッシュする
```sh
$ make build_image_ssh
$ # or
$ make build_image_https

$ make push_image tag=v1
```
ここで指定したタグを、`deployment/definition/task.json`の`"image"`属性に反映する。

なお通常のbuildコマンドはキャッシュを利用する設定としているが、変更反映時などでキャッシュの動きが怪しい場合は `docker build --no-cache ...` を利用するとよい。

#### 2. ECSのクラスター、タスク定義、サービスを作成/更新
```sh
# 1. クラスターを作成する（作成済みの場合スキップ）
$ make ecs_create_cluster env=dev

# 2. タスク定義を登録する
$ make ecs_register_task env=dev
```
2の結果として出力された`"revision"`の値を、`deployment/definition/service.json`の`"taskDefinition"`属性に反映する。

```sh
# 3. サービスを作成/更新する
# --- 初回のデプロイの場合
$ make ecs_create_service env=dev
# --- 更新の場合
$ make ecs_update_service env=dev
```

最後に、イメージのタグやタスク定義のリビジョン、その他設定値が変わった場合は差分をコミットする。

## モッククライアントを用いた動作確認手順
1. goの実行環境が必要なので、devcontainerを起動する。既に起動済みであれば、次の手順へ進む。
1. アクセス対象先にあわせて、`./mock-server/apitest/api_test.go`内の`address`部を変更する。
    - ローカルで起動しているサーバに対し、検証したい場合
      - ```address = "localhost"```
    - 特定のIPを持ったサーバに対して、検証したい場合
      - ```address = "xxx.xxx.xxx.xxx"(ECS タスクのパブリックIP,MMPFを起動している別PCのサーバのIP　等々)```
1. 次に、モッククライアントを用いて、サーバに対してリクエストする。サーバで待ち受けているポート番号に併せてdevcontainer内で以下コマンドを実行する。
    ```
    make run_mock_request PORT=xxxx
    ```
