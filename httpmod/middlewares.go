package httpmod

// Response
import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 通用插件相关内容

// / Default Plugin part
var headerXRequestID string = "X-Request-ID"

func PluginRequestId() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Get id from request
		rid := c.GetHeader(headerXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Request.Header.Add(headerXRequestID, rid)
		}
		// Set the id to ensure that the request-id is in the response
		c.Header(headerXRequestID, rid)
		c.Next()
	}
}

func GetRequestId(c *gin.Context) string {
	return c.Writer.Header().Get(headerXRequestID)
}

//

func PluginCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method //请求方法
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,Authorization-Token,AuthorizationToken")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}

// / --
// / default save the parmas into sqlite
type requestRecord struct {
	RequestId string `gorm:"column:request;size:64"`
	Method    string `gorm:"column:method;size:64"`
	Url       string `gorm:"column:url;size:1024"`
	Body      []byte `gorm:"column:body"`
	//TimeStamp uint64
}

func (r requestRecord) TableName() string {
	return "request"
}

// func getRequestRecordDb() (string, error) {
// 	if d, e := getResourceTypePath(ResourceTypeSqlite); e != nil {
// 		return d, e
// 	} else {
// 		return path.Join(d, "http_request.db"), nil
// 	}
// }

// var _pluginRequest = "_pluginRequest"

// func PluginRequestSnapShot() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		db, e := getRequestRecordDb()
// 		if e == nil {
// 			if !SqliteExists(_pluginRequest) {
// 				c := &conf.SQLiteConf{SQLitePath: db}
// 				initSqlite(_pluginRequest, c)
// 				SQLITE(_pluginRequest).AutoMigrate(requestRecord{})
// 			}

// 			go func() {
// 				body, e := io.ReadAll(c.Request.Body)
// 				if e == nil {
// 					SQLITE(_pluginRequest).Create(&requestRecord{
// 						RequestId: GetRequestId(c),
// 						Method:    c.Request.Method,
// 						Url:       fmt.Sprintf("%s?%s", c.Request.URL.Path, c.Request.URL.RawQuery),
// 						Body:      body,
// 					})
// 				}
// 				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
// 			}()
// 		}

// 		c.Next()
// 	}
// }
