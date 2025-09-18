package score

import (
	"fmt"
	"math"
	"sort"
	"time"

	"ytminer/domain/metrics"
	"ytminer/utils"
	"ytminer/youtube"
)

// Weights defines the contributions of each factor to the Opportunity Score.
type Weights struct {
	VPD   float64
	Like  float64
	Fresh float64
	Sat   float64 // penalty weight
	Slope float64
}

// Item is the result for a single video.
type Item struct {
	Title      string
	Channel    string
	URL        string
	Score      float64
	VPD        float64
	VPD7       float64
	VPD30      float64
	Slope      float64
	LikeRate   float64
	AgeDays    int
	Saturation float64
	Reasons    []string
}

// Compute computes a ranking of items using only the provided inputs.
// Pure function: no I/O, deterministic given (videos, weights, now).
func Compute(videos []youtube.Video, w Weights, now time.Time) []Item {
	if len(videos) == 0 {
		return []Item{}
	}

	vpdVals := make([]float64, 0, len(videos))
	likeRateVals := make([]float64, 0, len(videos))
	ageVals := make([]float64, 0, len(videos))
	slopeVals := make([]float64, 0, len(videos))

	tokenFreq := make(map[string]int)
	videoPrimaryToken := make([]string, len(videos))

	for i, v := range videos {
		vpd := v.VPD
		if vpd == 0 {
			vpd = metrics.VPD(v.Views, v.PublishedAt, now)
		}
		vpdVals = append(vpdVals, vpd)
		likeRate := metrics.LikeRatePerThousand(v.Views, v.Likes)
		likeRateVals = append(likeRateVals, likeRate)

		ageDays := metrics.AgeDays(now, v.PublishedAt)
		ageVals = append(ageVals, float64(ageDays))

		v7 := metrics.VPDWindow(v.Views, v.PublishedAt, now, 7)
		v30 := metrics.VPDWindow(v.Views, v.PublishedAt, now, 30)
		slope := metrics.SlopeVPD(v7, v30)
		slopeVals = append(slopeVals, slope)

		toks := utils.Tokenize(v.Title)
		primary := ""
		for _, t := range toks {
			if t == "" {
				continue
			}
			primary = t
			break
		}
		if primary == "" {
			primary = "_na_"
		}
		videoPrimaryToken[i] = primary
		tokenFreq[primary]++
	}

	items := make([]Item, 0, len(videos))
	for i, v := range videos {
		likeRate := metrics.LikeRatePerThousand(v.Views, v.Likes)
		ageDays := metrics.AgeDays(now, v.PublishedAt)
		freshness := metrics.FreshnessFromAges(ageVals, ageDays)
		saturation := metrics.NormalizeSaturation(tokenFreq[videoPrimaryToken[i]], len(videos))
		vpd := v.VPD
		if vpd == 0 {
			vpd = metrics.VPD(v.Views, v.PublishedAt, now)
		}
		v7 := metrics.VPDWindow(v.Views, v.PublishedAt, now, 7)
		v30 := metrics.VPDWindow(v.Views, v.PublishedAt, now, 30)
		slope := metrics.SlopeVPD(v7, v30)

		score := w.VPD*z(vpdVals, vpd) + w.Like*z(likeRateVals, likeRate) + w.Fresh*freshness - w.Sat*saturation + w.Slope*z(slopeVals, slope)

		reasons := []string{
			fmt.Sprintf("VPD=%s", utils.FormatVPD(vpd)),
			fmt.Sprintf("VPD7=%s", utils.FormatVPD(v7)),
			fmt.Sprintf("VPD30=%s", utils.FormatVPD(v30)),
			fmt.Sprintf("Slope=%.2f", slope),
			fmt.Sprintf("LikeRate=%.2f/1k", likeRate),
			fmt.Sprintf("Age=%dd", ageDays),
			fmt.Sprintf("Saturation=%.2f", saturation),
		}

		items = append(items, Item{
			Title:      v.Title,
			Channel:    v.Channel,
			URL:        v.URL,
			Score:      score,
			VPD:        vpd,
			VPD7:       v7,
			VPD30:      v30,
			Slope:      slope,
			LikeRate:   likeRate,
			AgeDays:    ageDays,
			Saturation: saturation,
			Reasons:    reasons,
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
	if len(items) > 20 {
		items = items[:20]
	}
	return items
}

func mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func stddev(vals []float64, m float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		dv := v - m
		sum += dv * dv
	}
	return math.Sqrt(sum / float64(len(vals)))
}

func z(vals []float64, x float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	m := mean(vals)
	sd := stddev(vals, m)
	if sd == 0 {
		return 0
	}
	return (x - m) / sd
}

func minMax(vals []float64, x float64, invert bool) float64 {
	if len(vals) == 0 {
		return 0
	}
	minV, maxV := vals[0], vals[0]
	for _, v := range vals {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if maxV == minV {
		return 0
	}
	norm := (x - minV) / (maxV - minV)
	if invert {
		return 1.0 - norm
	}
	return norm
}
