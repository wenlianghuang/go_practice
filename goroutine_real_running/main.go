package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type FetchResult struct {
	URL      string
	Status   int
	Duration time.Duration
	Err      error
}

// fetchWithRetry 會套用 QPS 速率限制、具備退避重試與 context 取消。
// 重試條件：網路錯誤、HTTP 429、5xx；4xx（非 429）視為終止錯誤不重試。
func fetchWithRetry(ctx context.Context, client *http.Client, rate <-chan time.Time, url string, attempts int, baseBackoff time.Duration) FetchResult {
	start := time.Now()
	for i := 1; i <= attempts; i++ {
		// 速率限制
		select {
		case <-ctx.Done():
			return FetchResult{URL: url, Err: ctx.Err()}
		case <-rate:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return FetchResult{URL: url, Err: err}
		}

		resp, err := client.Do(req)
		if err != nil {
			// 網路錯誤：考慮重試
			if i == attempts {
				return FetchResult{URL: url, Duration: time.Since(start), Err: fmt.Errorf("request failed: %w", err)}
			}
			waitWithJitter(ctx, baseBackoff, i)
			continue
		}

		// 確保連線可重用：把 body 讀完丟掉再關閉
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		// 狀態碼處理
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			// 429 或 5xx 可重試
			if i == attempts {
				return FetchResult{
					URL:      url,
					Status:   resp.StatusCode,
					Duration: time.Since(start),
					Err:      fmt.Errorf("status %d after retries", resp.StatusCode),
				}
			}
			waitWithJitter(ctx, baseBackoff, i)
			continue
		}
		if resp.StatusCode >= 400 {
			// 其他 4xx 視為終止錯誤
			return FetchResult{
				URL:      url,
				Status:   resp.StatusCode,
				Duration: time.Since(start),
				Err:      fmt.Errorf("client error %d", resp.StatusCode),
			}
		}

		// 成功
		return FetchResult{
			URL:      url,
			Status:   resp.StatusCode,
			Duration: time.Since(start),
			Err:      nil,
		}
	}

	// 理論不會到這裡
	return FetchResult{URL: url, Duration: time.Since(start), Err: fmt.Errorf("exhausted retries")}
}

func waitWithJitter(ctx context.Context, base time.Duration, attempt int) {
	// 指數退避 + 抖動：base * 2^(attempt-1) + jitter(20%)
	backoff := base * time.Duration(1<<(attempt-1))
	jitter := time.Duration(rand.Int63n(int64(backoff / 5)))
	timer := time.NewTimer(backoff + jitter)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

func main() {
	// 允許 Ctrl+C 優雅關閉
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 測試 URL（含成功與故意失敗）
	urls := []string{
		"https://example.com",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/503", // 會重試
		"https://httpbin.org/delay/2",
		"https://golang.org",
	}

	// 參數可依你的場景調整
	concurrency := 4 // 有界併發上限
	qps := 4         // 每秒最多請求數（簡單速率限制）
	attempts := 3    // 重試次數
	baseBackoff := 200 * time.Millisecond
	clientTimeout := 4 * time.Second

	// HTTP client 與連線池設定（實務依需求調整）
	client := &http.Client{
		Timeout: clientTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        64,
			MaxConnsPerHost:     16,
			IdleConnTimeout:     60 * time.Second,
			DisableCompression:  false,
			ForceAttemptHTTP2:   true,
			MaxIdleConnsPerHost: 16,
		},
	}

	// 速率限制（簡單 ticker 實作）
	rate := time.NewTicker(time.Second / time.Duration(qps))
	defer rate.Stop()

	// 有界併發 + 錯誤傳播
	g, ctx := errgroup.WithContext(rootCtx)
	g.SetLimit(concurrency)

	results := make(chan FetchResult, len(urls)) // buffer 避免 aggregator 變慢時阻塞
	var okCount, failCount int32

	for _, u := range urls {
		u := u
		g.Go(func() error {
			res := fetchWithRetry(ctx, client, rate.C, u, attempts, baseBackoff)
			// 將結果送出（若已取消則丟棄）
			select {
			case results <- res:
			case <-ctx.Done():
			}
			if res.Err != nil {
				atomic.AddInt32(&failCount, 1)
				return res.Err // 失敗即刻取消其他 goroutine（可改為不返回以採用「錯誤容忍」策略）
			}
			atomic.AddInt32(&okCount, 1)
			return nil
		})
	}

	// 結果彙整（fan-in）
	done := make(chan struct{})
	go func() {
		defer close(done)
		for r := range results {
			if r.Err != nil {
				fmt.Printf("ERR  | %-35s | %v (status=%d)\n", r.URL, r.Err, r.Status)
			} else {
				fmt.Printf("OK   | %-35s | %d | %v\n", r.URL, r.Status, r.Duration)
			}
		}
	}()

	// 等待並關閉結果通道
	if err := g.Wait(); err != nil {
		fmt.Println("canceled by error:", err)
	}
	close(results)
	<-done // 最後send to main goroutine

	fmt.Printf("summary: ok=%d fail=%d\n", okCount, failCount)
}
