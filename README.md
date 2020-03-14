### Golang
架构图：

https://www.processon.com/view/link/5b8f1968e4b06fc64ae4949f

https://www.processon.com/view/link/5d65f54de4b09965fad3d701

1. 项目目录介绍

    - `/build` **项目构建相关文件目录**
        - `/binary` **编译后的二进制包目录**
        - `Dockerfile` **Docker 镜像构建文件**
        - `Makefile` **自定义命令文件**
    - `/cmd` **项目服务目录**
        - `artisan` **命令中心**
        - `auth` **鉴权服务项目根目录**
        - `rxsc` **学生端 http 网关服务根目录**
        - `...`
    - `/config` **配置文件目录**
        - `config.yaml` **本地读取的配置文件**
    - `/internal` **服务间复用的包**
        - `api` **调用服务方式封装**
        - `cache` **缓存模块封装**
        - `config` **读写配置模块封装**
        - `core` **无状态组件封装**
        - `db` **数据库封装**
            - `mysql` **GROM封装**
            - `redis` **Redis客户端封装**
            - `es` **Elasticsearch服务封装**
        - `http` **Gin框架及中间件封装**
        - `jwt` **JWT Token加解密封装**
        - `log` **日志模块封装**
        - `wrong` **错误处理模块封装**
        - `util` **基础函数封装**

    - `/pkg` **依赖的外部包或服务**
    - `go.mod && go.sum` **依赖管理配置文件**
    - `README.md` **说明文件**
2. 第三方包列表 - 直接引用

    - 日志相关
        - `github.com/sirupsen/logrus` **支持Json记录日志包**
        - `github.com/lestrrat/go-file-rotatelogs` **logrus文件驱动**
    - `Http`相关
        - `github.com/gin-gonic/gin` **Gin框架**
        - `github.com/gin-contrib/cors` **Gin官方Cors跨域配置插件**
    - 数据库相关
        - `github.com/jinzhu/gorm` **SQL数据库ORM**
        - `github.com/go-redis/redis/v7` **Redis数据库客户端**
        - `github.com/olivere/elastic` **Elasticsearch客户端**
    - RPC相关
        - `github.com/golang/protobuf` **Proto协议解析包**
        - `google.golang.org/grpc` **Grpc框架**
    - 其他
        - `github.com/dgrijalva/jwt-go` **JWT Token加密解析包**
        - `github.com/spf13/cobra` **命令行生成工具**
        - `github.com/spf13/viper` **配置文件读写包**
        - `github.com/pkg/errors` **错误包**
3. 服务目录介绍

    - 网关服务
        - `/app` **项目内容**
            - `cmd` **命令定义**
            - `handler` **控制器**
        - `/route` **路由**
            - `route.go` **路由文件**
        - `main.go` **项目入口文件**
    - 业务服务
        - `/handler` **控制器 对外暴露服务**
        - `/model` **数据模型层**
        - `/service` **业务逻辑层**
        - `/proto` **Proto协议定义**
        - `/server` **Grpc客户端入口**
        - `main.go` **入口文件**
4. 启动服务

    - 项目入口 `/build/Makefile`
    - 本地调试启动命令：`make run`
    - 相见`Makefile`文件内的定义与配置

5. 部署项目

    - 本机编译： `make build` 可在 `build/binay`目录下生成对应二进制文件（`todo`: 交叉编译配置）

    - 使用`Docker`: `make deploy`
