package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"net/http"
	"time"
)

//定义接收数据的结构体
type Login struct {
	//binding:必选字段
	UserName string   `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	PassWord string   `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

//定义全局中间件
func MiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		t1 := time.Now()
		fmt.Println("中间件执行....")
		context.Set("request", "中间件")
		//执行
		context.Next()
		status := context.Writer.Status()
		fmt.Println("中间件执行完成status", status)
		t2 := time.Since(t1)
		fmt.Println("中间件执行时间", t2)
	}
}

func main() {
	r := gin.Default()

	//花括号内为这个中间件需要被那些请求执行
	r.Use(MiddleWare())
	{
		r.GET("/middleware", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message" : "middleware ok",
				"value" : context.Get("request"),
			})
		})

		r.GET("/middleware2", MiddleWare(), func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message" : "middleware ok",
				"value" : context.Get("request"),
			})
		})
	}

	//服务器向客户端写cookie
	r.GET("/cookie", func(context *gin.Context) {
		//客户端是否携带cookie
		cookie, err := context.Cookie("key")
		if err != nil {
			cookie = "NotSet"
			//name, value string,
			//maxAge int(s), 过期时间
			//path, 所在目录
			//domain string 域名
			//secure, 是否只能通过https访问
			//httpOnly bool,是否通过js获取自己的cookie
			context.SetCookie("key", "somecookie", 60, "/", "localhost", false, false)
		}

		fmt.Println("cookie 的值", cookie)
	})

	//表单处理
	r.POST("/form", func(c *gin.Context) {
		//表单参数设置默认值
		ctype := c.DefaultPostForm("type", "alert")
		//其他参数
		useName := c.PostForm("username")
		passWord := c.PostForm("password")
		//多选框
		hobbies := c.PostFormArray("hobby")
		c.String(http.StatusOK, fmt.Sprintf("ctype:%s, userName:%s, passWord:%s, hobbies:%v",
			                                         ctype, useName, passWord, hobbies))
	})

	//上传文件
	r.POST("/upload", func(context *gin.Context) {
		file, _ := context.FormFile("file")
		fmt.Println("filename", file.Filename)
		//上传至当前目录
		context.SaveUploadedFile(file, file.Filename)
		context.String(http.StatusOK, fmt.Sprintf("upload filename:%s ok", file.Filename))

	})

	//上传多个文件
	//限制表单上传大小8MB,gin默认设置为32MB
	r.MaxMultipartMemory = 8 << 20
	r.POST("/multipleupload", func(context *gin.Context) {
		form, _ := context.MultipartForm()

		//获取所有文件
		files := form.File["files"]
		for _, file := range files {
			//逐个存储
			context.SaveUploadedFile(file, file.Filename)
		}
		context.String(http.StatusOK, fmt.Sprintf("multi upload %d file ok", len(files)))

	})

	//路由组1:处理GET请求
	v1 := r.Group("/v1")
	{
		v1.GET("/login", login)
		v1.GET("/logout", logout)
	}

	//路由组2:处理POST请求
	v2 := r.Group("/v2")
	{
		v2.POST("/login", login)
		v2.POST("/logout", logout)
	}

	//JSON数据解析与绑定
	r.POST("/loginjson", func(context *gin.Context) {
		//body数据按照json格式解析到结构体
		var json Login
		if err := context.ShouldBindJSON(&json); err != nil {
			//gin.H 生成json的工具
			context.JSON(http.StatusBadRequest, gin.H{
				"code" : -1,
				"error": err.Error(),
			})
			return
		}

		if json.UserName != "lihui" && json.PassWord != "123456" {
			context.JSON(http.StatusBadRequest, gin.H{
				"code" : -2,
				"info": "username or password not ok",
			})

			return
		}

		context.JSON(http.StatusOK, gin.H{"code":200})


	})


	//表单数据绑定与解析
	r.POST("/loginform", func(context *gin.Context) {
		var form Login
		if err := context.Bind(&form); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error":err.Error(),
			})

			return
		}


		if form.UserName != "lihui" && form.PassWord != "123456" {
			context.JSON(http.StatusBadRequest, gin.H{
				"code" : -2,
				"info": "username or password not ok",
			})

			return
		}

		context.JSON(http.StatusOK, gin.H{"code":200})

	})

	//URI数据绑定与解析
	r.GET("/login/:user/:password", func(context *gin.Context) {
		var uri Login
		if err := context.ShouldBindUri(&uri); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error":err.Error(),
			})

			return
		}


		if uri.UserName != "lihui" && uri.PassWord != "123456" {
			context.JSON(http.StatusBadRequest, gin.H{
				"code" : -2,
				"info": "username or password not ok",
			})

			return
		}

		context.JSON(http.StatusOK, gin.H{"code":200})

	})

	//以JSON格式响应
	r.GET("/somejson", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message" : "ok",
		})
	})

	//以结构体响应
	r.GET("/somestruct", func(context *gin.Context) {
		s := Login{
			UserName:"lihui",
			PassWord:"123435",
		}

		context.JSON(200, s)

	})

	//以XML格式响应
	r.GET("/somexml", func(context *gin.Context) {
		context.XML(200, gin.H{
			"message" : "ok",
		})
	})

	//以YAML格式响应
	r.GET("/someyaml", func(context *gin.Context) {
		context.YAML(200, gin.H{
			"message" : "ok",
		})
	})

	//以protobuf格式响应
	r.GET("/someproto", func(context *gin.Context) {
		//构建protobuf格式数据
		label := "label"
		reps := []int64{1, 2}
		data := protoexample.Test{
			Label : &label,
			Reps:reps,

		}
		context.ProtoBuf(200, data)
	})

    //HTML渲染
    //加载模板文件
    r.LoadHTMLGlob("templates/*")
	//r.LoadHTMLFiles("templates/index.tmpl")
	r.GET("/index", func(context *gin.Context) {
		context.HTML(200, "index.tmpl", gin.H{
			"title" : "我的标题",
		})
	})


	//重定向
	r.GET("/redirect", func(context *gin.Context) {
		//支持内部和外部重定向
		context.Redirect(http.StatusMovedPermanently, "https://github.com")
	})

	//同步异步
	r.GET("/async", func(context *gin.Context) {
		copyContext := context.Copy()      //注意只能用上下文的副本
		go func(c *gin.Context) {
			time.Sleep(3 * time.Second)
			fmt.Println("异步执行...", c.Request.URL.Path)
		}(copyContext)

		context.JSON(200, gin.H{
			"message" : "ok",
		})
	})

	r.GET("/sync", func(context *gin.Context) {
		time.Sleep(3 * time.Second)
		fmt.Println("同步执行...", context.Request.URL.Path)
		context.JSON(200, gin.H{
			"message" : "ok",
		})
	})

	r.Run(":8888")
}


func login(ctx *gin.Context) {
	ctx.String(http.StatusOK, "login ok")
}

func logout(ctx *gin.Context) {
	ctx.String(http.StatusOK, "logout ok")
}