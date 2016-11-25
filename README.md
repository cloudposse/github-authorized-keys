# github-authorized-keys
Allow to provide ssh access to servers based on github teams

## Build docker image
```
docker build ./ -t github-authorized-keys
```

## Authorize command

```
docker run \                                                                      [system]
-e GITHUB_API_TOKEN={token} \
-e GITHUB_ORGANIZATION={organization} \
-e GITHUB_TEAM={team} \
github-authorized-keys authorize {user}
```

## Create users

```
docker run \                                                                      [system]
-e GITHUB_API_TOKEN={token} \
-e GITHUB_ORGANIZATION={organization} \
-e GITHUB_TEAM={team} \
-e SYNC_USERS_GID={gid OR empty} \
-e SYNC_USERS_GROUPS={comma separated groups OR empty} \
-e SYNC_USERS_SHELL={user shell} \
-v /etc/passwd:/etc/passwd \
-v /etc/shadow:/etc/shadow \
-v /etc/groups:/etc/groups \
-v /home:/home \
github-authorized-keys sync_users
```
