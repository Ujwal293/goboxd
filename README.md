# Goboxd

A lightweight multi-language code execution backend built with Go, Docker, Python, and C++ support.

## Features

- Execute Python code
- Execute C++ code
- REST API
- Dockerized runtime
- Health checks
- Readiness checks
- Temporary isolated execution
- JSON-based API

---

## Endpoints

### Health Check

GET /healthz

Response:

{
  "status": "ok"
}

---

### Readiness Check

GET /readyz

Response:

{
  "python": "ok",
  "g++": "ok",
  "status": "ready"
}

---

### Run Code

POST /run

Example Request:

{
  "language": "py3",
  "source": "print('hello world')"
}

---

## Run Locally
$env:Path += ";C:\msys64\ucrt64\bin"
go run main.go

---

## Run with Docker

docker compose up --build

---

## Supported Languages

- Python
- C++

---

## Tech Stack

- Go
- Docker
- Python
- g++
- REST API