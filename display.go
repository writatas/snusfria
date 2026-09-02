package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	bold    = color.New(color.Bold)
	green   = color.New(color.FgGreen, color.Bold)
	red     = color.New(color.FgRed, color.Bold)
	yellow  = color.New(color.FgYellow, color.Bold)
	cyan    = color.New(color.FgCyan, color.Bold)
	magenta = color.New(color.FgMagenta, color.Bold)
	white   = color.New(color.FgWhite, color.Bold)
	blue    = color.New(color.FgBlue, color.Bold)
)

const relapseImpactThreshold = 30 * 24 * time.Hour

func PrintStatus(data *Data) {
	cfg := data.Config
	now := time.Now()
	elapsed := now.Sub(cfg.QuitDate)
	days := elapsed.Hours() / 24

	fmt.Println()
	cyan.Println("╔══════════════════════════════════════════════╗")
	cyan.Println("║         🚭  SNUSFRIA STATUS  🚭              ║")
	cyan.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()

	white.Printf("  🏁 Attempt started: ")
	green.Printf("%s\n", cfg.QuitDate.Format("2006-01-02 15:04"))

	d := int(days)
	h := int(elapsed.Hours()) % 24
	m := int(elapsed.Minutes()) % 60

	white.Printf("  📅 Time clean:      ")
	if d > 0 {
		green.Printf("%d days, %d hours, %d minutes\n", d, h, m)
	} else {
		yellow.Printf("%d hours, %d minutes\n", h, m)
	}

	white.Printf("  💔 Relapses:        ")
	if len(data.Relapses) == 0 {
		green.Println("None! Keep it up! 🎉")
	} else {
		red.Printf("%d relapse(s)\n", len(data.Relapses))
	}

	// Compact streak summary
	periods := computeStreaks(data)
	longest := 0.0
	for _, p := range periods {
		if p.days > longest {
			longest = p.days
		}
	}
	current := periods[len(periods)-1].days
	white.Printf("  🔥 Current streak:  ")
	if current >= longest {
		green.Printf("%.1f days  ", current)
		green.Printf("🏆 personal best!\n")
	} else {
		cyan.Printf("%.1f days  ", current)
		white.Printf("(best: ")
		green.Printf("%.1f days", longest)
		white.Printf(")\n")
	}

	fmt.Println()
	PrintSavings(data)
	fmt.Println()
	PrintPouchStats(data)
	fmt.Println()
	PrintNextMilestone(data)

	// Show body impact of most recent relapse if within threshold
	if len(data.Relapses) > 0 {
		latest := data.Relapses[0]
		for _, r := range data.Relapses {
			if r.Time.After(latest.Time) {
				latest = r
			}
		}
		if time.Since(latest.Time) <= relapseImpactThreshold {
			PrintRelapseImpact(latest.Time)
		}
	}
}

func pouchesConsumed(data *Data) int {
	total := 0
	cfg := data.Config
	for _, r := range data.Relapses {
		if r.Time.After(cfg.QuitDate) {
			total += r.Pouches
		}
	}
	return total
}

func PrintPouchStats(data *Data) {
	cfg := data.Config
	if cfg.PouchesPerDosa == 0 || cfg.DosasPerDay == 0 {
		return
	}

	elapsed := time.Since(cfg.QuitDate)
	days := elapsed.Hours() / 24

	dailyPouches := cfg.DosasPerDay * float64(cfg.PouchesPerDosa)
	yearlyPouches := dailyPouches * 365
	wouldHaveConsumed := dailyPouches * days

	consumed := pouchesConsumed(data)
	avoided := wouldHaveConsumed - float64(consumed)
	if avoided < 0 {
		avoided = 0
	}

	var ratio float64
	if wouldHaveConsumed > 0 {
		ratio = avoided / wouldHaveConsumed
	}

	cyan.Println("  ─────────────── 🚫 POUCHES AVOIDED ───────────────")
	white.Printf("  Yearly baseline:     ")
	yellow.Printf("%.0f pouches/year\n", yearlyPouches)
	white.Printf("  Would have taken:    ")
	yellow.Printf("%.0f pouches\n", wouldHaveConsumed)
	white.Printf("  Actually consumed:   ")
	if consumed > 0 {
		red.Printf("%d pouches\n", consumed)
	} else {
		green.Printf("0 pouches 🎉\n")
	}
	white.Printf("  Pouches avoided:     ")
	if ratio >= 0.9 {
		green.Printf("%.0f 🌟\n", avoided)
	} else if ratio >= 0.5 {
		yellow.Printf("%.0f\n", avoided)
	} else {
		red.Printf("%.0f\n", avoided)
	}

	fmt.Println()
	PrintPouchBar(avoided, wouldHaveConsumed, ratio)
}

