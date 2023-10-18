module github.com/fumiama/NanoBot

go 1.20

require (
	github.com/FloatTech/floatbox v0.0.0-20231017134949-ae5059ebace7
	github.com/FloatTech/ttl v0.0.0-20230307105452-d6f7b2b647d1
	github.com/FloatTech/zbpctrl v1.6.0
	github.com/RomiChan/syncx v0.0.0-20221202055724-5f842c53020e
	github.com/RomiChan/websocket v1.4.3-0.20220227141055-9b2c6168c9c5
	github.com/fumiama/go-base16384 v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.1
	github.com/wdvxdr1123/ZeroBot v1.7.5-0.20231009162356-57f71b9f5258
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/FloatTech/sqlite v1.6.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fumiama/cron v1.3.0 // indirect
	github.com/fumiama/go-registry v0.2.6 // indirect
	github.com/fumiama/go-simple-protobuf v0.1.0 // indirect
	github.com/fumiama/gofastTEA v0.0.10 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	golang.org/x/sys v0.0.0-20220915200043-7b5979e65e41 // indirect
	golang.org/x/text v0.4.0 // indirect
	modernc.org/libc v1.21.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.4.0 // indirect
	modernc.org/sqlite v1.20.0 // indirect
)

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.20.0-with-win386

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b
