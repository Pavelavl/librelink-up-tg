# LibreLink Up Telegram Notifier

This application is designed to receive glucose data from LibreLink and send notifications to Telegram if the data exceeds the normal values. The app runs in the background and periodically checks the data.

## Installation

1. **Clone:**

```bash
git clone https://github.com/pavelavl/librelink-up-tg.git
cd librelink-up-tg
```

2. **Configure:**
Create a config/config.yml configuration file with the necessary parameters. [Example](./config/config.sample.yaml)

3. **Run:**
Make sure that you have Go installed (version 1.16 or higher). Then run the command:
```bash
go run ./cmd
```