func PrintPouchBar(avoided, wouldHave, ratio float64) {
	pct := ratio
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}

	barWidth := 40
	filled := int(math.Round(pct * float64(barWidth)))
	empty := barWidth - filled

	var barColor *color.Color
	switch {
	case pct >= 0.9:
		barColor = color.New(color.FgGreen, color.Bold)
	case pct >= 0.66:
		barColor = color.New(color.FgGreen)
	case pct >= 0.33:
		barColor = color.New(color.FgYellow)
	default:
		barColor = color.New(color.FgRed)
	}

	filledStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", empty)

	white.Printf("  🚫 Avoidance (%.0f / %.0f pouches)\n", avoided, wouldHave)
	white.Print("  [")
	barColor.Print(filledStr)
	color.New(color.FgWhite).Print(emptyStr)
	white.Printf("] ")

	displayPct := pct * 100
	if displayPct >= 100 {
		green.Printf("%.0f%% 🏆 PERFECT!\n", displayPct)
	} else {
		barColor.Printf("%.1f%%\n", displayPct)
	}
}

func PrintSavings(data *Data) {
	cfg := data.Config
	if cfg.DosaPrice == 0 {
		return
	}

	now := time.Now()
	elapsed := now.Sub(cfg.QuitDate)
	days := elapsed.Hours() / 24

	dailyCost := cfg.DosasPerDay * cfg.DosaPrice
	grossSavings := days * dailyCost

	var purchaseCost float64
	var purchasePouches int
	for _, p := range data.Purchases {
		if p.BoughtAt.After(cfg.QuitDate) {
			purchaseCost += p.Price
			purchasePouches += p.Pouches
		}
	}

	netSavings := grossSavings - purchaseCost

	baselinePerPouch := cfg.DosaPrice / float64(cfg.PouchesPerDosa)
	avgPerPouch := baselinePerPouch
	if purchasePouches > 0 {
		avgPerPouch = purchaseCost / float64(purchasePouches)
	}

	cyan.Println("  ─────────────── 💰 SAVINGS ───────────────")
	white.Printf("  Daily habit cost:    ")
	yellow.Printf("%.2f kr/day\n", dailyCost)
	white.Printf("  Baseline pouch cost: ")
	yellow.Printf("%.2f kr/pouch\n", baselinePerPouch)
	if purchasePouches > 0 {
		white.Printf("  Avg pouch cost now:  ")
		if avgPerPouch > baselinePerPouch {
			red.Printf("%.2f kr/pouch ↑\n", avgPerPouch)
		} else {
			green.Printf("%.2f kr/pouch ↓\n", avgPerPouch)
		}
	}
	white.Printf("  Gross saved:         ")
	green.Printf("%.2f kr\n", grossSavings)

	if purchaseCost > 0 {
		white.Printf("  Spent on purchases:  ")
		red.Printf("-%.2f kr\n", purchaseCost)
		white.Printf("  Net savings:         ")
		if netSavings >= 0 {
			green.Printf("%.2f kr ✨\n", netSavings)
		} else {
			red.Printf("%.2f kr\n", netSavings)
		}
	} else {
		white.Printf("  Net savings:         ")
		green.Printf("%.2f kr ✨\n", netSavings)
	}

	if cfg.Goal > 0 {
		fmt.Println()
		PrintGoalBar(netSavings, cfg.Goal, cfg.GoalDescription)
	}

	fmt.Println()
	PrintPouchCostTrend(data)
}

