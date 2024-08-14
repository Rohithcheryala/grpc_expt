# CLIENT

## GO

`go run main.go add --sku "12345" --name "New Item" --description "A new item" --price 29.99 --quantity 100`

`go run main.go remove --sku "12345"`

`go run main.go get --sku "12345"`

`go run main.go update-quantity --sku "12345" --change 10`

`go run main.go update-price --sku "12345" --price 25.99`

`go run main.go watch --sku "12345"`

## RUST

`cargo run --release --bin cli -- add --sku "item001" --price 12.99 --quantity 10 --name "Item Name" --description "Item Description"`

`cargo run --release --bin cli -- remove --sku "item001"`

`cargo run --release --bin cli -- get --sku "item001"`

`cargo run --release --bin cli -- update-quantity --sku "item001" --change 5`

`cargo run --release --bin cli -- update-price --sku "item001" --price 15.99`

`cargo run --release --bin cli -- watch --sku "item001"`

# SERVER

## GO

`go run src\server\server.go`

## RUST

`cargo run --release --bin server`
