package rest

import (
	"net/http"

	"github.com/ethanburkett/goadmin/app/metrics"
	"github.com/gin-gonic/gin"
)

func RegisterMetricsRoutes(r *gin.Engine, api *Api) {
	metricsGroup := r.Group("/metrics")
	{
		// Prometheus metrics endpoint (no auth for scraping)
		metricsGroup.GET("", getPrometheusMetrics(api))

		// JSON metrics endpoint (requires auth)
		metricsGroup.GET("/json", AuthMiddleware(), getJSONMetrics(api))
	}
}

func getPrometheusMetrics(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		m, err := metrics.GetMetrics()
		if err != nil {
			c.String(http.StatusInternalServerError, "# Error fetching metrics\n")
			return
		}

		c.Header("Content-Type", "text/plain; version=0.0.4")
		c.String(http.StatusOK, m.PrometheusFormat())
	}
}

func getJSONMetrics(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		m, err := metrics.GetMetrics()
		if err != nil {
			c.Set("error", "Failed to fetch metrics")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", m)
		c.Status(http.StatusOK)
	}
}
