# エンドユーザー向けのMMPFインストール用zipファイル作成および解凍＆デプロイ手順

エンドユーザー向けの各サービスイメージとデプロイ用のファイル群（以下、インストール用zipファイル）の作成手順を示します。また、インストール用zipファイルの解凍およびMMPFのデプロイ手順を示します。

## 想定するデリバリーの流れ

| #   | 作業者     | 作業内容                                   |
| --- | ------- | -------------------------------------- |
| 1   | 開発者     | エンドユーザーの環境に合わせてIF等を改修し、各サービスイメージをビルドする |
| 2   | 開発者     | サービスイメージとデプロイに必要なファイル群をエンドユーザーに送付する   |
| 3   | エンドユーザー | サービスイメージを解凍し、MMPFをデプロイします            |

※開発者とエンドユーザーは、インターネットにつながる環境で作業することを前提とします。

※エンドユーザー向けのMMPFインストール用zipファイルで配布するMMPFは最小構成（サービス＋Redis）を想定しています。

## （開発者作業）インストール用zipファイルの作成手順

1. サービスとRedisのtarファイルを作成します。
下記コマンドを実行してtarファイルやmmpf-monolithicリポジトリをクローンし、zip化します。
    ```sh
    <!-- 権限付与 -->
    chmod -R 755 ./install
    <!-- ghcrへログイン -->
    make login
    <!-- 実行 -->
    ./install/create_module.sh
    ```
1. SLAMに必要なファイル群を格納します。

    ライブラリ・設定ファイルの用意を`/tmp/mmpf/mmpf-monolithic/lib`に格納してください。[ファイル配置(インターネット接続・未接続共通)](/README.md#ファイル配置インターネット接続・未接続共通)を参照下さい。

1. MMPF起動に必要な設定ファイルを編集します。

    devicesフォルダ配下のenvファイルを編集してください。[起動パラメータの編集(インターネット接続・未接続共通)](/README.md#起動パラーメータの編集インターネット接続・未接続共通)を参照下さい。

1. 必要であれば、ログレベルを変更します。

   [ログレベルの変更(インターネット接続・未接続共通)](/README.md#ログレベルの変更インターネット接続・未接続共通)

1. ファイル群をzipに圧縮してユーザーへ送付します。

   下記コマンドを実行してmmpf-service.tar、mmpf-other-services.tar、mmpf-monolithic、zip圧縮してください。
   ```sh
   $ ./create_module_zip.sh
   <!-- 以下のパス配下に.zipと.shが格納される為、そこからコピーする。 -->
   $ cd /app/delivery
   ```
   `mmpf_modules_YYYYMMDD.zip` と mmpf-monolithic/installディレクトリ直下の`unzip_module.sh`をエンドユーザーへ送付します。

## （エンドユーザー様作業）インストール用zipファイルの解凍・MMPFの起動手順
以下手順は、`Ubuntu20.04`を想定しています。

1. ユーザーを作成します。
  
   以下の手順は、ユーザー名`mmpf`が作成済みであることを前提とします。未作成の場合は、ユーザーを作成し、ログインしてください。
   ```sh
   $ useradd -m mmpf
   $ passwd mmpf <任意のパスワード>
   $ su mmpf
   ```

1. 必要なツール群をインストールします。(※要ネットワーク接続)
  
   unzip、docker、docker-compose(ver1.29以降)が未インストールの場合は、以下コマンドを実行します。これにより、必要なツールをインストールし、`docker`コマンドを`mmpf`ユーザ権限で使用できるようにします。
   ```sh
   $ sudo apt-get update
   $ sudo apt-get install zip
   $ sudo apt-get install docker
   $ sudo apt-get install docker.io
   $ sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   $ sudo chmod +x /usr/local/bin/docker-compose
   $ sudo gpasswd -a mmpf docker
   $ sudo chmod 666 /var/run/docker.sock
   ```


1. zipを解凍します。

   `mmpf_modules_YYYYMMDD.zip` と `unzip_module.sh`をサーバーの `/tmp` にコピーし、ターミナルで下記コマンドを実行してください。
   ```sh
   $ chmod 755 /tmp/unzip_module.sh /tmp/mmpf_modules_YYYYMMDD.zip
   $ /tmp/unzip_module.sh /tmp/mmpf_modules_YYYYMMDD.zip
   $ cd /home/mmpf/mmpf_modules/mmpf-monolithic
   $ ./install/load_modules.sh
   ```

1. MMPFを起動します。

   ターミナルで下記コマンドを実行して、MMPFを起動してください。
   ```sh
   $ cd /home/mmpf/mmpf_modules/mmpf-monolithic
   $ ./mmpf.sh deploy -mmpf
   ```

   停止する場合は、下記コマンドを実行してください。
   ```sh
   $ cd /home/mmpf/mmpf_modules/mmpf-monolithic
   $ ./mmpf.sh undeploy
   ```

   その他のコマンドについては、、[シェルコマンド実行手順](/README.md#シェルコマンド実行手順)に記載しています。
