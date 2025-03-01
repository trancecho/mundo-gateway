package pkg

import (
	"encoding/json"
	"fmt"
)

// ConfigRegistrySDK 表示配置注册中心的 SDK
type ConfigRegistrySDK struct {
	baseURL string
	client  *resty.Client
}

// NewConfigRegistrySDK 创建一个新的 SDK 实例
func NewConfigRegistrySDK(baseURL string) *ConfigRegistrySDK {
	client := resty.New()
	return &ConfigRegistrySDK{
		baseURL: baseURL,
		client:  client,
	}
}

// Ping 测试与配置注册中心的连接
func (s *ConfigRegistrySDK) Ping() (string, error) {
	url := fmt.Sprintf("%s/ping", s.baseURL)
	resp, err := s.client.R().Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("ping failed with status code: %d", resp.StatusCode())
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	message, ok := result["message"].(string)
	if !ok {
		return "", fmt.Errorf("failed to parse ping response")
	}

	return message, nil
}

// CreateNamespace 创建命名空间
func (s *ConfigRegistrySDK) CreateNamespace(namespaceData interface{}) error {
	url := fmt.Sprintf("%s/namespace/create", s.baseURL)
	return s.postRequest(url, namespaceData)
}

// GetAllNamespaces 获取所有命名空间
func (s *ConfigRegistrySDK) GetAllNamespaces() ([]interface{}, error) {
	url := fmt.Sprintf("%s/namespace/getAll", s.baseURL)
	return s.getRequest(url)
}

// DeleteNamespace 删除命名空间
func (s *ConfigRegistrySDK) DeleteNamespace(namespaceName string) error {
	url := fmt.Sprintf("%s/namespace/delete?name=%s", s.baseURL, namespaceName)
	return s.getRequestWithErrorCheck(url)
}

// CreateKV 创建键值对
func (s *ConfigRegistrySDK) CreateKV(kvData interface{}) error {
	url := fmt.Sprintf("%s/kv/create", s.baseURL)
	return s.postRequest(url, kvData)
}

// GetKV 获取键值对
func (s *ConfigRegistrySDK) GetKV(key string) (interface{}, error) {
	url := fmt.Sprintf("%s/kv/get?key=%s", s.baseURL, key)
	results, err := s.getRequest(url)
	if err != nil {
		return nil, err
	}
	if len(results) > 0 {
		return results[0], nil
	}
	return nil, nil
}

// DeleteKV 删除键值对
func (s *ConfigRegistrySDK) DeleteKV(key string) error {
	url := fmt.Sprintf("%s/kv/delete?key=%s", s.baseURL, key)
	return s.getRequestWithErrorCheck(url)
}

// RegisterService 注册服务
func (s *ConfigRegistrySDK) RegisterService(serviceData interface{}) error {
	url := fmt.Sprintf("%s/service/register", s.baseURL)
	return s.postRequest(url, serviceData)
}

// UnRegisterService 注销服务
func (s *ConfigRegistrySDK) UnRegisterService(serviceData interface{}) error {
	url := fmt.Sprintf("%s/service/unregister", s.baseURL)
	return s.postRequest(url, serviceData)
}

// GetServices 获取服务
func (s *ConfigRegistrySDK) GetServices(serviceName string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/service/get?name=%s", s.baseURL, serviceName)
	return s.getRequest(url)
}

// GetAllServices 获取所有服务
func (s *ConfigRegistrySDK) GetAllServices() ([]interface{}, error) {
	url := fmt.Sprintf("%s/service/getAll", s.baseURL)
	return s.getRequest(url)
}

// postRequest 发送 POST 请求
func (s *ConfigRegistrySDK) postRequest(url string, data interface{}) error {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("post request failed with status code: %d", resp.StatusCode())
	}

	return nil
}

// getRequest 发送 GET 请求并返回结果
func (s *ConfigRegistrySDK) getRequest(url string) ([]interface{}, error) {
	resp, err := s.client.R().Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get request failed with status code: %d", resp.StatusCode())
	}

	var results []interface{}
	err = json.Unmarshal(resp.Body(), &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// getRequestWithErrorCheck 发送 GET 请求并检查错误
func (s *ConfigRegistrySDK) getRequestWithErrorCheck(url string) error {
	resp, err := s.client.R().Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("get request failed with status code: %d", resp.StatusCode())
	}

	return nil
}
