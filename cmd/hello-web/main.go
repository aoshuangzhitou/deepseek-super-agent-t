// Command hello-web serves a single page that outputs hello.
package main

import (
	"flag"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// helloBody is the plain-text content served to non-browser clients.
const helloBody = "hello"

// helloWebPage is the document served to browsers that ask for HTML. It is one
// self-contained page with no external assets, so it still renders offline and
// the API contract (plain text "hello") stays intact for curl. The
// pageAddressPlaceholder is replaced with the address the server listens on.
const helloWebPage = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="color-scheme" content="dark">
<meta name="description" content="hello — a tiny Go web service that greets you.">
<title>hello · go web service</title>
<link rel="icon" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 64 64'%3E%3Crect width='64' height='64' rx='16' fill='%230b1020'/%3E%3Ctext x='32' y='43' font-family='system-ui,sans-serif' font-size='32' font-weight='700' fill='%238b9cf6' text-anchor='middle'%3Eh%3C/text%3E%3C/svg%3E">
<style>
  *, *::before, *::after { box-sizing: border-box; }

  :root {
    --bg: #070b14;
    --ink: #eef2ff;
    --muted: #98a2c3;
    --line: rgba(148, 163, 255, 0.16);
    --cyan: #67e8f9;
    --mono: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
  }

  html, body { height: 100%; }

  body {
    margin: 0;
    min-height: 100vh;
    display: grid;
    place-items: center;
    padding: clamp(20px, 6vw, 56px) clamp(16px, 5vw, 40px);
    background:
      radial-gradient(1200px 600px at 12% -12%, rgba(99, 102, 241, 0.34), transparent 60%),
      radial-gradient(900px 520px at 110% 8%, rgba(6, 182, 212, 0.18), transparent 55%),
      radial-gradient(760px 520px at 52% 118%, rgba(168, 85, 247, 0.18), transparent 60%),
      var(--bg);
    color: var(--ink);
    font-family: ui-sans-serif, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    overflow-x: hidden;
  }

  .glow {
    position: fixed;
    z-index: 0;
    border-radius: 50%;
    filter: blur(72px);
    pointer-events: none;
  }
  .glow-a {
    width: 380px; height: 380px;
    top: -120px; left: -90px;
    background: radial-gradient(circle, rgba(129, 140, 248, 0.42), transparent 70%);
    animation: drift-a 22s ease-in-out infinite alternate;
  }
  .glow-b {
    width: 320px; height: 320px;
    right: -100px; bottom: -120px;
    background: radial-gradient(circle, rgba(34, 211, 238, 0.28), transparent 70%);
    animation: drift-b 26s ease-in-out infinite alternate;
  }

  .card {
    position: relative;
    z-index: 1;
    width: min(100%, 620px);
    padding: clamp(30px, 6vw, 56px);
    border-radius: 28px;
    border: 1px solid var(--line);
    background: linear-gradient(180deg, rgba(30, 41, 84, 0.55), rgba(9, 14, 30, 0.68));
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.06),
      0 40px 80px -40px rgba(2, 6, 23, 0.95),
      0 0 0 1px rgba(2, 6, 23, 0.35);
    -webkit-backdrop-filter: blur(18px) saturate(120%);
    backdrop-filter: blur(18px) saturate(120%);
    animation: rise 0.7s cubic-bezier(0.2, 0.7, 0.2, 1) both;
    overflow: hidden;
  }
  .card::after {
    content: "";
    position: absolute;
    inset: 0 0 auto;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(165, 180, 252, 0.6), transparent);
  }

  .eyebrow {
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 0;
    font-family: var(--mono);
    font-size: 11.5px;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: var(--muted);
  }
  .dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: #34d399;
    box-shadow: 0 0 0 4px rgba(52, 211, 153, 0.14);
    animation: pulse 2.6s ease-in-out infinite;
  }

  h1 {
    margin: 20px 0 14px;
    font-size: clamp(64px, 15vw, 112px);
    line-height: 0.95;
    font-weight: 750;
    letter-spacing: -0.05em;
    color: var(--ink);
  }
  /* Gradient text only where background-clip works; otherwise the solid
     colour above keeps the wordmark readable. */
  @supports ((background-clip: text) or (-webkit-background-clip: text)) {
    h1 {
      background: linear-gradient(108deg, #ffffff 8%, #c7d2fe 46%, #67e8f9 92%);
      -webkit-background-clip: text;
      background-clip: text;
      -webkit-text-fill-color: transparent;
      color: transparent;
    }
  }

  .lede {
    margin: 0 0 28px;
    max-width: 44ch;
    font-size: clamp(15px, 2.2vw, 17px);
    line-height: 1.65;
    color: var(--muted);
  }

  .chips {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    margin: 0 0 28px;
    padding: 0;
    list-style: none;
  }
  .chips li {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 8px 13px;
    border-radius: 999px;
    border: 1px solid var(--line);
    background: rgba(148, 163, 255, 0.07);
    font-size: 12.5px;
    color: #dbe3ff;
  }
  .chips code { font-family: var(--mono); font-size: 0.95em; color: var(--cyan); }
  .chips .ok { color: #34d399; font-size: 9px; line-height: 1; }

  .terminal {
    margin: 0;
    padding: 16px 18px;
    border-radius: 16px;
    border: 1px solid rgba(148, 163, 255, 0.12);
    background: rgba(2, 6, 23, 0.72);
    font-family: var(--mono);
    font-size: 13.5px;
    line-height: 1.75;
    color: #cbd5f5;
    white-space: pre-wrap;
    overflow-wrap: anywhere;
  }
  .terminal .prompt { color: var(--cyan); }
  .terminal .out { color: #86efac; }

  footer {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    gap: 10px 18px;
    margin-top: 26px;
    padding-top: 18px;
    border-top: 1px solid var(--line);
    font-size: 12px;
    color: #7c87a8;
  }
  footer code { font-family: var(--mono); color: #a9b4d6; }

  @keyframes rise {
    from { opacity: 0; transform: translateY(14px) scale(0.985); }
    to { opacity: 1; transform: none; }
  }
  @keyframes pulse {
    50% { box-shadow: 0 0 0 8px rgba(52, 211, 153, 0); }
  }
  @keyframes drift-a {
    to { transform: translate(70px, 50px) scale(1.12); }
  }
  @keyframes drift-b {
    to { transform: translate(-60px, -40px) scale(1.08); }
  }

  @media (max-width: 480px) {
    .card { border-radius: 22px; }
    footer { justify-content: flex-start; }
  }
  @media (prefers-reduced-motion: reduce) {
    *, *::before, *::after {
      animation: none !important;
      transition: none !important;
    }
  }
  /* Keep text visible in high-contrast mode, where background images and
     gradient fills are dropped. */
  @media (forced-colors: active) {
    h1 {
      background: none;
      -webkit-text-fill-color: currentColor;
      color: CanvasText;
    }
    .card, .terminal { border-color: CanvasText; }
  }
</style>
</head>
<body>
  <div class="glow glow-a" aria-hidden="true"></div>
  <div class="glow glow-b" aria-hidden="true"></div>

  <main class="card">
    <p class="eyebrow"><span class="dot" aria-hidden="true"></span>go · net/http · {{ADDR}}</p>
    <h1>hello</h1>
    <p class="lede">A tiny web service that greets you — plain text on the wire, a little prettier in the browser.</p>

    <ul class="chips" aria-label="Service contract">
      <li><code>GET /</code></li>
      <li><span class="ok" aria-hidden="true">●</span>200 OK</li>
      <li><code>text/plain</code> for API clients</li>
    </ul>

    <pre class="terminal"><span class="prompt">$</span> curl {{ADDR}}/
<span class="out">hello</span></pre>

    <footer>
      <span>no frameworks · no build step</span>
      <span>other paths → <code>404</code></span>
    </footer>
  </main>
</body>
</html>
`

// pageAddressPlaceholder is filled with the host the browser requested.
const pageAddressPlaceholder = "{{ADDR}}"

// defaultHost is shown when a request carries no Host header, which only
// happens for hand-written HTTP requests.
const defaultHost = "127.0.0.1:8123"

// pageHTML returns the styled page with host interpolated into the examples.
// The host comes from the request, so running the service on another port or
// behind a proxy still shows the address the reader actually used.
func pageHTML(host string) string {
	if host == "" {
		host = defaultHost
	}
	return strings.ReplaceAll(helloWebPage, pageAddressPlaceholder, html.EscapeString(host))
}

// newHandler returns the HTTP handler for the hello page.
func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHello)
	return mux
}

// handleHello writes the hello page for GET and HEAD requests to the root
// path. Browsers that accept HTML get the styled page; every other client
// keeps receiving the plain-text body. Other paths return 404 and other
// methods return 405.
func handleHello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	body := helloBody
	contentType := "text/plain; charset=utf-8"
	if wantsHTML(r) {
		body = pageHTML(r.Host)
		contentType = "text/html; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodGet {
		_, _ = w.Write([]byte(body))
	}
}

// wantsHTML reports whether the request's Accept header lists text/html,
// which is how browsers identify themselves. CLI clients such as curl send
// */* and keep getting the plain-text response.
func wantsHTML(r *http.Request) bool {
	for _, part := range strings.Split(r.Header.Get("Accept"), ",") {
		media := part
		if i := strings.IndexByte(media, ';'); i >= 0 {
			media = media[:i]
		}
		if strings.EqualFold(strings.TrimSpace(media), "text/html") {
			return true
		}
	}
	return false
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8123", "address to listen on")
	flag.Parse()

	log.Printf("hello-web listening on http://%s/", *addr)
	srv := &http.Server{
		Addr:              *addr,
		Handler:           newHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
