package kindle

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"

	pkgnews "assistant/pkg/news"
)

var (
	chartGPSC []chartPoint
	chartSH   []chartPoint
	chartGC   []chartPoint
	chartTime time.Time
	chartMu   sync.Mutex

	weatherData string
	weatherTime time.Time
	weatherMu   sync.Mutex
)

type Service struct {
	news *pkgnews.Collector
}

func NewService() *Service {
	return &Service{news: pkgnews.New()}
}

type weatherResponse struct {
	CurrentCondition []struct {
		TempC       string `json:"temp_C"`
		Humidity    string `json:"humidity"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`
}

const weatherURL = "https://wttr.in?format=j1"

func (s *Service) FetchWeather() string {
	weatherMu.Lock()
	defer weatherMu.Unlock()
	if weatherData != "" && time.Since(weatherTime) < 10*time.Minute {
		return weatherData
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(weatherURL)
	if err != nil {
		return weatherData
	}
	defer resp.Body.Close()
	var w weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return weatherData
	}
	if len(w.CurrentCondition) > 0 {
		cc := w.CurrentCondition[0]
		desc := ""
		if len(cc.WeatherDesc) > 0 {
			desc = cc.WeatherDesc[0].Value
		}
		weatherData = fmt.Sprintf("%s  %s°C  %s%%", desc, cc.TempC, cc.Humidity)
		weatherTime = time.Now()
	}
	return weatherData
}

func (s *Service) BuildNewsHTML() string {
	items, err := s.news.Fetch(context.Background(), "top-news", 15)
	if err != nil {
		return ""
	}
	htmlItems := ""
	for _, item := range items {
		htmlItems += fmt.Sprintf(`<li><a target="_blank" href="%s">%s</a></li>`, html.EscapeString(item.Link), html.EscapeString(item.Title))
	}
	return htmlItems
}

func (s *Service) BuildCalendar(now time.Time) string {
	year, month, today := now.Year(), now.Month(), now.Day()
	first := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	last := first.AddDate(0, 1, -1)
	wd := int(first.Weekday())
	if wd == 0 {
		wd = 7
	}
	wd--
	var b strings.Builder
	monthLabel := fmt.Sprintf("%d年%d月", year, month)
	b.WriteString(`<div class="cal-wrap"><div class="cal-month">` + monthLabel + `</div><table class="cal"><tr>`)
	days := []string{"一", "二", "三", "四", "五", "六", "日"}
	for i, d := range days {
		cls := ""
		if i >= 5 {
			cls = ` class="wd"`
		}
		b.WriteString(fmt.Sprintf(`<th%s>%s</th>`, cls, d))
	}
	b.WriteString(`</tr><tr>`)
	for i := 0; i < wd; i++ {
		b.WriteString(`<td></td>`)
	}
	for d := 1; d <= last.Day(); d++ {
		cls := ""
		if d == today {
			cls = ` class="today"`
		} else if (wd+d-1)%7 >= 5 {
			cls = ` class="weekend"`
		}
		b.WriteString(fmt.Sprintf(`<td%s>%d</td>`, cls, d))
		if (wd+d)%7 == 0 && d != last.Day() {
			b.WriteString(`</tr><tr>`)
		}
	}
	b.WriteString(`</tr></table></div>`)
	return b.String()
}

// ---------- charts ----------

var (
	chartCache   map[string][]chartPoint
	chartCacheMu sync.Mutex
)

const yahooURL = "https://query1.finance.yahoo.com/v8/finance/chart/%s?range=1mo&interval=1d"
const tencentURL = "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?param=%s,day,,,30,qfq"

func fetchYahoo(symbol string) []chartPoint {
	url := fmt.Sprintf(yahooURL, symbol)
	c := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := c.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var d struct {
		Chart struct {
			Result []struct {
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Close []float64 `json:"close"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
		} `json:"chart"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil
	}
	if len(d.Chart.Result) == 0 {
		return nil
	}
	r := d.Chart.Result[0]
	var pts []chartPoint
	for i, ts := range r.Timestamp {
		if i >= len(r.Indicators.Quote[0].Close) {
			break
		}
		v := r.Indicators.Quote[0].Close[i]
		if v == 0 {
			continue
		}
		pts = append(pts, chartPoint{
			Label: time.Unix(ts, 0).Format("01/02"),
			Value: v,
		})
	}
	return pts
}

func fetchTencent(param string) []chartPoint {
	url := fmt.Sprintf(tencentURL, param)
	c := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := c.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var d struct {
		Data map[string]struct {
			Day [][]any `json:"day"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil
	}
	for _, v := range d.Data {
		var pts []chartPoint
		for _, row := range v.Day {
			if len(row) < 3 {
				continue
			}
			date, _ := row[0].(string)
			closeStr, _ := row[2].(string)
			if date == "" || closeStr == "" {
				continue
			}
			var closeVal float64
			fmt.Sscanf(closeStr, "%f", &closeVal)
			if closeVal == 0 {
				continue
			}
			t, err := time.Parse("2006-01-02", date)
			if err != nil {
				continue
			}
			pts = append(pts, chartPoint{
				Label: t.Format("01/02"),
				Value: closeVal,
			})
		}
		return pts
	}
	return nil
}

func (s *Service) MarketChart() string {
	chartMu.Lock()
	if time.Since(chartTime) > 30*time.Minute || chartGPSC == nil {
		chartGPSC = fetchYahoo("%5EGSPC")
		chartSH = fetchTencent("sh000001")
		chartGC = fetchYahoo("GC%3DF")
		chartTime = time.Now()
	}
	gspc, sh, gc := chartGPSC, chartSH, chartGC
	chartMu.Unlock()

	var all []chartSeries
	if len(gspc) > 0 {
		all = append(all, chartSeries{Name: "S&P 500", Color: "#000", Points: gspc})
	}
	if len(sh) > 0 {
		all = append(all, chartSeries{Name: "上证指数", Color: "#666", Points: sh})
	}
	if len(gc) > 0 {
		all = append(all, chartSeries{Name: "黄金", Color: "#999", Points: gc})
	}
	if len(all) < 2 {
		return ""
	}
	return svgMultiLineChart(all, 600, 220, "")
}
