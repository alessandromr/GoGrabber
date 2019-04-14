# GoGrabber
GoGrabber is a "WebCrawler". The main objective of GoGrabber is grabbing URL from websites recursively.

# Main Requirements
- Golang
- Docker
- Docker Compose

##### Docker Containers require
- Redis
- MySQL
- PHPMyAdmin

# Installation
## Docker
```
#Start Docker
docker-compose up -d

#Stop Docker
docker-compose down
```
## Go
```
#Start
cd src
go run .

#Stop
CTRL + C
```

# How it works
GoGrabber is based on Worker and Jobs pattern and tries to use vastly Go Concurrency.

## Queue

- URLs Queue
- Recent URLs Queue

In this version Queues are stored in memory. This choise can be easily swapped with some dedicated Message Queue.
