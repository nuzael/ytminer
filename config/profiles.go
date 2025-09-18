package config

import "fmt"

// WeightProfile represents a predefined set of weights for Opportunity Score
type WeightProfile struct {
	Name        string
	Description string
	Weights     OpportunityWeights
}

// OpportunityWeights contains all weight components
type OpportunityWeights struct {
	VPD   float64
	Like  float64
	Fresh float64
	Sat   float64
	Slope float64
}

// Predefined weight profiles for different strategies
var WeightProfiles = map[string]WeightProfile{
	"exploration": {
		Name:        "Exploration",
		Description: "Discover new niches and emerging trends",
		Weights: OpportunityWeights{
			VPD:   0.30, // Lower VPD weight
			Like:  0.20, // Lower engagement weight
			Fresh: 0.35, // High freshness (new content)
			Sat:   0.10, // Low saturation penalty
			Slope: 0.25, // High slope (trending up)
		},
	},
	"evergreen": {
		Name:        "Evergreen",
		Description: "Focus on timeless, high-quality content",
		Weights: OpportunityWeights{
			VPD:   0.25, // Lower VPD weight
			Like:  0.40, // High engagement weight
			Fresh: 0.10, // Low freshness weight
			Sat:   0.20, // Medium saturation penalty
			Slope: 0.05, // Low slope weight
		},
	},
	"trending": {
		Name:        "Trending",
		Description: "Catch viral content and momentum",
		Weights: OpportunityWeights{
			VPD:   0.50, // High VPD weight
			Like:  0.15, // Lower engagement weight
			Fresh: 0.20, // Medium freshness
			Sat:   0.05, // Low saturation penalty
			Slope: 0.30, // High slope weight
		},
	},
	"balanced": {
		Name:        "Balanced",
		Description: "Default balanced approach",
		Weights: OpportunityWeights{
			VPD:   0.45,
			Like:  0.25,
			Fresh: 0.20,
			Sat:   0.30,
			Slope: 0.15,
		},
	},
}

// GetProfile returns a weight profile by name, or default if not found
func GetProfile(name string) WeightProfile {
	if profile, exists := WeightProfiles[name]; exists {
		return profile
	}
	return WeightProfiles["balanced"]
}

// ListProfiles returns all available profiles
func ListProfiles() []WeightProfile {
	profiles := make([]WeightProfile, 0, len(WeightProfiles))
	for _, profile := range WeightProfiles {
		profiles = append(profiles, profile)
	}
	return profiles
}

// ApplyProfile applies a weight profile to the config
func (c *AppConfig) ApplyProfile(profileName string) error {
	profile := GetProfile(profileName)
	c.OppWeightVPD = profile.Weights.VPD
	c.OppWeightLike = profile.Weights.Like
	c.OppWeightFresh = profile.Weights.Fresh
	c.OppWeightSatPen = profile.Weights.Sat
	c.OppWeightSlope = profile.Weights.Slope
	return nil
}

// GetActiveProfileName returns the name of the currently active profile
func (c *AppConfig) GetActiveProfileName() string {
	weights := OpportunityWeights{
		VPD:   c.OppWeightVPD,
		Like:  c.OppWeightLike,
		Fresh: c.OppWeightFresh,
		Sat:   c.OppWeightSatPen,
		Slope: c.OppWeightSlope,
	}

	for name, profile := range WeightProfiles {
		if profile.Weights == weights {
			return name
		}
	}
	return "custom"
}

// DisplayProfiles prints all available profiles
func DisplayProfiles() {
	fmt.Println("Available Weight Profiles:")
	fmt.Println("========================")
	for name, profile := range WeightProfiles {
		fmt.Printf("%s: %s\n", name, profile.Description)
		fmt.Printf("  Weights: VPD=%.2f, LIKE=%.2f, FRESH=%.2f, SAT=%.2f, SLOPE=%.2f\n",
			profile.Weights.VPD, profile.Weights.Like, profile.Weights.Fresh,
			profile.Weights.Sat, profile.Weights.Slope)
		fmt.Println()
	}
}
