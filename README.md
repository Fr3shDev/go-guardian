# go-guardian
A command‑line tool that monitors the health of websites. It performs real‑time checks for:

- Uptime and response time: Makes an HTTP GET request to verify the website is online and measures how fast it responds.
- SSL certificate expiration: For HTTPS websites, it inspects the SSL certificate and warns you if it’s nearing expiration (for example, within 30 days).

This tool can help developers, system administrators, or even nontechnical users who run websites to stay on top of website health and security.

# Next Steps & Enhancements
- Configuration File: Instead of hardcoding URLs, read them from a JSON or YAML config file.
- Notifications: Integrate email or Slack notifications when a website is down or an SSL certificate is expiring.
- Logging: Add file logging for historical data and alert trends.
- CLI Flags: Use the flag package to allow users to set parameters (like check interval or threshold) via command-line arguments 