# 🚭 Snusfria

**Snusfria** is a command-line companion for anyone quitting snus (Swedish smokeless tobacco). Track your progress, log relapses and moods, visualize your health recovery, and see your financial savings—all in your terminal.

---

## Features

- **Status Dashboard:** See your clean time, relapses, savings, and health milestones.
- **Mood & Craving Logs:** Track your mental state and cravings over time.
- **Relapse Tracking:** Log relapses with notes and pouches consumed.
- **Purchase Logging:** Track dosa purchases and see cost trends.
- **Savings Goals:** Set and visualize financial goals.
- **Health Timeline:** See your recovery milestones and what’s happening in your body.
- **Streaks:** Visualize your clean streaks and personal bests.

---

## Installation

1. **Clone the repo:**
   ```sh
   git clone https://github.com/yourusername/snusfria.git
   cd snusfria
   ```

2. **Build:**
   ```sh
   go build -o snusfria
   ```

3. **Run:**
   ```sh
   ./snusfria
   ```

---

## Usage

Initialize your quit journey:
```sh
./snusfria init --price 49.90 --dosas 1 --pouches 24
```

Show your current status:
```sh
./snusfria status
```

Log a mood or craving:
```sh
./snusfria mood "Craving after lunch but resisted"
```

Log a relapse:
```sh
./snusfria relapse
```

Log a purchase:
```sh
./snusfria purchase
```

Set a savings goal:
```sh
./snusfria goal --amount 1000 --for "New headphones"
```

Show your health recovery timeline:
```sh
./snusfria timeline
```

Show your clean streaks:
```sh
./snusfria streak
```

See all logs and history:
```sh
./snusfria logs
./snusfria history
```

---

## Data Storage

All your data is stored locally in `~/.snusfria/data.json`.

---

## Requirements

- Go 1.19+ (see `go.mod` for details)
- Linux, macOS, or Windows (terminal support)

---

## Credits

- [fatih/color](https://github.com/fatih/color) for colored terminal output
- [spf13/cobra](https://github.com/spf13/cobra) for CLI framework

---

## License

MIT License

---

Stay strong and snus-free! 💪
