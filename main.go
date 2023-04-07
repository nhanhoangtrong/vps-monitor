package main

import (
	"fmt"
	"net/http"

	human "github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/disk"
)

func printUsage() ([]string, error) {
	ret := []string{"FStype Total Used Free Percent Mountpoint"}
	formatter := "%-14s %7s %7s %7s %4s %s"
	parts, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	for _, p := range parts {
		device := p.Mountpoint
		s, err := disk.Usage(device)
		if err != nil {
			return nil, err
		}

		if s.Total == 0 {
			continue
		}

		percent := fmt.Sprintf("%2.f%%", s.UsedPercent)

		ret = append(ret, fmt.Sprintf(formatter,
			s.Fstype,
			human.Bytes(s.Total),
			human.Bytes(s.Used),
			human.Bytes(s.Free),
			percent,
			p.Mountpoint,
		))
	}
	return ret, nil
}

func main() {
	// Initialize gin
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Get disk info
	r.GET("/disk", func(c *gin.Context) {
		usages, err := printUsage()
		if err != nil {
			c.JSON(http.StatusFailedDependency, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.HTML(http.StatusOK, "disk.html", gin.H{
			"usages": usages,
		})
	})

	r.Run()
}
