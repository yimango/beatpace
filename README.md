# BeatPace

BeatPace is a rhythm-based Spotify web application that generates playlists customized to your running pace. It analyzes your target speed, height, and gender to estimate your ideal cadence, then curates music to match your stride.

Built with a modern stack—Next.js, Tailwind CSS, and Go with Gin—BeatPace offers a clean developer-first architecture that connects seamlessly with Spotify’s API. It enables an immersive running experience by syncing music tempo with footfalls.

> **Note:** Due to recent changes to the Spotify Developer API policy, public deployment of BeatPace is currently not possible. Only development-mode usage is supported with Spotify account whitelisting.

---

## Features

- **Spotify Login via OAuth 2.0**  
  Secure login and token handling with refresh token management.

- **Cadence-Aware Playlist Generation**  
  Converts user pace to BPM range, then builds a tempo-matched playlist using Spotify’s recommendation API.

- **Persistent User Sessions (JWT)**  
  Sessions are issued and validated using signed JWTs, with server-side revocation via database.

- **Modern UI**  
  Built with Next.js App Router, Radix UI, and Tailwind CSS for accessible, responsive interactions.

- **MySQL Integration**  
  Tokens, sessions, and user profiles are stored and managed through a robust MySQL backend.

---

## Installation (Development Mode Only)

Due to Spotify policy changes, BeatPace cannot be deployed publicly without app verification. However, developers can run a local version:

### 1. Clone the Repository
```bash
git clone https://github.com/yimango/beatpace.git
cd beatpace
```

### 2. Set Up Your Spotify Developer App
- Create an app at [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
- Add the redirect URI: `http://localhost:3001/api/callback`

### 3. Configure Environment Variables
In `beatpace-backend`, create a `.env` file:
```env
CLIENT_ID=your_spotify_client_id
CLIENT_SECRET=your_spotify_secret
REDIRECT_URI=http://localhost:3001/api/callback
JWT_SECRET=your_jwt_secret
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_mysql_user
DB_PASS=your_mysql_password
DB_NAME=beatpace
FRONTEND_URL=http://localhost:3000
```

### 4. Prepare the Database
Ensure MySQL is running, then apply [migrations.sql](beatpace-backend/db/migrations.sql).

### 5. Run Backend
```bash
cd beatpace-backend
go run .
```

### 6. Run Frontend
```bash
cd ../beatpace-frontend
npm install
npm run dev
```

---

## Usage

1. **Log in via Spotify**
   - Redirects to Spotify’s OAuth screen and returns with a JWT.

2. **Enter Your Details**
   - Target pace (e.g., `5:00/km`), height, and gender.

3. **Generate Playlist**
   - Backend computes target BPM and fetches matching tracks via Spotify API.

4. **View Your Playlist**
   - Playlist is saved to your Spotify account and linked in the UI.
  
---

## Architecture Overview

- **Frontend:**
  - Next.js App Router
  - Tailwind CSS, Radix UI
  - Auth state managed client-side using JWT

- **Backend:**
  - Go + Gin for API routing
  - JWT auth middleware
  - Spotify API integration via `zmb3/spotify`
  - MySQL for persistent user/session/token storage

- **Database:**
  - Tables: `users`, `spotify_tokens`, `sessions`
  - Migrations provided in SQL format

- **Deployment:**
  - Development mode only (due to Spotify API limits)
  - Requires account whitelisting via Spotify Dashboard

---

## License

This project is licensed under the [MIT License](LICENSE).

---

> Built by [@yimango](https://github.com/yimango). For devs, runners, and rhythm chasers alike.
