package main

import "time"

type Milestone struct {
	Duration    time.Duration
	Label       string
	Description string
	Emoji       string
}

type RelapseEffect struct {
	Duration    time.Duration
	Label       string
	Description string
	Emoji       string
}

var healthMilestones = []Milestone{
	{
		Duration:    24 * time.Hour,
		Label:       "24 hours",
		Description: "Blood pressure and heart rate begin to normalize. Nicotine starts leaving your bloodstream.",
		Emoji:       "❤️",
	},
	{
		Duration:    3 * 24 * time.Hour,
		Label:       "3 days",
		Description: "Nicotine is fully cleared from your body. Withdrawal symptoms peak — headaches, irritability. Push through!",
		Emoji:       "🔥",
	},
	{
		Duration:    7 * 24 * time.Hour,
		Label:       "1 week",
		Description: "Taste and smell sharpen noticeably. Gum inflammation starts reducing. Mouth feels cleaner.",
		Emoji:       "👅",
	},
	{
		Duration:    14 * 24 * time.Hour,
		Label:       "2 weeks",
		Description: "Blood circulation improving. Energy levels rising. Gum tissue beginning to regenerate.",
		Emoji:       "🩸",
	},
	{
		Duration:    30 * 24 * time.Hour,
		Label:       "1 month",
		Description: "Gum health significantly improved. Reduced risk of gum recession. Stomach lining recovering from nicotine irritation.",
		Emoji:       "🦷",
	},
	{
		Duration:    90 * 24 * time.Hour,
		Label:       "3 months",
		Description: "Mouth sores healing. Leukoplakia (white patches from snus) begin to fade. Reduced pre-cancerous cell activity.",
		Emoji:       "🌿",
	},
	{
		Duration:    180 * 24 * time.Hour,
		Label:       "6 months",
		Description: "Significant reduction in gum disease risk. Pancreatic enzyme production normalizing. Stomach acid levels balanced.",
		Emoji:       "💚",
	},
	{
		Duration:    365 * 24 * time.Hour,
		Label:       "1 year",
		Description: "Risk of oral cancer halved. Gum line stabilized. Risk of esophageal cancer reduced by ~30%.",
		Emoji:       "🏆",
	},
	{
		Duration:    2 * 365 * 24 * time.Hour,
		Label:       "2 years",
		Description: "Oral mucosal tissue largely normalized. Cardiovascular risk approaching that of a non-user.",
		Emoji:       "🫀",
	},
	{
		Duration:    5 * 365 * 24 * time.Hour,
		Label:       "5 years",
		Description: "Risk of stroke similar to a non-user. Gum recession halted. Tooth loss risk dramatically reduced.",
		Emoji:       "🧠",
	},
	{
		Duration:    10 * 365 * 24 * time.Hour,
		Label:       "10 years",
		Description: "Risk of pancreatic cancer significantly reduced, approaching non-user levels. Oral cancer risk minimal.",
		Emoji:       "🌟",
	},
}

var relapseEffects = []RelapseEffect{
	{
		Duration:    10 * time.Second,
		Label:       "Within seconds",
		Description: "Nicotine floods your bloodstream. Dopamine spikes — this is the trap your brain remembers.",
		Emoji:       "⚡",
	},
	{
		Duration:    30 * time.Minute,
		Label:       "30 minutes",
		Description: "Blood pressure and heart rate elevated. Adrenaline released, reducing insulin effectiveness.",
		Emoji:       "💓",
	},
	{
		Duration:    2 * time.Hour,
		Label:       "2 hours",
		Description: "Nicotine levels drop sharply. Your brain signals craving — the cycle tries to restart itself.",
		Emoji:       "📉",
	},
	{
		Duration:    6 * time.Hour,
		Label:       "6 hours",
		Description: "Blood flow to gums reduced. Saliva production altered, leaving mouth drier and more vulnerable.",
		Emoji:       "🦷",
	},
	{
		Duration:    24 * time.Hour,
		Label:       "24 hours",
		Description: "Nicotine still in your system. Sleep may be disrupted. Mood dips as dopamine baseline drops.",
		Emoji:       "😔",
	},
	{
		Duration:    3 * 24 * time.Hour,
		Label:       "3 days",
		Description: "Nicotine fully cleared again, but withdrawal restarts. Headaches, irritability, difficulty concentrating.",
		Emoji:       "🤕",
	},
	{
		Duration:    7 * 24 * time.Hour,
		Label:       "1 week",
		Description: "Gum tissue inflammation rises again. Any healing progress from before starts to reverse.",
		Emoji:       "🔴",
	},
	{
		Duration:    14 * 24 * time.Hour,
		Label:       "2 weeks",
		Description: "Continued use re-exposes oral tissue to carcinogens. Cell regeneration slows again.",
		Emoji:       "⚠️",
	},
	{
		Duration:    30 * 24 * time.Hour,
		Label:       "1 month",
		Description: "Risk of gum recession increases again. Nicotine dependence re-establishing at a neurological level.",
		Emoji:       "🧠",
	},
}
