# API for concerts service

## How to setup
To use this concerts-api server you also need to setup keycloak server for authorization.
There are a lot of setup, but you should do it only once if you will not delete keycloak database.

1. Clone the repo:
```sh
git clone https://github.com/mykykh/concerts-api.git
```

2. Create .env file:
```sh
cp example.env .env
```

3. Add keycloak to /etc/hosts so concerts-api would redirect you to valid domain:
```sh
127.0.0.1 keycloak
```

4. Launch docker containsers:
```sh
docker compose up
```

Now you need to setup keycloak:

5. Login to keycloak admin pannel at keycloak:8081
6. Create new realm
7. Create new client
8. Add realm, client id and client secret to .env file (see example.env)

Now keycloak server is ready, but you also need to give permisions for concerts-api to your tokens:
9. Add client scope concerts-api-scope with access to created client id.
10. Add client scope concerts-api-scope as default to access token in client
11. Relaunch docker compose

12. Now you can use concerts-api, but if you register users they will not have almost any permisions, so add this roles
in concerts-api client in keycloak admin panel:
```
write:concerts
read:concerts
update:concerts
update:ownConcerts

write:tickets
read:tickets
read:ownTickets
update:tickets
update:ownTickets
```

13. Now you will want to create any user group to add this roles. For example create Admins group with all roles.
14. You also will want to enable user registration. So go into realm settings, login page and turn on user registration.
15. And lastly you will need to add roles to default registred users. Go into realm settings, user registration page,
default group, and add any group that you want. For testing you can add Admin group as default

16. Now everything is ready. Go into localhost:8080/auth/login. It should redirect you into
keycloak server where you can login or register. After that you should be redirected back
to localhost:8080/auth/callback, where you will get your access token.

17. If you want to regenerate your token, go again into /auth/login, because refreshing /auth/callback will not work.
18. You can use your token with swagger ui in localhost:8080/swagger/index.html
19. Now, if you don't delete keycloak-db docker container, you don't need to redo all of the setup above to update concerts-api.
Just do: docker compose up --build, and everything should work.
