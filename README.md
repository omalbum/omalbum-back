# teleOMA API

## Setup

1. Install go (1.14.1+) following the [instructions](https://golang.org/dl/).
2. Setup your [IDE](https://www.jetbrains.com/go/download/)
3. Clone the repo

## Build and Run
If you want to build and run the project from the command line, follow this steps.

1. CD to the appropriate directory

```sh
cd cmd/api
```

2. Build the project
```sh
go build -v .
```

3. Run the project
```sh
./api
```

## Run with Docker

### Project
```sh
docker build -t teleoma .
docker run -it --rm --network host --name teleoma teleoma
```

### MySQL

```sh
docker run -it --rm --name mysql-teleoma -p 3306:3306 -v db-data:/var/lib/mysql-teleoma -e MYSQL_ROOT_PASSWORD=teleoma -e MYSQL_PASSWORD=teleoma -e MYSQL_USER=teleoma -e MYSQL_DATABASE=teleoma mysql:8.0.19
```

## Build and Run on Testing Server 
PENDING


## Build and Run on Production Server
1. Tag a commit as a release version with the following rule:
``` vM.m.p ``` being M: Major version, m: minor version and p: patch
2. Push tag:
``` git push --tag ```

Example:
1. ``` git tag v1.0.2 ```
2. ``` git push --tag ```
 
