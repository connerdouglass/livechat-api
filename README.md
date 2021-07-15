# livechat-api

Live chat API server written in Go.

## Development
To develop this project locally, you'll need to have Go 1.16 or newer installed on your local machine.

In the root directory of the project, create a file named `.env` with the following contents:

```env
DB_URL="mysql://root@tcp(127.0.0.1:3306)/livestream?charset=utf8&parseTime=True&loc=UTC"
AUTH_TOKEN_SIGNING_PEPPER=8yfds34rtyui98uygfdw45yuouy1
RTMP_SERVER_PASSCODE=3454f56ygfdsertyuio076rseryui76
CORS_ALLOW_ORIGINS=http://localhost:4200
```

These are just example values. You'll probably want to change `DB_URL` for your local environment.

You can even use SQLite, if you want. An example SQLite setup would look like:

```env
DB_URL="sqlite://data.db"
```

Note: Data can safely be stored in `data.db`, and it won't be checked into the Git repo.

Once you've got a `.env` file, just run this to start the server:

```sh
go run .
```
