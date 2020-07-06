package api

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/hirakiuc/site-meta-go/sitemeta"
)

func respondErr(c *gin.Context, status int, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"code":    status,
		"message": err.Error(),
	})
}

func respondMeta(c *gin.Context, url string, meta *sitemeta.SiteMeta) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"url":    url,
		"meta":   meta.Attrs,
	})
}

func MetaHandler(c *gin.Context) {
	url, err := url.Parse(c.Query("url"))
	if err != nil {
		respondErr(c, http.StatusBadRequest, err)

		return
	}

	meta, err := sitemeta.Parse(c, url.String())
	if err != nil {
		respondErr(c, http.StatusInternalServerError, err)

		return
	}

	respondMeta(c, url.String(), meta)
}
