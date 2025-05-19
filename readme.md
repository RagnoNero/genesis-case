# 🌤️ Weather Subscription Service

This project is a **RESTful API** that allows users to:
- Get current weather for a specific city.
- Subscribe to periodic weather updates via email.
- Confirm subscriptions through a unique token.
- Automatically send weather forecasts according to user-defined frequency (hourly or daily).

## 📦 Stack

- **Golang** — main backend logic.
- **PostgreSQL** — relational database for storing subscriptions.
- **Docker / Docker Compose** — containerization and orchestration.
- **SMTP** — email notifications.
- **OpenWeatherMap API** — weather data source.
- **Go Testing** — unit and component test coverage.

---

## ▶️ How to Run

### 1. Clone the repository

```bash
git clone https://github.com/your-username/weather-subscription.git
cd genesis-case
```

### 2. Create `.env` file

```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=weather_subscriptions

WEATHER_API_KEY=your_openweather_api_key

SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_EMAIL=your@email.com
SMTP_USER=your@email.com
SMTP_PASSWORD=your_password
```

### 3. Build and run with Docker Compose

```bash
docker compose up --build
```

The app will be available at: `http://localhost:8080`

---

## 🛠️ Project Structure

```
.
├── api/                # REST API routes and handlers
├── cmd/                # Application entry point
│   └── main.go
├── config/             # App configuration loader and types
├── di/                 # Dependency injection and setup
├── email/              # Email sending logic
├── http/               # Custom HTTP client
├── migrations/         # SQL migration scripts
├── models/             # Domain models
├── scheduler/          # Background scheduler for sending emails
├── sql/                # DB repository logic
│   └── dto/            # Data transfer objects
├── static/             # HTML form for weather and subscription
├── tests/              # Test suites
│   ├── unit/           # Unit tests
│   └── component/      # Component/integration tests
├── token/              # Token generation logic
├── weatherapi/         # Weather API client
├── Dockerfile
├── docker-compose.yml
└── .env                # Environment configuration
```

---

## 🔄 Workflow Explanation

1. **User opens the web form** served from `/static/index.html`.
2. **If city is selected**, the app fetches and displays the weather using OpenWeatherMap API.
3. **If user fills in email/city/frequency**, a POST request is made to subscribe them.
4. **A token is sent to the user via email** for confirmation.
5. **A background scheduler** checks confirmed subscriptions and sends emails based on frequency.

---

## 🧪 Running Tests

Run all tests:

```bash
go test ./...
```

Run only unit tests:

```bash
go test ./tests/unit/...
```

Run only component tests:

```bash
go test ./tests/component/...
```

---

## 📌 Notes

- Subscriptions are persisted in PostgreSQL with timestamps.
- Emails are sent only after confirmation via token.
- The app uses a memory cache to avoid DB polling every minute.
- Weather is fetched live before sending each email.