func PrintPouchCostTrend(data *Data) {
	cfg := data.Config
	if cfg.DosaPrice == 0 || cfg.PouchesPerDosa == 0 {
		return
	}

	type point struct {
		label    string
		perPouch float64
	}

	baseline := cfg.DosaPrice / float64(cfg.PouchesPerDosa)
	points := []point{{"init", baseline}}

	purchases := make([]Purchase, 0, len(data.Purchases))
	for _, p := range data.Purchases {
		if p.BoughtAt.After(cfg.QuitDate) {
			purchases = append(purchases, p)
		}
	}
	sort.Slice(purchases, func(i, j int) bool {
		return purchases[i].BoughtAt.Before(purchases[j].BoughtAt)
	})
	for _, p := range purchases {
		points = append(points, point{p.BoughtAt.Format("01-02"), p.PerPouch})
	}

	if len(points) < 2 {
		return
	}

	minVal, maxVal := points[0].perPouch, points[0].perPouch
	for _, pt := range points {
		if pt.perPouch < minVal {
			minVal = pt.perPouch
		}
		if pt.perPouch > maxVal {
			maxVal = pt.perPouch
		}
	}

	sparkChars := []rune("▁▂▃▄▅▆▇█")
	sparkRange := maxVal - minVal

	cyan.Println("  ─────────────── 📈 POUCH COST TREND ───────────────")
	white.Printf("  ")

	for _, pt := range points {
		var idx int
		if sparkRange == 0 {
			idx = 3
		} else {
			idx = int(math.Round((pt.perPouch - minVal) / sparkRange * float64(len(sparkChars)-1)))
		}
		ch := string(sparkChars[idx])

		if pt.perPouch > baseline {
			red.Printf("%s", ch)
		} else if pt.perPouch < baseline {
			green.Printf("%s", ch)
		} else {
			yellow.Printf("%s", ch)
		}
	}

	fmt.Println()

	first := points[0]
	last := points[len(points)-1]
	white.Printf("  init: ")
	yellow.Printf("%.2f kr", first.perPouch)
	white.Printf("  →  latest: ")
	if last.perPouch > first.perPouch {
		red.Printf("%.2f kr ↑  (+%.2f kr/pouch)\n", last.perPouch, last.perPouch-first.perPouch)
	} else if last.perPouch < first.perPouch {
		green.Printf("%.2f kr ↓  (-%.2f kr/pouch)\n", last.perPouch, first.perPouch-last.perPouch)
	} else {
		yellow.Printf("%.2f kr  (no change)\n", last.perPouch)
	}
	white.Printf("  min: ")
	green.Printf("%.2f kr", minVal)
	white.Printf("  max: ")
	red.Printf("%.2f kr\n", maxVal)
}

func PrintGoalBar(current, goal float64, description string) {
	percent := current / goal
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	barWidth := 40
	filled := int(math.Round(percent * float64(barWidth)))
	empty := barWidth - filled

	var barColor *color.Color
	switch {
	case percent >= 1.0:
		barColor = color.New(color.FgGreen, color.Bold)
	case percent >= 0.66:
		barColor = color.New(color.FgGreen)
	case percent >= 0.33:
		barColor = color.New(color.FgYellow)
	default:
		barColor = color.New(color.FgRed)
	}

	filledStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", empty)

	label := "Goal"
	if description != "" {
		label = description
	}

	white.Printf("  🎯 %s (%.2f kr)\n", label, goal)
	white.Print("  [")
	barColor.Print(filledStr)
	color.New(color.FgWhite).Print(emptyStr)
	white.Printf("] ")

	pct := percent * 100
	if pct >= 100 {
		green.Printf("%.0f%% 🎉 GOAL REACHED!\n", pct)
	} else {
		barColor.Printf("%.1f%%\n", pct)
		remaining := goal - current
		if remaining > 0 {
			white.Printf("  💡 %.2f kr to go!\n", remaining)
		}
	}
}

func PrintNextMilestone(data *Data) {
	cfg := data.Config
	now := time.Now()
	elapsed := now.Sub(cfg.QuitDate)

	cyan.Println("  ─────────────── 🏥 HEALTH TIMELINE ───────────────")
	fmt.Println()

	var nextUnlocked *Milestone
	for i := range healthMilestones {
		m := &healthMilestones[i]
		reached := elapsed >= m.Duration
		if reached {
			green.Printf("  %s ✅ %s\n", m.Emoji, m.Label)
			green.Printf("     %s\n\n", m.Description)
		} else {
			if nextUnlocked == nil {
				nextUnlocked = m
				timeLeft := m.Duration - elapsed
				daysLeft := timeLeft.Hours() / 24
				yellow.Printf("  %s ⏳ %s  (in ", m.Emoji, m.Label)
				if daysLeft >= 1 {
					yellow.Printf("%.0f days)\n", daysLeft)
				} else {
					yellow.Printf("%.0f hours)\n", timeLeft.Hours())
				}
				white.Printf("     %s\n\n", m.Description)
			}
		}
	}
}

