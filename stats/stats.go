package stats

import (
	"time"

	"github.com/gin-gonic/gin"
)

const ginStats = "_ginStats"

//Stats is a container for information
type Stats struct {
	TotalRequests    int64         `json:"toalRequests"`
	TotalRequestTime time.Duration `json:"totalRequestTime"`
	StartTime        time.Time     `json:"startTime"`
	Responses        map[int]int   `json:"responses"`
	AverageTime      time.Duration `json:"averageTimeMs"`
}

func GetStats(c *gin.Context) *Stats {
	v, exists := c.Get(ginStats)
	if !exists {
		return nil
	}
	return v.(*Stats)
}

func InjectStats() gin.HandlerFunc {
	var _stats Stats
	_stats.StartTime = time.Now()
	_stats.Responses = make(map[int]int)

	return func(c *gin.Context) {
		start := time.Now()
		_stats.TotalRequests++
		c.Set(ginStats, &_stats)
		c.Next()
		duration := time.Now().Sub(start)

		status := c.Writer.Status()

		if _, ok := _stats.Responses[status]; ok {
			_stats.Responses[status]++
		} else {
			_stats.Responses[status] = 1
		}

		_stats.TotalRequestTime = _stats.TotalRequestTime + duration
		_stats.AverageTime = time.Duration(_stats.TotalRequests / _stats.AverageTime.Nanoseconds())
	}
}
