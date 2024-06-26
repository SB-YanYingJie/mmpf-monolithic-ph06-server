# EC2インスタンスでモノリシック版MMPFを起動する手順

## 前提
本手順は事前に以下の事項が揃っていることを前提とする。
- デプロイ対象のAWS環境にマネジメントコンソールからアクセスできる
- ローカルPCのSSHログイン用の設定が整っている
- mmpf-monolithicのGitHubリポジトリ（github.com/machinemapplatform/mmpf-monolithic ）にアクセスできる
- GHCRからイメージをpullするためのアクセストークンを作成済みである
  - 未作成の場合は https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token を参考に作成しておく
    - アクセストークンには以下の権限を付与する
      ```
      read:packages（Download packages from GitHub Package Registry）
      ```

## 手順

### EC2インスタンスの起動

AWSコンソールで EC2のダッシュボードを開く。
  - 画面右上で、「東京 (ap-northeast-1)」がリージョンに選択されていることを確認する

「インスタンスを起動」をクリックし、続けて表示されるプルダウンメニューでも「インスタンスを起動」をクリックする。

インスタンス起動画面で各入力欄に以下の値を入力する。（各値は例）

| 分類 | 項目名 | 値 | 備考 |
| -----  | ---- | -- | ---- |
| 名前とタグ | 名前 | mmpf-monolithic-server | - |
| アプリケーションおよびOSイメージ（Amazon マシンイメージ） | Amazon マシンイメージ（AMI） | Ubuntu Server 20.04 LTS (HVM), SSD Volume Type | - |
|^ | Architecture | 64 ビット (x86) | 初期値のまま |
| インスタンスタイプ | インスタンスタイプ | m5.large | - |
| キーペア（ログイン） | キーペア名 | mmpf-monolithic-server-key | 適切なキーペアが無ければ「新しいキーペアの作成」から新規作成すること |
| ネットワーク設定 | VPC | [EC2インスタンス設置対象のVPC。無ければデフォルトのままでも可] | - |
|^ | サブネット | [EC2インスタンス設置対象のパブリックサブネット。無ければデフォルトのままでも可] | - |
|^ | パブリックIPの自動割り当て | 無効化 | 後でElastic IPを紐づけるため |
|^ | ファイアウォール（セキュリティグループ） | セキュリティグループを作成する | - |
|^ | セキュリティグループ名 | mmpf-monolithic-server-sg | - |
|^ | 説明 | mmpf-monolithic-server-sg created 2023-02-22T12:13:49.766Z | 時刻は一例 |
|^ | インバウンドセキュリティグループのルール |> | [※後述](#インバウンドセキュリティグループのルール) |
| ストレージを設定 | - | 1x 16GiB gp2 | 初期値の8GiBだとGHCRからイメージをpullする際に容量が足らなくなる |
| 高度な詳細 | - | - | 全て初期値のまま |

「インスタンスを起動」をクリックしてインスタンスを起動する。
  - 次の画面で「インスタンスの起動を正常に開始しました」と表示されることを確認する

#### インバウンドセキュリティグループのルール
デフォルトで設定されている「セキュリティグループルール１」は削除し、以下の三つの通信を許可するセキュリティグループルールを追加する。

| タイプ | プロトコル | ポート範囲  | ソースタイプ | ソース | 説明 |
| -----  | ---------- | ---------- | ------------ | ------ | ---- |
| ssh | TCP | 22 | カスタム | [ログイン元PCのIPアドレス] | 作業者のPCからEC2インスタンスへのSSHログイン |
| カスタム TCP | TCP | 50051-50060 | カスタム | 0.0.0.0/0 | クライアント端末からMMPFサーバーへのリクエスト |
| カスタム TCP | TCP | 6379 | カスタム | 0.0.0.0/0 | viewerからredisへのリクエスト |

### Elastic IPの紐づけ

EC2のダッシュボードを開く。
画面左のサイドメニューで「Elastic IP」をクリックする。
画面右上の「Elastic IP アドレスを割り当てる」をクリックする。
入力欄の値は全て初期値のまま、画面右下の「割り当て」をクリックする。
画面右上の「アクション」から表示されるメニューの「Elastic IP アドレスの関連付け」をクリックする。
以下の入力欄にMMPF用に起動したEC2インスタンスの値を入力する。
  - インスタンス
   - プライベート IP アドレス

「関連付ける」をクリックする。
  - MMPF用に起動したEC2インスタンスの概要画面を開き、「パブリック IPv4 アドレス」の欄にIPアドレスが表示されていることを確認する

### 各種ツールのインストール

TeraTermなどを用いてEC2インスタンスにSSHログインする。

インスタンスにログインした状態で以下の通りコマンドを実行し、git、docker、docker-composeをインストールする。
  ```
  sudo apt-get update
  sudo apt-get install git
  sudo apt-get install docker
  sudo apt-get install docker.io
  sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
  sudo chmod +x /usr/local/bin/docker-compose
  ```

### ソースコードのクローン

インスタンスにログインした状態で以下の通りコマンドを実行し、GitHubリポジトリからmmpf-monolithicのソースコードをクローンする。
  ```
  cd ~/
  git clone https://{GITHUB_USERNAME}@github.com/machinemapplatform/mmpf-monolithic.git
  ```

### MMPFサーバーの起動

インスタンスにログインした状態で以下の通りコマンドを実行し、サーバー起動時に必要な権限を設定する。
  ```
  cd ~/mmpf-monolithic/
  chmod +x mmpf.sh
  chmod +x -R scripts/
  sudo gpasswd -a ubuntu docker
  sudo chmod 666 /var/run/docker.sock
  ```

続けて以下のコマンドを実行し、GHCRにログインしておく。
  ```
  echo $(GITHUB_ACCESS_TOKEN) | docker login ghcr.io -u $(GITHUB_USERNAME) --password-stdin
  ```

[インターネット接続環境の手順](../README.md#インターネット接続環境の手順)に従って、モノリシック版MMPFサーバーを起動する。
