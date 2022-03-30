package configs

import (
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var (
	AuthorizationHeader = "Authorization"
	Config              *ConfigType
	origConfig          ConfigType
	secretVersion       int64 = 0
)

type ConfigType struct {
	Port                   int    `env:"PORT" envDefault:"8080"`
	MongoURI               string `env:"MONGO_URI" envDefault:""`
	MongoDB                string `env:"MONGO_DB" envDefault:""`
	GBucketName            string `env:"GCLOUD_BUCKET_NAME" envDefault:""`
	GProjectID             string `env:"GCLOUD_PROJECT_ID" envDefault:""`
	LineChannelSecret      string `env:"LINE_CHANNEL_SECRET" envDefault:""`
	LineChannelAccessToken string `env:"LINE_CHANNEL_ACCESS_TOKEN" envDefault:""`
	RebrandlyKey           string `env:"REBRANDLY_API" envDefault:""`
	RebrandlyWorkspace     string `env:"REBRANDLY_WORDSPACE" envDefault:""`
	AirVisualKey           string `env:"AIR_API_KEY" envDefault:""`
	SecretPassphrase       string `env:"SECRET_PASSPHRASE" envDefault:""`
	DISCORD_TOKEN          string `env:"DISCORD_TOKEN" envDefault:""`
	DISCORD_ERROR_CHANNEL  string `env:"DISCORD_ERROR_CHANNEL" envDefault:""`
	DISCORD_SERVER_ID      string `env:"DISCORD_SERVER_ID" envDefault:""`
	LineNotifyURL          string `env:"MICROSERVICE_BINANCE_LINE_NOTIFY_URL" envDefault:""`

	LineAPIBroadcast        string `envDefault:"https://api.line.me/v2/bot/message/broadcast"`
	LineAPIReply            string `envDefault:"https://api.line.me/v2/bot/message/reply"`
	AirVisualAPINearestCity string `envDefault:"http://api.airvisual.com/v2/nearest_city"`
	AirVisualAPICity        string `envDefault:"http://api.airvisual.com/v2/city"`
	BinanceAccountAPI       string `envDefault:"https://api.binance.com/api/v3/account"`
	RebrandlyAPI            string `envDefault:"https://api.rebrandly.com/v1/links"`
	LogFilePath             string `env:"LOG_PATH" envDefault:"/var/log/cjtim-backend-go.log"`
	GCLOUD_CREDENTIAL       string `env:"GCLOUD_CREDENTIAL" envDefault:"./configs/serviceAcc.json"`
}

func init() {
	log.Default().Println("Initial config...")
	fp, err := os.Create("/var/log/cjtim-backend-go.log")
	if err != nil {
		os.Setenv("LOG_PATH", "./log/cjtim-backend.go.log")
	}
	defer fp.Close()

	client, err := newVault()
	if err != nil {
		log.Default().Println("Vault secret error:", err.Error())
		return
	}

	loadEnv(client)

	c := cron.New()
	c.AddFunc("* * * * *", cronVault(client))
	c.Start()

}

func loadEnv(client *api.Client) {
	log.Default().Println("Loading Vault secret...")
	err := readVault(client)
	if err != nil {
		log.Default().Println("Vault secret error:", err.Error())
	}

	cfg := ConfigType{}
	_ = godotenv.Load()
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	Config = &cfg
	origConfig = cfg
}

func RestoreConfigMock() {
	Config = &origConfig
}
