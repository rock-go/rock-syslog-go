# rock-syslog-go
syslog 模块

#  syslog.*
- 函数: syslog.server( table )
- 参数: 配置参数服务的必要参数 
```lua
    local kafka = kafka.producer{}
    local file = rock.file{}
    local syslog = syslog.server{
        protocol  = "udp", -- udp , tcp , udp/tcp
        listen    = "0.0.0.0:514",
        
        -- RFC3178,RFC5424,RFC6587,Auto
        format    = syslog.Auto ,
        
        -- json , raw  数据保存格式
        format    = "raw", 
        
        output = {kafka , file } -- lua.Writer 接口的方法 
    }

    proc.start(kafka , file , syslog)-- 启动
```