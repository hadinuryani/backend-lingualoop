# Backend LinguaLoop

---

## Installation

Clone repository

```bash
git clone https://github.com/hadinuryani/backend-lingualoop.git
cd backend-lingualoop
```

Install dependency

```bash
go mod tidy
```

Copy environment

```bash
cp .env.example .env
```

setup  `.env`.

---

## Database Setup

### 1. Jalankan Migration

```bash
go run cmd/migration/main.go
```

### 2. Jalankan Seeder

```bash
go run cmd/seed/main.go
```

---

##  Run Application

```bash
go run cmd/api/main.go
```

---

## API Documentation

iki nk kene:

```
http://localhost:8080/swagger/index.html
```

---

```text
                        ____             _                  _
                        | __ )  __ _  ___| | _____ _ __   __| |
                        |  _ \ / _` |/ __| |/ / _ \ '_ \ / _` |
                        | |_) | (_| | (__|   <  __/ | | | (_| |
                        |____/ \__,_|\___|_|\_\___|_| |_|\__,_|

                _     _                         _
                | |   (_)_ __   __ _ _   _  __ _| |    ___   ___  _ __
                | |   | | '_ \ / _` | | | |/ _` | |   / _ \ / _ \| '_ \
                | |___| | | | | (_| | |_| | (_| | |__| (_) | (_) | |_) |
                |_____|_|_| |_|\__, |\__,_|\__,_|_____\___/ \___/| .__/
                                |___/                            |_|
```