func PrintTimeline(data *Data) {
	fmt.Println()
	magenta.Println("╔══════════════════════════════════════════════════╗")
	magenta.Println("║       🏥  FULL HEALTH RECOVERY TIMELINE  🏥     ║")
	magenta.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	cfg := data.Config
	elapsed := time.Since(cfg.QuitDate)

	for _, m := range healthMilestones {
		reached := elapsed >= m.Duration
		if reached {
			green.Printf("  %s ✅ ", m.Emoji)
			green.Printf("[ACHIEVED] %s\n", m.Label)
		} else {
			timeLeft := m.Duration - elapsed
			daysLeft := timeLeft.Hours() / 24
			yellow.Printf("  %s ⏳ ", m.Emoji)
			if daysLeft >= 1 {
				yellow.Printf("[%.0f days away] %s\n", daysLeft, m.Label)
			} else {
				yellow.Printf("[%.0f hours away] %s\n", timeLeft.Hours(), m.Label)
			}
		}
		white.Printf("     %s\n\n", m.Description)
	}
}

func PrintRelapseImpact(relapseTime time.Time) {
	now := time.Now()
	elapsed := now.Sub(relapseTime)

	fmt.Println()
	red.Println("╔══════════════════════════════════════════════════╗")
	red.Println("║       ⚠️   WHAT'S HAPPENING IN YOUR BODY  ⚠️    ║")
	red.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	white.Printf("  Since your relapse at %s:\n\n", relapseTime.Format("2006-01-02 15:04"))

	for _, e := range relapseEffects {
		reached := elapsed >= e.Duration
		if reached {
			red.Printf("  %s ⚠️  %s\n", e.Emoji, e.Label)
			white.Printf("     %s\n\n", e.Description)
		} else {
			timeLeft := e.Duration - elapsed
			var timeStr string
			if timeLeft.Hours() >= 1 {
				timeStr = fmt.Sprintf("%.0f hours", timeLeft.Hours())
			} else {
				timeStr = fmt.Sprintf("%.0f minutes", timeLeft.Minutes())
			}
			yellow.Printf("  %s ⏳ %s  (in %s)\n", e.Emoji, e.Label, timeStr)
			white.Printf("     %s\n\n", e.Description)
			break
		}
	}

	fmt.Println()
	cyan.Println("  💡 Every minute without another pouch starts reversing this.")
	fmt.Println()
}

func PrintHistory(data *Data) {
	fmt.Println()
	cyan.Println("╔══════════════════════════════════════════════╗")
	cyan.Println("║              📖  HISTORY                    ║")
	cyan.Println("╚══════════════════════════════════════════════╝")

	if len(data.Logs) == 0 && len(data.Relapses) == 0 && len(data.Purchases) == 0 {
		yellow.Println("\n  No entries yet.")
		return
	}

	if len(data.Logs) > 0 {
		fmt.Println()
		cyan.Println("  📝 MOOD / CRAVING LOGS")
		for _, l := range data.Logs {
			white.Printf("  [%s] ", l.Time.Format("2006-01-02 15:04"))
			fmt.Printf("%s\n", l.Message)
		}
	}

	if len(data.Relapses) > 0 {
		fmt.Println()
		red.Println("  💔 RELAPSES")
		for _, r := range data.Relapses {
			red.Printf("  [%s] ", r.Time.Format("2006-01-02 15:04"))
			if r.Pouches > 0 {
				fmt.Printf("%s — %d pouches\n", r.Message, r.Pouches)
			} else {
				fmt.Printf("%s\n", r.Message)
			}
		}
	}

	if len(data.Purchases) > 0 {
		fmt.Println()
		yellow.Println("  💸 PURCHASES")
		for _, p := range data.Purchases {
			yellow.Printf("  [%s] ", p.BoughtAt.Format("2006-01-02 15:04"))
			fmt.Printf("%.2f kr — %d pouches @ %.2f kr each\n", p.Price, p.Pouches, p.PerPouch)
		}
	}

	fmt.Println()
}

