Shorten URL:

    Request: POST /shorten
    Body: {"url": "https://example.com"}
    Response: {"short_url": "http://localhost:8080/abc123"}

Retrieve URL:

    Request: GET /abc123
    Response: Redirects to https://example.com
