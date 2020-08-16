package teacache

import (
	"bufio"
	"bytes"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"io"
	"net/http"
	"strings"
)

// 请求之前处理
func ProcessBeforeRequest(req *teaproxy.Request, writer *teaproxy.ResponseWriter) bool {
	req.SetVarMapping("cache.status", "BYPASS")
	req.SetVarMapping("cache.policy.name", "")
	req.SetVarMapping("cache.policy.type", "")

	cacheConfig := req.CachePolicy()
	if cacheConfig == nil || !cacheConfig.On {
		return true
	}

	// 匹配条件
	if len(cacheConfig.Cond) > 0 {
		for _, cond := range cacheConfig.Cond {
			if !cond.Match(req.Format) {
				return true
			}
		}
	}

	// 支持请求中使用Pragma或Cache-Control
	if cacheConfig.EnableRequestCachePragma {
		if req.Raw().Header.Get("Cache-Control") == "no-cache" || req.Raw().Header.Get("Pragma") == "no-cache" {
			return true
		}
	}

	req.SetVarMapping("cache.policy.name", cacheConfig.Name)
	req.SetVarMapping("cache.policy.type", cacheConfig.Type)

	cachePolicyMapLocker.RLock()
	cache, found := cachePolicyMap[cacheConfig.Filename]
	cachePolicyMapLocker.RUnlock()
	if !found {
		cachePolicyMapLocker.Lock()

		// find again
		cache, found = cachePolicyMap[cacheConfig.Filename]
		if found {
			cachePolicyMapLocker.Unlock()
		} else {
			cacheConfig = shared.NewCachePolicyFromFile(cacheConfig.Filename)
			if cacheConfig == nil {
				cachePolicyMapLocker.Unlock()
				return true
			}
			cache = NewManagerFromConfig(cacheConfig)
			if cache == nil {
				cachePolicyMapLocker.Unlock()
				return true
			}

			logs.Println("[cache]create cache policy instance:", cacheConfig.Name+"("+cacheConfig.Type+")")
			cachePolicyMap[cacheConfig.Filename] = cache
			cachePolicyMapLocker.Unlock()
		}
	}

	// key
	if len(cacheConfig.Key) == 0 {
		return true
	}
	key := req.Format(cacheConfig.Key)

	// 是否为清除缓存
	rawReq := req.Raw()
	teaKey := rawReq.Header.Get("Tea-Key")
	if rawReq.Header.Get("Tea-Cache-Purge") == "1" {
		req.SetVarMapping("cache.status", "PURGE")

		if len(teaKey) == 0 {
			writer.Write([]byte("ERROR:'Tea-Key' should be set in header"))
			return false
		}

		if configs.SharedAdminConfig().FindUserWithKey(teaKey) == nil {
			writer.Write([]byte("ERROR:Tea-Key:'" + teaKey + "' is incorrect"))
			return false
		}

		err := cache.Delete(key)
		if err != nil {
			writer.Write([]byte("ERROR:" + err.Error()))
		} else {
			writer.Write([]byte("ok"))
		}
		return false
	}

	// 读取缓存
	data, err := cache.Read(key)
	req.SetVarMapping("cache.status", "MISS")
	if err != nil {
		if err != ErrNotFound {
			logs.Error(err)
		} else {
			req.SetCacheEnabled()
			writer.SetBodyCopying(true)
		}
		return true
	}

	if len(data) <= 8 {
		return true
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data[8:])), nil)
	if err != nil {
		logs.Error(err)
		return true
	}

	defer resp.Body.Close()

	for k, vs := range resp.Header {
		if k == "Connection" {
			continue
		}
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}

	req.SetAttr("cache.cached", "1")
	req.SetAttr("cache.policy.name", cacheConfig.Name)
	req.SetAttr("cache.policy.type", cacheConfig.Type)

	// 添加变量
	req.SetVarMapping("cache.status", "HIT")

	// 自定义Response
	req.WriteResponseHeaders(writer, resp.StatusCode)

	writer.WriteHeader(resp.StatusCode)

	buf := make([]byte, 1024)
	_, err = io.CopyBuffer(writer, resp.Body, buf)
	if err != nil {
		logs.Error(err)
	}

	return false
}

// 请求之后处理
func ProcessAfterRequest(req *teaproxy.Request, writer *teaproxy.ResponseWriter) bool {
	if !req.IsCacheEnabled() {
		return true
	}

	cacheConfig := req.CachePolicy()
	if cacheConfig == nil {
		return true
	}

	//check status
	if writer.StatusCode() == http.StatusNotModified { // 如果没有修改就不会有body，会有陷阱，所以这里不加入缓存
		return true
	}
	if len(cacheConfig.Status) == 0 {
		cacheConfig.Status = []int{http.StatusOK}
	}
	if !lists.ContainsInt(cacheConfig.Status, writer.StatusCode()) {
		return true
	}

	// check length
	contentLength := types.Int(writer.Header().Get("Content-Length"))
	if contentLength != len(writer.Body()) && writer.Header().Get("Content-Encoding") != "gzip" {
		return true
	}

	// validate cache control
	if len(cacheConfig.SkipResponseCacheControlValues) > 0 {
		cacheControl := writer.Header().Get("Cache-Control")
		if len(cacheControl) > 0 {
			values := strings.Split(cacheControl, ",")
			for _, value := range values {
				if cacheConfig.ContainsCacheControl(strings.TrimSpace(value)) {
					return true
				}
			}
		}
	}

	// validate set cookie
	if cacheConfig.SkipResponseSetCookie && len(writer.Header().Get("Set-Cookie")) > 0 {
		return true
	}

	cachePolicyMapLocker.RLock()
	cache, found := cachePolicyMap[cacheConfig.Filename]
	cachePolicyMapLocker.RUnlock()
	if !found {
		return true
	}

	key := req.Format(cacheConfig.Key)
	headerData := writer.HeaderData()
	item := &Item{
		Header: headerData,
		Body:   writer.Body(),
	}
	if len(headerData) == 0 {
		return true
	}
	data := item.Encode()
	if cacheConfig.MaxDataSize() > 0 && float64(len(data)) > cacheConfig.MaxDataSize() {
		return true
	}
	err := cache.Write(key, data)
	if err != nil {
		logs.Error(err)
	}
	return true
}
