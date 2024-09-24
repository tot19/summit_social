# SummitSocial
For climbers who want to bond over their favorite hobby. Invite climbers to your climb sessions and SEND together!

![image](https://github.com/user-attachments/assets/27417e29-c8a9-483c-b0b6-a99b8a28a91a)


## MUST INSTALL
- podman
- podman-compose

## OR DOCKER
- docker
- docker-compose

## Building locally
```
podman-compose build --no-cache
podman-compose up -d
podman-compose ps
```

## or docker-compose
```
docker-compose build --no-cache
docker-compose up -d
docker-compose ps
```

## Deleting local build
```
podman-compose down
podman pod rm -f pod_summit_social
podman-compose ps
```

## or docker-compose
```
docker-compose down
docker-compose ps
```

## Checking logs
```
podman-compose logs -f
podman-compose logs -f frontend
podman-compose logs -f backend
```

## or docker-compose
```
docker-compose logs -f
docker-compose logs -f frontend
docker-compose logs -f backend
```