func PrintLogs(data *Data) {
	type entry struct {
		t       time.Time
		kind    string
		message string
		pouches int
	}

	var entries []entry
	for _, l := range data.Logs {
		entries = append(entries, entry{l.Time, "mood", l.Message, 0})
	}
	for _, r := range data.Relapses {
		entries = append(entries, entry{r.Time, "relapse", r.Message, r.Pouches})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].t.Before(entries[j].t)
	})

	fmt.Println()
	cyan.Println("╔══════════════════════════════════════════════╗")
	cyan.Println("║         📋  MOOD & RELAPSE LOGS             ║")
	cyan.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()

	if len(entries) == 0 {
		yellow.Println("  No logs yet.")
		fmt.Println()
		return
	}

	for _, e := range entries {
		ts := white.Sprintf("[%s]", e.t.Format("2006-01-02 15:04"))
		switch e.kind {
		case "mood":
			fmt.Printf("  %s ", ts)
			cyan.Printf("📝 mood    ")
			fmt.Printf("%s\n", e.message)
		case "relapse":
			fmt.Printf("  %s ", ts)
			red.Printf("💔 relapse ")
			if e.pouches > 0 {
				fmt.Printf("%s — %d pouches\n", e.message, e.pouches)
			} else {
				fmt.Printf("%s\n", e.message)
			}
		}
	}
	fmt.Println()
}

type streakPeriod struct {
	start     time.Time
	end       time.Time
	days      float64
	isCurrent bool
}

func computeStreaks(data *Data) []streakPeriod {
	var relapses []time.Time
	for _, r := range data.Relapses {
		relapses = append(relapses, r.Time)
	}
	sort.Slice(relapses, func(i, j int) bool {
		return relapses[i].Before(relapses[j])
	})

	var periods []streakPeriod
	start := data.Config.QuitDate
	now := time.Now()

	for _, r := range relapses {
		if r.Before(start) {
			continue
		}
		d := r.Sub(start).Hours() / 24
		periods = append(periods, streakPeriod{
			start: start,
			end:   r,
			days:  d,
		})
		start = r
	}

	d := now.Sub(start).Hours() / 24
	periods = append(periods, streakPeriod{
		start:     start,
		end:       now,
		days:      d,
		isCurrent: true,
	})

	return periods
}

func PrintStreaks(data *Data) {
	periods := computeStreaks(data)

	fmt.Println()
	cyan.Println("╔══════════════════════════════════════════════╗")
	cyan.Println("║           🔥  CLEAN STREAKS                 ║")
	cyan.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()

	longest := 0.0
	longestIdx := 0
	for i, p := range periods {
		if p.days > longest {
			longest = p.days
			longestIdx = i
		}
	}

	maxDays := longest
	if maxDays == 0 {
		maxDays = 1
	}

	for i, p := range periods {
		isLongest := i == longestIdx
		label := fmt.Sprintf("Attempt %d", i+1)
		if p.isCurrent {
			label = "Current  "
		}

		ratio := p.days / maxDays
		barWidth := 30
		filled := int(math.Round(ratio * float64(barWidth)))
		if filled < 1 && p.days > 0 {
			filled = 1
		}
		empty := barWidth - filled

		var barColor *color.Color
		if p.isCurrent {
			barColor = color.New(color.FgCyan, color.Bold)
		} else if isLongest {
			barColor = color.New(color.FgGreen, color.Bold)
		} else {
			barColor = color.New(color.FgYellow)
		}

		filledStr := strings.Repeat("█", filled)
		emptyStr := strings.Repeat("░", empty)

		suffix := ""
		if isLongest && !p.isCurrent {
			suffix = " 🏆 best"
		} else if p.isCurrent && isLongest {
			suffix = " 🏆 best & current"
		} else if p.isCurrent {
			suffix = " ← now"
		}

		white.Printf("  %s  [", label)
		barColor.Print(filledStr)
		color.New(color.FgWhite).Print(emptyStr)
		white.Printf("] ")
		barColor.Printf("%.1f days", p.days)
		if suffix != "" {
			green.Printf("%s", suffix)
		}
		fmt.Println()

		if !p.isCurrent {
			white.Printf("             %s → %s\n",
				p.start.Format("2006-01-02"),
				p.end.Format("2006-01-02"))
		} else {
			white.Printf("             %s → now\n", p.start.Format("2006-01-02"))
		}
		fmt.Println()
	}

	cyan.Println("  ──────────────────────────────────────────")
	white.Printf("  🏆 Longest streak:  ")
	green.Printf("%.1f days\n", longest)
	white.Printf("  🔥 Current streak:  ")
	current := periods[len(periods)-1].days
	if current >= longest {
		green.Printf("%.1f days 🌟\n", current)
	} else {
		cyan.Printf("%.1f days\n", current)
	}
	fmt.Println()
}
