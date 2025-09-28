package plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
)

// ExecutionContext 插件执行上下文
type ExecutionContext struct {
	PluginID  int64  `json:"plugin_id"` // 插件ID
	EnvID     int64  `json:"env_id"`    // 环境变量ID
	EnvValue  string `json:"env_value"` // 环境变量值
	Config    []byte `json:"config"`    // 插件配置
	Timestamp int64  `json:"timestamp"` // 时间戳
}

// ExecutionResult 插件执行结果
type ExecutionResult struct {
	Success       bool   `json:"success"`        // 执行是否成功
	OutputData    []byte `json:"output_data"`    // 输出数据
	ErrorMessage  string `json:"error_message"`  // 错误信息
	ExecutionTime int    `json:"execution_time"` // 执行耗时(毫秒)
	StackTrace    string `json:"stack_trace"`    // 错误堆栈
}

// Engine 插件执行引擎
type Engine struct {
	timeout time.Duration // 默认超时时间
}

// NewEngine 创建插件执行引擎
func NewEngine(defaultTimeout time.Duration) *Engine {
	return &Engine{
		timeout: defaultTimeout,
	}
}

// Execute 执行插件脚本
func (e *Engine) Execute(ctx context.Context, script string, execCtx *ExecutionContext, timeout time.Duration) *ExecutionResult {
	startTime := time.Now()

	// 使用传入的超时时间，如果为0则使用默认超时时间
	if timeout == 0 {
		timeout = e.timeout
	}

	// 创建带超时的上下文
	execContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 创建执行结果
	result := &ExecutionResult{
		Success: false,
	}

	// 在goroutine中执行脚本
	done := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				result.ErrorMessage = fmt.Sprintf("插件执行发生panic: %v", r)
				result.StackTrace = fmt.Sprintf("%+v", r)
			}
			result.ExecutionTime = int(time.Since(startTime).Milliseconds())
			close(done)
		}()

		// 执行脚本
		e.executeScript(script, execCtx, result)
	}()

	// 等待执行完成或超时
	select {
	case <-done:
		// 执行完成
		return result
	case <-execContext.Done():
		// 超时
		result.ErrorMessage = "插件执行超时"
		result.ExecutionTime = int(timeout.Milliseconds())
		return result
	}
}

// executeScript 执行JavaScript脚本
func (e *Engine) executeScript(script string, execCtx *ExecutionContext, result *ExecutionResult) {
	// 创建goja运行时
	vm := goja.New()

	// 设置全局变量和函数
	e.setupGlobalObjects(vm, execCtx)

	// 包装用户脚本，确保有正确的入口函数
	wrappedScript := e.wrapScript(script, execCtx.EnvValue)

	defer func() {
		if r := recover(); r != nil {
			if jsErr, ok := r.(*goja.Exception); ok {
				result.ErrorMessage = jsErr.Error()
				result.StackTrace = jsErr.String()
			} else {
				result.ErrorMessage = fmt.Sprintf("脚本执行错误: %v", r)
				result.StackTrace = fmt.Sprintf("%+v", r)
			}
		}
	}()

	// 编译并执行脚本
	resultValue, err := vm.RunString(wrappedScript)
	if err != nil {
		result.ErrorMessage = err.Error()
		var jsErr *goja.Exception
		if errors.As(err, &jsErr) {
			result.StackTrace = jsErr.String()
		}
		return
	}

	// 处理返回值
	if resultValue != nil && !goja.IsUndefined(resultValue) {
		exported := resultValue.Export()
		if outputBytes, err := config.JSON.Marshal(exported); err == nil {
			result.OutputData = outputBytes
		} else {
			result.ErrorMessage = fmt.Sprintf("无法序列化输出数据: %v", err)
			return
		}
	}

	result.Success = true
}

// setupGlobalObjects 设置全局对象和函数
func (e *Engine) setupGlobalObjects(vm *goja.Runtime, execCtx *ExecutionContext) {
	// 设置执行上下文
	if err := vm.Set("context", map[string]interface{}{
		"pluginId":  execCtx.PluginID,
		"envId":     execCtx.EnvID,
		"timestamp": execCtx.Timestamp,
	}); err != nil {
		config.Log.Warn(err.Error()) // 仅做错误记录
	}

	// 设置配置数据
	if len(execCtx.Config) > 0 {
		var configData interface{}
		if err := config.JSON.Unmarshal(execCtx.Config, &configData); err == nil {
			if errSet := vm.Set("config", configData); errSet != nil {
				config.Log.Warn(errSet.Error()) // 仅做错误记录
			}
		}
	}

	// 设置工具函数
	if errSet := vm.Set("console", map[string]interface{}{
		"log": func(args ...interface{}) {
			// 这里可以集成到日志系统
			logArgs := make([]interface{}, 0, len(args)+1)
			logArgs = append(logArgs, "[Plugin Log]")
			logArgs = append(logArgs, args...)
			fmt.Println(logArgs...)
		},
		"error": func(args ...interface{}) {
			logArgs := make([]interface{}, 0, len(args)+1)
			logArgs = append(logArgs, "[Plugin Error]")
			logArgs = append(logArgs, args...)
			fmt.Println(logArgs...)
		},
	}); errSet != nil {
		config.Log.Warn(errSet.Error()) // 仅做错误记录
	}

	// 设置网络请求函数
	if errSet := vm.Set("request", func(options map[string]interface{}) interface{} {
		return e.makeHTTPRequest(options)
	}); errSet != nil {
		config.Log.Warn(errSet.Error()) // 仅做错误记录
	}

	// 设置JSON工具
	if errSet := vm.Set("JSON", map[string]interface{}{
		"parse": func(str string) interface{} {
			var result interface{}
			if err := config.JSON.Unmarshal([]byte(str), &result); err != nil {
				panic(vm.ToValue(err.Error()))
			}
			return result
		},
		"stringify": func(obj interface{}) string {
			bytes, err := config.JSON.Marshal(obj)
			if err != nil {
				panic(vm.ToValue(err.Error()))
			}
			return string(bytes)
		},
	}); errSet != nil {
		config.Log.Warn(errSet.Error()) // 仅做错误记录
	}

	// 设置时间工具
	if errSet := vm.Set("Date", map[string]interface{}{
		"now": func() int64 {
			return time.Now().UnixMilli()
		},
	}); errSet != nil {
		config.Log.Warn(errSet.Error()) // 仅做错误记录
	}
}

