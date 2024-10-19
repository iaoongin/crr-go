package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type TemplateProcessor struct{}

func NewTemplateProcessor() *TemplateProcessor {
	return &TemplateProcessor{}
}

func processTemplate(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	remoteConfig, err := fetchRemoteYAML(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch remote YAML: %v", err)})
		return
	}

	localConfig, err := readLocalYAML(templateFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read local YAML: %v", err)})
		return
	}

	mergedConfig := mergeConfigs(localConfig, remoteConfig)

	result, err := yaml.Marshal(mergedConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to marshal merged YAML: %v", err)})
		return
	}

	c.String(http.StatusOK, string(result))
}

func fetchRemoteYAML(url string) (ClashConfig, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Clash")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var config ClashConfig
	err = yaml.Unmarshal(body, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func readLocalYAML(filename string) (ClashConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config ClashConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func mergeConfigs(local, remote ClashConfig) ClashConfig {
	merged := make(ClashConfig)

	// 复制本地配置的所有字段
	for k, v := range local {
		merged[k] = v
	}

	// 只替换 proxies 和 proxy-groups
	if proxies, ok := remote["proxies"]; ok {
		merged["proxies"] = proxies
	}
	if proxyGroups, ok := remote["proxy-groups"]; ok {
		merged["proxy-groups"] = proxyGroups
	}

	return merged
}
