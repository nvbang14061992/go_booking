Here is the complete flow formatted as clean Markdown content that you can copy into your documentation:

---

# HTTP Redirect Flow: Login → Homepage (Practical Example)

This document explains how HTTP redirects work in a real-world login scenario using a Go backend.

Assumptions:

- Domain: `https://app.example.com`
- Backend written in Go
- Possibly behind a reverse proxy like NGINX
- Client is a browser

---

# Scenario

User tries to access a protected page:

```
/dashboard
```

If not authenticated:

1. Redirect to `/login`
2. User logs in
3. Redirect to `/` (homepage)

---

# Step-by-Step HTTP Flow

---

## Step 1 — User Requests Protected Page

Browser sends:

```
GET /dashboard HTTP/1.1
Host: app.example.com
Cookie: session=abc123
```

---

## Step 2 — Backend Detects User Is Not Logged In

Go handler:

```go
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    if !isAuthenticated(r) {
        http.Redirect(w, r, "/login", http.StatusFound) // 302
        return
    }

    // Serve dashboard
}
```

Server responds:

```
HTTP/1.1 302 Found
Location: /login
```

Important:

- No HTML body required
- No browser logic involved
- Only status code + Location header

---

## Step 3 — Browser Automatically Follows Redirect

Browser receives:

```
302 Location: /login
```

Browser then makes a new request automatically:

```
GET /login HTTP/1.1
Host: app.example.com
```

Because `/login` is relative, the browser resolves it as:

```
https://app.example.com/login
```

---

## Step 4 — Backend Serves Login Page

Server responds:

```
HTTP/1.1 200 OK
Content-Type: text/html
```

Login page is displayed.

---

## Step 5 — User Submits Login Form

Browser sends:

```
POST /login HTTP/1.1
Host: app.example.com
Content-Type: application/x-www-form-urlencoded
```

Body:

```
username=alice&password=123
```

---

## Step 6 — Backend Authenticates and Sets Session

Go handler:

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    // Validate credentials

    http.SetCookie(w, &http.Cookie{
        Name:  "session",
        Value: "new-session-token",
        Path:  "/",
    })

    http.Redirect(w, r, "/", http.StatusSeeOther) // 303
}
```

Server responds:

```
HTTP/1.1 303 See Other
Set-Cookie: session=new-session-token
Location: /
```

Why `303 See Other`?

- Prevents form re-submission
- Forces browser to use GET for the next request

---

## Step 7 — Browser Redirects to Homepage

Browser automatically sends:

```
GET / HTTP/1.1
Host: app.example.com
Cookie: session=new-session-token
```

---

## Step 8 — Homepage Loads Successfully

Server responds:

```
HTTP/1.1 200 OK
```

User is now logged in and sees homepage.

---

# Visual Flow Summary

```
GET /dashboard
    ↓
302 → /login
    ↓
GET /login
    ↓
POST /login
    ↓
303 → /
    ↓
GET /
    ↓
200 OK
```

---

# Key Concepts

## 1. Redirect Is Not Internal Forwarding

A redirect:

- Ends the current request
- Sends a 3xx response
- Client initiates a new request

There are always two separate HTTP requests.

---

## 2. Backend Does Not “Send” the User Anywhere

The backend only responds with:

```
Location: /login
```

The client (browser) decides to make the next request.

---

## 3. Why Relative URLs Work

Because the browser already knows the origin:

```
https://app.example.com
```

So:

```
/login → https://app.example.com/login
```

---

## 4. Reverse Proxy Does Not Change Redirect Logic

Even if architecture is:

```
Browser (HTTPS)
    ↓
Reverse Proxy
    ↓
Go Backend (HTTP)
```

Redirect behavior is unchanged because:

- Browser handles redirects
- Backend only returns status + Location header

---

# Final Mental Model

An HTTP redirect always follows this pattern:

1. Server responds with 3xx status
2. Response includes Location header
3. Client sends a new request to that Location

That is the complete mechanism.

---

If you want, I can also generate a version tailored for:

- SPA + fetch()
- REST API login (JSON instead of HTML)
- OAuth flow
- Kubernetes / cloud environment

Just let me know your actual architecture.
