package goyadisk

const (
	baseURL = "https://cloud-api.yandex.net/v1/disk"
)

type Yadisk struct {
	token  string
	appDir string
}

func New(token string, appDir string) *Yadisk {
	return &Yadisk{token: token, appDir: appDir}
}
