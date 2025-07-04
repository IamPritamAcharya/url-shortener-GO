# 🔗 Go URL Shortener

A simple URL shortener built with Go, Gin, and PostgreSQL.

## 🚀 Features

* Shorten long URLs (`POST /shorten`)
* Redirect using short codes (`GET /:code`)

## 🛠 Tech Stack

* Go
* Gin
* PostgreSQL

## 📦 Setup

1. Create a PostgreSQL table:

   CREATE TABLE urls (
   id SERIAL PRIMARY KEY,
   original\_url TEXT NOT NULL,
   short\_code VARCHAR(20) UNIQUE NOT NULL
   );

2. Update your PostgreSQL credentials in `db/connect.go`.

3. Run the app with:

   go run main.go

## 📬 API Usage

* To shorten a URL:

  * Send a POST request to `/shorten` with JSON body:
    { "url": "[https://example.com](https://example.com)" }

* To redirect:

  * Open or GET `/short_code`, e.g., `/a1B2c3` → Redirects to original URL

