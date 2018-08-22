# Getting started

1. Install [dep](https://github.com/golang/dep)
2. Run `make install && make run`

# API

## Create user

```
curl --header "Content-Type: application/json" \
 --request POST \
 --data '{
        "name": "One",
        "surname": "Eno",
        "balance": 100,
        "id": 1
}' \
http://localhost:8000/account
```

## Create another user 

```
curl --header "Content-Type: application/json" \
 --request POST \
 --data '{
        "name": "Two",
        "surname": "Three",
        "balance": 100,
        "id": 2
}' \
http://localhost:8000/account
```

## Get user info

```
curl http://localhost:8000/account/1
```

## Send money 

```
curl --header "Content-Type: application/json" \
 --request POST \
 --data '{
        "from": 1,
        "to": 2,
        "amount": 42
}' \
http://localhost:8000/sendMoney
```