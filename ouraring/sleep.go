package ouraring

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

type SleepService service

type Sleep struct {
	SummaryDate       string    `json:"summary_date"`
	PeriodID          int       `json:"period_id"`
	IsLongest         int       `json:"is_longest"`
	Timezone          int       `json:"timezone"`
	BedtimeStart      time.Time `json:"bedtime_start"`
	BedtimeEnd        time.Time `json:"bedtime_end"`
	Score             int       `json:"score"`
	ScoreTotal        int       `json:"score_total"`
	ScoreDisturbances int       `json:"score_disturbances"`
	ScoreEfficiency   int       `json:"score_efficiency"`
	ScoreLatency      int       `json:"score_latency"`
	ScoreRem          int       `json:"score_rem"`
	ScoreDeep         int       `json:"score_deep"`
	ScoreAlignment    int       `json:"score_alignment"`
	Total             int       `json:"total"`
	Duration          int       `json:"duration"`
	Awake             int       `json:"awake"`
	Light             int       `json:"light"`
	Rem               int       `json:"rem"`
	Deep              int       `json:"deep"`
	OnsetLatency      int       `json:"onset_latency"`
	Restless          int       `json:"restless"`
	Efficiency        int       `json:"efficiency"`
	MidpointTime      int       `json:"midpoint_time"`
	HrLowest          int       `json:"hr_lowest"`
	HrAverage         float64   `json:"hr_average"`
	Rmssd             int       `json:"rmssd"`
	BreathAverage     float64       `json:"breath_average"`
	TemperatureDelta  float64   `json:"temperature_delta"`
	Hypnogram5Min     string    `json:"hypnogram_5min"`
	Hr5Min            []int     `json:"hr_5min"`
	Rmssd5Min         []int     `json:"rmssd_5min"`
}

type SleepResponse struct {
	Sleep []Sleep `json:"sleep"`
}


func (s *SleepService) List(ctx context.Context, opts ListOptions)([]Sleep, *http.Response, error){
    p := fmt.Sprintf("v1/sleep")
    u, err := url.Parse(p)
    if err != nil {
        return nil, nil, err
    }

    qs, err := query.Values(opts)
    if err != nil {
        return nil, nil, err
    }

    u.RawQuery = qs.Encode()

    req, err := s.client.NewRequest(u.String(), nil)
    if err != nil {
        return nil, nil, err
    }

    var sleepResp SleepResponse

    resp, err := s.client.Do(ctx, req, &sleepResp)
    if err != nil {
        return nil, resp, err
    }

    return sleepResp.Sleep, resp, nil
}
