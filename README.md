## Infra
```shell
make infra-up
make infra-down
```

## Migrations
```shell
go run cmd/main.go migration --up 0
go run cmd/main.go migration --down 0
```

## SFarmer
```shell
go run cmd/main.go sfarmer
go run cmd/main.go sfarmer --test=false
```

## SFarmer In Container
```shell
make sfarmer-up
```

## SymList Updater
```shell
go run cmd/main.go symlist
```

## Wavetrend Calculator
```shell
go run cmd/main.go wtmomentum --symlist=./files/bought.txt
```