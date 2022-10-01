# Oauth2 Golang Mysql

- Mysql to store token

## install

```bash
git clone https://github.com/danangkonang/oauth2-golang.git

cd oauth2-golang

go mod tidy

#database(optional)
docker-compose up -d

make up

go run main.go
```

## response

```json
{
  "access_token": "sDrk..",
  "token_type": "Bearer",
  "refresh_token": "jkudy8..",
  "expiry": 7195,
}
```
