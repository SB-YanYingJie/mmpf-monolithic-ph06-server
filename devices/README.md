# devices
このドキュメントの対象読者は、MMPF開発者（以下、開発者）と、MMPFサービス運用者（以下、運用者）である。

このフォルダはクライアントデバイス固有のenvファイルと、MMPFサーバ共通のenvファイルを格納するフォルダである。envファイルでMMPFの起動に必要なパラメータ（以下、起動パラメータ）を定義する。

### デプロイ時のenvファイルの設定方法
1. devicesディレクトリ配下に、`.env.xxxx`のようなファイル名で、MMPFを利用する端末ごとにenvファイルを用意する。サーバ共通の設定値に変更があれば`common.env`も編集する。
   1. envファイルの具体的な書き方は[envファイルの書き方](#envファイルの書き方)を参照すること
2. devicesディレクトリ配下に、`env_list.txt`を作成する。
   1. デプロイしたいサービスのenvファイル名(`".env.xxxx"`)を記述する。`"#"`でコメントアウトすることでデプロイ対象から外すことが可能。

#### envファイルの書き方
`リポジトリルート階層/devices/.env.sampleXXXX`を参考に新規`.env.xxxx`ファイルを作成する。各起動パラメータの説明は以下一覧を参照のこと。

| クライアント固有パラメータ                             | 説明                                                                                                   | デフォルト値    |
| --------------------------------- | ---------------------------------------------------------------------------------------------------- | --------- |
| MAPID                             | SLAMに使うマップを一意に特定するID                                                                                 | testMAPID         | 
| MMID                              | MMPFを一意に特定するID                                                                                       | -         |
| HOST_PORT                         | クライアントアプリのリクエストを待ち受けるポート番号                                                                           | -         |
| PORT                         | ホストのポートを転送し、コンテナで待ち受けるポート番号                                                                           | -         | 
| REDIS_PUBSUB_CHANNEL_POSE         | POSEチャネルのチャネル名                                                                                       | -         |
| SAVE_IMAGE_EXT                    | 画像ファイルを保存するときの拡張子                                                                                    | -         |
| NUMBER_OF_LENSES                  | クライアントのカメラ個数(mono/stereo)                                                                            | -         |
| TRIMMING_PARAMETER                | クライアントから送信される画像が結合されている場合、それを分割するためのパラメータ<br>(左上を基準として、切り取り開始点x, 切り取り開始点y, 切り取り幅width, 切り取り高さheight) | `0:0:0:0` |
| SEND_IMAGE_TYPE                   | クライアントから送信される画像形式 (mono, stereo_separated, stereo_merged)                                            | -         |
| KD_CALIB_PATH                     | calibrationファイルのパス                                                                                   | `/app/lib`         |
| KD_MAP_PATH                       | mapファイルのパス                                                                                           | `/app/lib`         |
| KD_MAP_EXPAND                     | map を自動拡張するフラグ                                    | false  |
| TARGET_FPS                        | フレームレート                                                                                              | -         |
| （以下、開発用パラメータ）                     |                                                                                                      | -         |
| DEV_DISPLAY_RAW_IMAGE             | 読み込んだ画像そのものを、サーバーの画面に表示するフラグ。trueの場合は、実行環境に応じてDISPLAYを変える必要がある。デフォルトはfalse                           | false     |
| DEV_DISPLAY_DEBUG_IMAGE           | 読み込んだ画像（kdslamのデバッグ画像）を、サーバーの画面に表示するフラグ。trueの場合は、実行環境に応じてDISPLAYを変える必要がある。デフォルトはfalse                | false       | false     |
##### 【開発者向け】起動パラメータ補足
起動パラメータにはdocker containerの立ち上げに必要なパラメータとサーバーアプリ内で使用するパラメータが存在する。上記一覧のうち`HOST_PORT`はdocker containerの立ち上げに必要なパラメータであり、それ以外はサーバーアプリ内で使用するパラメータである。

---
`リポジトリルート階層/devices/common.env`に記載の起動パラーメータの説明は、以下一覧を参照のこと。
| サーバ共通パラメータ | 説明 | デフォルト値 |
|-|-|-|
| GOPRIVATE                         | インポート可能なプライベートライブラリのURL                           | `github.com/machinemapplatform/*,github.com/KudanJP/*`      |
| LOG_SETTINGS_FILE_PATH            | ログレベルの設定ファイルのパス                                   | `/app/zap/prd/config.yml`      |
| SERVICE_NAME_MONOLITHIC                 | serviceのホスト名                                      | machinemapplatform-monolithic-service      |
| KD_VOCAB_PATH                     | vocabularyファイルのパス                                 | `/usr/local/lib/ORBvoc.kdbow`      |
| LD_LIBRARY_PATH                   | KdSlamGoとgocvが依存するファイルの参照先                        | `/usr/local/lib:/app/lib`      |
| CGO_LDFLAGS                       | KdSlamが依存するファイルの参照先                               | `"-L/usr/local/lib -lKdSlam"`      |
| PKG_CONFIG_PATH                   | gocvが依存するファイルの参照先                                 | `/usr/lib/x86_64-linux-gnu/pkgconfig:/usr/local/lib/pkgconfig`      |
| IMAGE_STORE_REDIS_TTL             | redisの画像保存期間(秒)。 0 : TTLを設定しない | 3      |
| REDIS_ADDRESS                     | redisの接続先                                         | `redis:6379`      |
| REDIS_IDLE_TIMEOUT_SECONDS        | コネクションがアイドル状態で残存可能な時間（秒）                          | 3      |
| REDIS_MAX_IDLE                    | アイドル状態のコネクションがプールに接続できる最大数                        | 5      |
| REDIS_PUBSUB_DB                   | DBを特定するための数値                                      | 2      |
| （以下、開発用パラメータ）                     | -                                                 | -      |
| DISPLAY                           | （開発用）KdSlamの実行結果を描画する際に使用                         | -      |
