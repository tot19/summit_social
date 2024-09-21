# summit_social

## MUST INSTALL
- podman
- podman-compose

## Building locally
```
podman-compose build --no-cache
podman-compose up -d
podman-compose ps
```

## Deleting local build
```
podman-compose down
podman pod rm -f pod_summit_social
podman-compose ps
```


## Checking logs
```
podman-compose logs -f
podman-compose logs -f frontend
podman-compose logs -f backend
```