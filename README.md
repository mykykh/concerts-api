# API for concerts service

## How to setup
1. Clone the repo:
```sh
git clone https://github.com/mykykh/concerts-api.git
```

2. Create .env file:
```sh
cp example.env .env
```

3. Add keycloak to /etc/hosts:
```sh
127.0.0.1 keycloak
```

4. Launch docker containsers:
```sh
docker compose up
```

5. Login to keycloak admin pannel at keycloak:8081
6. Create new realm test
7. Create new client with id concerts-api
8. Add client id and client secret to .env file
9. Add client scope concerts-api-scope with access to created client id.
10. Add client scope concerts-api-scope as default to access token in client
11. Relaunch docker compose
