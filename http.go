package ezgo

// steal from https://github.com/shenghui0779/yiigo/blob/master/http.go

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"time"
)

type httpSetting struct {
	headers map[string]string
	cookies []*http.Cookie
	close   bool
}

// HTTPOption HTTP请求选项
type HTTPOption func(s *httpSetting)

// WithHTTPHeader 设置HTTP请求头
func WithHTTPHeader(key, value string) HTTPOption {
	return func(s *httpSetting) {
		s.headers[key] = value
	}
}

// WithHTTPCookies 设置HTTP请求Cookie
func WithHTTPCookies(cookies ...*http.Cookie) HTTPOption {
	return func(s *httpSetting) {
		s.cookies = cookies
	}
}

// WithHTTPClose 请求结束后关闭请求
func WithHTTPClose() HTTPOption {
	return func(s *httpSetting) {
		s.close = true
	}
}

// UploadForm HTTP文件上传表单
type UploadForm interface {
	// Write 将表单字段写入流
	Write(w *multipart.Writer) error
}

// FormFileFunc 将上传文件写入流
type FormFileFunc func(w io.Writer) error

type formfile struct {
	fieldname string
	filename  string
	filefunc  FormFileFunc
}

type uploadform struct {
	formfiles  []*formfile
	formfields map[string]string
}

func (f *uploadform) Write(w *multipart.Writer) error {
	if len(f.formfiles) == 0 {
		return errors.New("empty file field")
	}

	for _, v := range f.formfiles {
		part, err := w.CreateFormFile(v.fieldname, v.filename)

		if err != nil {
			return err
		}

		if err = v.filefunc(part); err != nil {
			return err
		}
	}

	for name, value := range f.formfields {
		if err := w.WriteField(name, value); err != nil {
			return err
		}
	}

	return nil
}

// UploadField 文件上传表单字段
type UploadField func(f *uploadform)

// WithFormFile 设置表单文件字段
func WithFormFile(fieldname, filename string, fn FormFileFunc) UploadField {
	return func(f *uploadform) {
		f.formfiles = append(f.formfiles, &formfile{
			fieldname: fieldname,
			filename:  filename,
			filefunc:  fn,
		})
	}
}

// WithFormField 设置表单普通字段
func WithFormField(fieldname, fieldvalue string) UploadField {
	return func(u *uploadform) {
		u.formfields[fieldname] = fieldvalue
	}
}

// NewUploadForm 生成一个文件上传表单
func NewUploadForm(fields ...UploadField) UploadForm {
	form := &uploadform{
		formfiles:  make([]*formfile, 0),
		formfields: make(map[string]string),
	}

	for _, f := range fields {
		f(form)
	}

	return form
}

// HTTPClient HTTP客户端
type HTTPClient interface {
	// Do 发送HTTP请求
	// 注意：应该使用Context设置请求超时时间
	Do(ctx context.Context, method, reqURL string, body []byte, options ...HTTPOption) (*http.Response, error)

	// Upload 上传文件
	// 注意：应该使用Context设置请求超时时间
	Upload(ctx context.Context, reqURL string, form UploadForm, options ...HTTPOption) (*http.Response, error)
}

type httpclient struct {
	client *http.Client
}

func (c *httpclient) Do(ctx context.Context, method, reqURL string, body []byte, options ...HTTPOption) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	setting := new(httpSetting)

	if len(options) != 0 {
		setting.headers = make(map[string]string)

		for _, f := range options {
			f(setting)
		}
	}

	// headers
	if len(setting.headers) != 0 {
		for k, v := range setting.headers {
			req.Header.Set(k, v)
		}
	}

	// cookies
	if len(setting.cookies) != 0 {
		for _, v := range setting.cookies {
			req.AddCookie(v)
		}
	}

	if setting.close {
		req.Close = true
	}

	resp, err := c.client.Do(req)

	if err != nil {
		// If the context has been canceled, the context's error is probably more useful.
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}

		return nil, err
	}

	return resp, nil
}

func (c *httpclient) Upload(ctx context.Context, reqURL string, form UploadForm, options ...HTTPOption) (*http.Response, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 20<<10)) // 20kb
	w := multipart.NewWriter(buf)

	if err := form.Write(w); err != nil {
		return nil, err
	}

	options = append(options, WithHTTPHeader("Content-Type", w.FormDataContentType()))

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	if err := w.Close(); err != nil {
		return nil, err
	}

	return c.Do(ctx, http.MethodPost, reqURL, buf.Bytes(), options...)
}

// NewDefaultHTTPClient 生成一个默认的HTTP客户端
func NewDefaultHTTPClient() HTTPClient {
	return &httpclient{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 60 * time.Second,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				MaxIdleConns:          0,
				MaxIdleConnsPerHost:   1000,
				MaxConnsPerHost:       1000,
				IdleConnTimeout:       60 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: time.Second,
			},
		},
	}
}

// NewHTTPClient 通过官方 `http.Client` 生成一个HTTP客户端
func NewHTTPClient(client *http.Client) HTTPClient {
	return &httpclient{
		client: client,
	}
}

var defaultHTTPClient = NewDefaultHTTPClient()

// HTTPGet 发送GET请求
func HTTPGet(ctx context.Context, reqURL string, options ...HTTPOption) (*http.Response, error) {
	return defaultHTTPClient.Do(ctx, http.MethodGet, reqURL, nil, options...)
}

// HTTPPost 发送POST请求
func HTTPPost(ctx context.Context, reqURL string, body []byte, options ...HTTPOption) (*http.Response, error) {
	return defaultHTTPClient.Do(ctx, http.MethodPost, reqURL, body, options...)
}

// HTTPPostForm 发送POST表单请求
func HTTPPostForm(ctx context.Context, reqURL string, data url.Values, options ...HTTPOption) (*http.Response, error) {
	options = append(options, WithHTTPHeader("Content-Type", "application/x-www-form-urlencoded"))

	return defaultHTTPClient.Do(ctx, http.MethodPost, reqURL, []byte(data.Encode()), options...)
}

// HTTPUpload 文件上传
func HTTPUpload(ctx context.Context, reqURL string, form UploadForm, options ...HTTPOption) (*http.Response, error) {
	return defaultHTTPClient.Upload(ctx, reqURL, form, options...)
}

// HTTPDo 发送HTTP请求
func HTTPDo(ctx context.Context, method, reqURL string, body []byte, options ...HTTPOption) (*http.Response, error) {
	return defaultHTTPClient.Do(ctx, method, reqURL, body, options...)
}
