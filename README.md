# GoGrabber

GoGrabber is a "WebCrawler". The main objective of GoGrabber is grabbing URL from websites recursively.

## Important

Using **high rates** and floods server with requests can let you incur in **legal issues**.  
Be sure to check the website's Terms of Service.  
Follow `robots.txt` indications.  

## Version

GoGrabber is still unstable and the first stable version is not out yet.  

### Stable

- `-`

### Dev

- `0.0.1-alpha`

## Main Requirements

- Golang
- Docker
- Docker Compose

### Docker Containers require

- Redis
- MySQL
- PHPMyAdmin

## Installation

### Docker

```bash
#Start Docker
docker-compose up -d

#Stop Docker
docker-compose down
```

### Env Vars

```bash
cp ./src/.env.example ./src/.env
```

### Go

```bash
#Start
cd src
go run .

#Stop
CTRL + C
```

## How it works

GoGrabber is based on Worker and Jobs pattern and tries to use vastly Go Concurrency.

### Queue

- URLs Queue
- Recent URLs Queue

In this version Queues are stored in memory. This choise can be easily swapped with some dedicated Message Queue.
