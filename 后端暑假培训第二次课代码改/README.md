## 接口：

增加/add接口来增加商品,传入格式(json):

`{`

  `"name":string,`

  `"price":int,`

  `"num":int`

`}`

补充:
```go
    //数据库加载到Map 
    func set_inMap() {
    
    	//假设每天凌晨12点加载
    	t := time.Now()
    	//获取明天的时间
    	next := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
    	//设置timer
    	timer := time.NewTimer(next.Sub(t))
    	
    	<-timer.C
    	//添加
    	initMap()
    
    }
```

对这个函数开启协程，则会在每晚上12点进行刷新
