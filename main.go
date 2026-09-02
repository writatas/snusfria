package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snusfria",
	Short: color.CyanString("🚭 Snusfria - Your snus quitting companion"),
	Long:  color.CyanString("Track your journey to a snus-free life. Stay strong! 💪"),
}

func main() {
	rootCmd.AddCommand(
		initCmd(),
		statusCmd(),
		logMoodCmd(),
		relapseCmd(),
		purchaseCmd(),
		setGoalCmd(),
		historyCmd(),
		timelineCmd(),
		logsCmd(),
		streakCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		os.Exit(1)
	}
}

func initCmd() *cobra.Command {
	var quitDateStr string
	var dosasPerDay float64
	var dosaPrice float64
	var pouchesPerDosa int

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize your quit journey",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, _ := LoadData()

			var quitDate time.Time
			if quitDateStr == "" || quitDateStr == "today" {
				quitDate = time.Now()
			} else {
				var err error
				quitDate, err = time.Parse("2006-01-02", quitDateStr)
				if err != nil {
					return fmt.Errorf("invalid date format, use YYYY-MM-DD")
				}
			}

			data.Config.QuitDate = quitDate
			data.Config.DosasPerDay = dosasPerDay
			data.Config.DosaPrice = dosaPrice
			data.Config.PouchesPerDosa = pouchesPerDosa

			if err := SaveData(data); err != nil {
				return err
			}

			perPouch := dosaPrice / float64(pouchesPerDosa)
			dailyCost := dosasPerDay * dosaPrice

			fmt.Println(color.GreenString("✅ Quit journey initialized!"))
			fmt.Printf("  %s %s\n", color.YellowString("Quit date:"), quitDate.Format("2006-01-02"))
			fmt.Printf("  %s %.2f kr/day\n", color.YellowString("Daily habit cost:"), dailyCost)
			fmt.Printf("  %s %.2f kr/pouch\n", color.YellowString("Per pouch cost:"), perPouch)
			fmt.Println(color.MagentaString("\n  Every day clean = %.2f kr saved 💰", dailyCost))
			return nil
		},
	}

	cmd.Flags().StringVarP(&quitDateStr, "date", "d", "today", "Quit start date (YYYY-MM-DD or 'today')")
	cmd.Flags().Float64VarP(&dosasPerDay, "dosas", "n", 1.0, "Dosas consumed per day (habit baseline)")
	cmd.Flags().Float64VarP(&dosaPrice, "price", "p", 0, "Price per dosa (kr)")
	cmd.Flags().IntVarP(&pouchesPerDosa, "pouches", "c", 24, "Pouches per dosa")
	cmd.MarkFlagRequired("price")

	return cmd
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show your current quit status",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil || data.Config.QuitDate.IsZero() {
				fmt.Println(color.RedString("❌ Not initialized. Run: snusfria init --price <price>"))
				return nil
			}

			PrintStatus(data)
			return nil
		},
	}
}

func logMoodCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mood [message]",
		Short: "Log your mental state / craving",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil {
				return err
			}

			entry := LogEntry{
				Time:    time.Now(),
				Message: strings.Join(args, " "),
				Type:    "mood",
			}
			data.Logs = append(data.Logs, entry)

			if err := SaveData(data); err != nil {
				return err
			}

			fmt.Println(color.CyanString("📝 Mood logged: \"%s\"", entry.Message))
			fmt.Println(color.GreenString("   Logged at %s", entry.Time.Format("2006-01-02 15:04")))
			return nil
		},
	}
}

func relapseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "relapse",
		Short: "Log a relapse (pouches consumed)",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(os.Stdin)

			yellow.Print("\n  Add a note (optional, press Enter to skip): ")
			scanner.Scan()
			note := strings.TrimSpace(scanner.Text())
			if note == "" {
				note = "No note"
			}

			yellow.Print("  How many pouches did you take? ")
			scanner.Scan()
			pouches := 0
			if input := strings.TrimSpace(scanner.Text()); input != "" {
				if v, err := strconv.Atoi(input); err == nil {
					pouches = v
				}
			}

			entry := LogEntry{
				Time:    time.Now(),
				Message: note,
				Type:    "relapse",
				Pouches: pouches,
			}
			data.Relapses = append(data.Relapses, entry)

			if err := SaveData(data); err != nil {
				return err
			}

			fmt.Println()
			fmt.Println(color.RedString("💔 Relapse logged. It happens. Get back up."))
			fmt.Printf("  %s\n", color.YellowString("Note: %s", note))
			if pouches > 0 {
				fmt.Printf("  %s\n", color.YellowString("Pouches taken: %d", pouches))
			}
			fmt.Println(color.GreenString("  ✨ Stay strong — your streak keeps going! 💪"))

			PrintRelapseImpact(entry.Time)
			return nil
		},
	}
}

func purchaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "purchase",
		Short: "Log a dosa purchase (affects pouch cost trend and savings)",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(os.Stdin)

			price := data.Config.DosaPrice
			pouches := data.Config.PouchesPerDosa

			yellow.Printf("\n  Dosa price (default %.2f kr): ", price)
			scanner.Scan()
			if input := strings.TrimSpace(scanner.Text()); input != "" {
				if v, err := strconv.ParseFloat(input, 64); err == nil {
					price = v
				}
			}

			yellow.Printf("  Number of pouches (default %d): ", pouches)
			scanner.Scan()
			if input := strings.TrimSpace(scanner.Text()); input != "" {
				if v, err := strconv.Atoi(input); err == nil {
					pouches = v
				}
			}

			perPouch := price / float64(pouches)
			dailyPouches := data.Config.DosasPerDay * float64(data.Config.PouchesPerDosa)

			purchase := Purchase{
				Price:    price,
				Pouches:  pouches,
				BoughtAt: time.Now(),
				PerPouch: perPouch,
			}
			data.Purchases = append(data.Purchases, purchase)

			if err := SaveData(data); err != nil {
				return err
			}

			fmt.Println()
			red.Printf("  💸 Purchase logged: %.2f kr — %d pouches @ %.2f kr each\n", price, pouches, perPouch)
			if dailyPouches > 0 {
				white.Printf("  📊 At your baseline of %.0f pouches/day, this dosa lasts ~%.1f days\n",
					dailyPouches, float64(pouches)/dailyPouches)
			}

			baselinePerPouch := data.Config.DosaPrice / float64(data.Config.PouchesPerDosa)
			if perPouch > baselinePerPouch {
				red.Printf("  📈 Price up %.2f kr from baseline (%.2f kr/pouch)\n", perPouch-baselinePerPouch, baselinePerPouch)
			} else if perPouch < baselinePerPouch {
				green.Printf("  📉 Price down %.2f kr from baseline (%.2f kr/pouch)\n", baselinePerPouch-perPouch, baselinePerPouch)
			} else {
				yellow.Printf("  📊 Same as baseline price\n")
			}

			PrintSavings(data)
			return nil
		},
	}
}

func setGoalCmd() *cobra.Command {
	var amount float64
	var description string

	cmd := &cobra.Command{
		Use:   "goal",
		Short: "Set a savings goal",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil {
				return err
			}

			data.Config.Goal = amount
			data.Config.GoalDescription = description

			if err := SaveData(data); err != nil {
				return err
			}

			fmt.Println(color.GreenString("🎯 Goal set: %.2f kr", amount))
			if description != "" {
				fmt.Printf("  %s\n", color.CyanString("For: %s", description))
			}
			return nil
		},
	}

	cmd.Flags().Float64VarP(&amount, "amount", "a", 0, "Goal amount in kr")
	cmd.Flags().StringVarP(&description, "for", "f", "", "What is the goal for?")
	cmd.MarkFlagRequired("amount")
	return cmd
}

func historyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "Show mood and relapse history",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil {
				return err
			}
			PrintHistory(data)
			return nil
		},
	}
}

func timelineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "timeline",
		Short: "Show health recovery timeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil || data.Config.QuitDate.IsZero() {
				fmt.Println(color.RedString("❌ Not initialized. Run: snusfria init --price <price>"))
				return nil
			}
			PrintTimeline(data)
			return nil
		},
	}
}

func logsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logs",
		Short: "Show timestamped mood and relapse logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil {
				return err
			}
			PrintLogs(data)
			return nil
		},
	}
}

func streakCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "streak",
		Short: "Show your clean streaks across all attempts",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := LoadData()
			if err != nil || data.Config.QuitDate.IsZero() {
				fmt.Println(color.RedString("❌ Not initialized. Run: snusfria init --price <price>"))
				return nil
			}
			PrintStreaks(data)
			return nil
		},
	}
}

// suppress unused strconv
var _ = strconv.Itoa