// wrapScript 包装用户脚本，确保有正确的结构
func (e *Engine) wrapScript(userScript string, envValue string) string {
	return fmt.Sprintf(`
// 用户脚本
%s

// 确保有main函数
if (typeof main !== 'function') {
    throw new Error('脚本必须定义一个main函数作为入口点');
}

// 执行main函数并传入环境变量值
var __result = main("%s");
__result;
`, userScript, envValue)
}

// ValidateScript 验证脚本语法
func (e *Engine) ValidateScript(script string) error {
	vm := goja.New()
	wrappedScript := e.wrapScript(script, "test_env_value")

	_, err := vm.RunString(wrappedScript)
	return err
}

// TestScript 测试脚本执行
func (e *Engine) TestScript(script string, envValue string) *ExecutionResult {
	execCtx := &ExecutionContext{
		PluginID:  0,
		EnvID:     0,
		EnvValue:  envValue,
		Config:    []byte(`{}`),
		Timestamp: time.Now().Unix(),
	}

	return e.Execute(context.Background(), script, execCtx, 10*time.Second)
}

// makeHTTPRequest 执行HTTP请求
func (e *Engine) makeHTTPRequest(options map[string]interface{}) interface{} {
	// 安全检查：只允许特定的域名和协议
	urlStr, ok := options["url"].(string)
	if !ok || urlStr == "" {
		return map[string]interface{}{
			"error": "URL is required",
		}
	}

	// 检查URL是否安全
	if !e.isURLSafe(urlStr) {
		return map[string]interface{}{
			"error": "URL not allowed",
		}
	}

	method, ok := options["method"].(string)
	if !ok {
		method = "GET"
	}
	method = strings.ToUpper(method)

	// 创建HTTP客户端，设置超时
	timeout := 10 * time.Second
	if timeoutVal, ok := options["timeout"].(int); ok && timeoutVal > 0 {
		timeout = time.Duration(timeoutVal) * time.Millisecond
		// 限制最大超时时间为30秒
		if timeout > 30*time.Second {
			timeout = 30 * time.Second
		}
	}

	client := &http.Client{
		Timeout: timeout,
	}

	// 准备请求体
	var body io.Reader
	if data, ok := options["data"]; ok && data != nil {
		if dataStr, ok := data.(string); ok {
			body = strings.NewReader(dataStr)
		} else {
			// 尝试JSON序列化
			if jsonData, err := config.JSON.Marshal(data); err == nil {
				body = strings.NewReader(string(jsonData))
			}
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	// 设置请求头
	if headers, ok := options["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if valueStr, ok := value.(string); ok {
				req.Header.Set(key, valueStr)
			}
		}
	}

	// 设置默认User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "QLToolsV2-Plugin/1.0")
	}

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}
	defer func(Body io.ReadCloser) {
		if errClose := Body.Close(); errClose != nil {
			config.Log.Error(errClose.Error())
		}
	}(resp.Body)

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	// 尝试解析JSON响应
	var jsonResp interface{}
	if err := config.JSON.Unmarshal(respBody, &jsonResp); err == nil {
		return jsonResp
	}

	// 如果不是JSON，返回字符串
	return map[string]interface{}{
		"status": resp.StatusCode,
		"body":   string(respBody),
		"headers": func() map[string]string {
			headers := make(map[string]string)
			for key, values := range resp.Header {
				if len(values) > 0 {
					headers[key] = values[0]
				}
			}
			return headers
		}(),
	}
}

// isURLSafe 检查URL是否安全
func (e *Engine) isURLSafe(urlStr string) bool {
	// 只允许HTTP和HTTPS协议
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	// 禁止访问内网地址
	if strings.Contains(urlStr, "localhost") ||
		strings.Contains(urlStr, "127.0.0.1") ||
		strings.Contains(urlStr, "0.0.0.0") ||
		strings.Contains(urlStr, "10.") ||
		strings.Contains(urlStr, "192.168.") ||
		strings.Contains(urlStr, "172.16.") ||
		strings.Contains(urlStr, "172.17.") ||
		strings.Contains(urlStr, "172.18.") ||
		strings.Contains(urlStr, "172.19.") ||
		strings.Contains(urlStr, "172.20.") ||
		strings.Contains(urlStr, "172.21.") ||
		strings.Contains(urlStr, "172.22.") ||
		strings.Contains(urlStr, "172.23.") ||
		strings.Contains(urlStr, "172.24.") ||
		strings.Contains(urlStr, "172.25.") ||
		strings.Contains(urlStr, "172.26.") ||
		strings.Contains(urlStr, "172.27.") ||
		strings.Contains(urlStr, "172.28.") ||
		strings.Contains(urlStr, "172.29.") ||
		strings.Contains(urlStr, "172.30.") ||
		strings.Contains(urlStr, "172.31.") {
		return false
	}

	return true
}
