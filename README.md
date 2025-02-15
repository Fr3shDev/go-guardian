# go-guardian
A command‑line tool that monitors the health of websites. It performs real‑time checks for:

- Uptime and response time: Makes an HTTP GET request to verify the website is online and measures how fast it responds.
- SSL certificate expiration: For HTTPS websites, it inspects the SSL certificate and warns you if it’s nearing expiration (for example, within 30 days).

This tool can help developers, system administrators, or even nontechnical users who run websites to stay on top of website health and security.