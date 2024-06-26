# mmpf-monolithic
Machinemap Project 2023-ph5用のリポジトリです。
## 目次
- [開発準備](#開発準備)
- [ディレクトリ構成](#ディレクトリ構成)
- [ビルド方法](#ビルド方法)
- [デプロイ方法](#デプロイ方法)
- [モッククライアントを用いた動作確認方法](#モッククライアントを用いた動作確認手順)
- [クライアントから受け取った画像を表示する方法](#クライアントから受け取った画像を表示する方法)
## 開発準備
- 本プロジェクトでは、docker imageをghcr（github container repository）に保存します。
ghcrからイメージをpush/pullするためにあらかじめgithub tokenを作成する必要があります。[リンク先](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)を参考に、token(classic)を作成してください。

  tokenには必ず以下の権限を付与してください。

  ```
  write:packages（Upload packages to GitHub Package Registry）
  ```
- [machinemapplatform](https://github.com/machinemapplatform) へのアクセス権限が付与されているかを確認して下さい。アクセスできない場合は、管理者に連絡し、リポジトリへのアクセス権限付与を依頼してください。
machinemapplatform管理者にアクセス権限の付与を依頼してください。
- [KudanJP/KdSlamGo](https://github.com/KudanJP/KdSlamGo)に対してアクセスできるか確認して下さい。アクセスできない場合は、管理者に連絡し、リポジトリへのアクセス権限付与を依頼して下さい。
- `env/devcontainer.env.sample`をコピーして、`devcontainer.env`を作成し、各種必要な環境変数を埋めてください。devcontainer起動時に必要になるのは、以下の環境変数のみです。それ以外はデフォルト値で構いません。
  - `GITHUB_USERNAME`
  - `GITHUB_ACCESS_TOKEN`
## ディレクトリ構成
directory | 内容
---|---
`/.devcontainer` | VSCode Remote - devcontainerの設定ファイルを格納
`/.github`       | 各フォーマットやGitHub Actionsの設定を格納
`/.vscode`       | vscodeの設定を格納
`/aws`           | aws関連のファイルを格納
`/cmd`           | golangのメインプログラム main.go および環境変数取得用の config.go を格納
`/devices` | クライアント毎に設定するサーバーの環境変数設定ファイルを格納
`/env` | 環境変数の定義ファイルを格納
`/install` | インターネット未接続環境でのデプロイ方法に関するファイルを格納
`/internal` | ロジックの本体を格納
`/lib` | KdSlamの実行に必要なライブラリを格納
`/mock-server`   | mock-serverのソースやデプロイに関するファイルを格納
`/ops-utils` | コンテナに関するメトリクスを計測するためのプログラムを格納
`/pkg` | util系のプログラムを格納
`/scripts` | サーバ起動に必要なスクリプト群を格納
`/zap` | ログレベルの設定ファイルを格納

## devcontainer起動方法
本手順を実行する前に、[開発準備](#開発準備)が実施済であることを前提とします。
1. VSCodeで`mmpf-monolithic`フォルダーを開きます。既に開いている場合は次の手順へ進みます。
2. 右下に表示される、`Reopen in Container`を選択します。もしくは、VSCode左下の緑色のマークを押下し、`Reopen in Container`を選択します。
3. 起動が完了するまで待ちます。右下の`Starting Dev Container(show log)`をクリックすると、状況が確認できます。

## ビルド方法
本手順を実行する前に、[devcontainer起動方法](#devcontainer起動方法)の手順が実施済であることを前提とします。
1. devcontainerを起動します。既に起動済みであれば次の手順へ進みます。
2. `リポジトリルート階層上`で、以下コマンドを実行し、imageをbuildし、ghcrへpushします。(※既に最新版がpush済みであれば、以下コマンドを実行する必要はありません。)
    ```sh
    $ make login
    $ make build_image
    $ make push_image
    ```
## デプロイ方法
本手順を実行する前に、[ビルド](#ビルド方法)が実施済みであることを前提とします。
ここでは、MMPFサーバを起動する手順のみを記載しています。Viewerの起動方法に関しては、[Viewerの起動手順](https://github.com/machinemapplatform/viewer#mmviewer%E6%93%8D%E4%BD%9C%E6%96%B9%E6%B3%95)を参照下さい。

### インターネット未接続環境の手順
 [エンドユーザー向けのMMPFインストール用zipファイル作成および解凍&デプロイ手順](./install/README.md)を参照下さい。
### インターネット接続環境の手順
以下手順に沿って進めてください。なお、ここからはdevcontainerは使用しません。
### デプロイ前の準備
#### インストール
1. MMPFサーバ起動に必要なミドルウェアをインストールしておきます。インストール済みであれば不要です。

    | ミドルウェア | version |
    | - | - |
    | docker | 20.10.22 |
    | docker-compose | 1.29.0 |  

1. MMPFサーバの実行環境に、最新のブランチの内容をcloneもしくは、pullしておきます。最新化済みであれば、不要です。
#### ファイル配置(インターネット接続・未接続共通)

デプロイ前にSLAMに必要なファイルを格納します。各ファイルの概要は以下の通りです。基本的には、MAPファイルとキャリブレーションファイルのみを`リポジトリルート階層/lib`配下へ格納します。

| ファイルの種類 | 概要 |
| - | - |
| MAPファイル       | kdslam用のmapファイル。拡張子はkdmp                     |
| キャリブレーションファイル | クライアントごとに作成される、kdslam用のキャリブレーションファイル。拡張子はini | 

※ライブラリの実行ファイルやボキャブラリーファイルは、docker image内で管理しており、特に更新がなければ改めて格納する必要はありません。更新する場合は、[build-push方法](https://github.com/machinemapplatform/base-image)を参照下さい。

#### 起動パラーメータの編集(インターネット接続・未接続共通)

MMPFサーバの起動パラメータをenvファイルを用いて設定します。envファイルには、MMPFサーバ共通の設定値(`devices/common.env`)と、MMPFを利用する端末毎の設定値(`devices/.env.xxxx`)があり、これらを作成・編集することで起動パラメータを決定します。各起動パラーメータの説明や設定方法は[envファイル設定方法](./devices/README.md#devices)を参照下さい。
#### ログレベルの変更(インターネット接続・未接続共通)

`リポジトリルート階層/zap/prd/config.yml`の`level`の設定値を変更することで、アプリケーションログのレベルを変更できます。`debug`レベルだと、ログの出力数が多く、パフォーマンスが低下する恐れがあるため、デフォルトは`info`としています。必要であれば、変更してください。
  - debug -> アプリケーションに仕込まれた全てのログを表示します。
  - info -> ErrorログやWarnログのみを表示します。

### デプロイ手順
#### シェルコマンド実行手順
リポジトリルート階層にある`mmpf.sh`を使用し、デプロイを実施します。
`mmpf.sh`や`mmpf.sh`が実行するシェルスクリプトに対して、以下コマンドを実行し、権限を付与してください。
```sh
$ chmod 755 ./mmpf.sh 
$ chmod -R 755 ./scripts ./ops-utils
```
それぞれのユースケースに併せて、以下のいずれかのコマンドを実行する。(※コンテナ内からコマンドを実行しないでください。devcontainerを起動している場合は、devcontainer外のターミナルでコマンドを実行します。)
- MMPFサーバのみを起動する場合
  ```sh
  $ ./mmpf.sh deploy -mmpf
  ```
- MMPFサーバを起動し、コンテナのメトリクスを計測する場合(※要注意)
 ※`リポジトリルート階層/ops-utils/sh/output.csv`が既にある場合は、本コマンド実行時に削除されてしまうため、事前に退避させておく必要があります。
  ```sh
  $ ./mmpf.sh deploy -mmpf-with-measurement
  ```
- MMPFサーバを停止する場合
 ※コンテナのメトリクスを計測している場合、その処理も停止されます。
  ```sh
  $ ./mmpf.sh undeploy
  ```
#### メトリクスについて
サーバを起動しているコンテナ毎のCPU使用率やメモリ使用率等を計測することができます。
- リアルタイムで動作しているコンテナのメトリクスを閲覧する場合は、以下コマンドを実行します。
  ```sh
  $ docker stats
    <!-- 実行結果 -->
    CONTAINER ID   NAME                                 CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O   PIDS
    d9a846baa4ee   angry_varahamihira                   19.79%    27.17MiB / 15.51GiB   0.17%     512MB / 432MB     0B / 0B     55
    9c8790a405e6   suspicious_goldstine                 0.00%     1.379MiB / 15.51GiB   0.01%     2.73kB / 858B     0B / 0B     2
    793ad5579605   docker                               1.72%     56.09MiB / 15.51GiB   0.35%     432MB / 512MB     0B / 0B     42
    85cbbca87c28   mmpf-monolithic_devcontainer-app-1   4.48%     1.899GiB / 15.51GiB   12.25%    6.22MB / 6.07MB   0B / 0B     292
  ```
- [シェルコマンド実行手順](#シェルコマンド実行手順)に記載の方法で収集したコンテナメトリクスをグラフ加工する場合は、[コンテナ毎のCPU使用率、メモリ使用率のメトリクスをグラフ化する手順](./ops-utils/docs/README.md)を参照下さい。

## モッククライアントを用いた動作確認手順

手元にクライアント端末がなかったり、特定のアクセス先への疎通を確認したい場合等にモッククライアントを用いて、動作確認を行なうことができます。確認手順に関しては、[モッククライアントを用いた動作確認手順](./mock-server/README.md#モッククライアントを用いた動作確認手順)を参照下さい。

## クライアントから受け取った画像を表示する方法

MMPFの実行環境で環境変数とXwindowの設定をすることで、クライアントから受け取った画像を指定したディスプレイ上に表示することができます。以下の手順において、画像を表示する端末はWindowsを想定しています。


#### 環境変数の設定

|環境変数|設定するためのコマンド|
|-|-|
|DEV_DISPLAY_RAW_IMAGE|`export DEV_DISPLAY_RAW_IMAGE=true`|
|DEV_DISPLAY_DEBUG_IMAGE|`export DEV_DISPLAY_DEBUG_IMAGE=true`|
|DISPLAY（devcontainerの場合）|`export DISPLAY=host.docker.internal:0`|
|DISPLAY (特定IPに対し、表示したい場合) |`export DISPLAY=xxx.xxx.xxx.xxx:0 (例 192.168.11.100:0)`|

#### Xwindowの設定

X Server の事前準備が必要です。 以下に準備の手順を示します。

1. X Server をインストールします
   1. https://sourceforge.net/projects/vcxsrv/ からダウンロードしてインストールします。いくつか選択肢が出ますが、デフォルトでOK
2. X Server の初期設定をする
   1. https://rin-ka.net/windows-x-server/ の「VcXsrv（Xサーバー）の初期設定」を参考に初期設定します
3. x11-appsを実行環境にインストールします (devcontainerの場合は不要)
  `sudo apt install x11-apps`
4. 実行環境の環境変数を設定します (devcontainerの場合は不要)
   上記表の環境変数値に基づいて設定します
5. 以下コマンドで、X Server を起動します
  `xeyes`
6. 「目」のウィンドウが立ち上がればOK（プログラムから X Serverを起動すると画像が表示される）
