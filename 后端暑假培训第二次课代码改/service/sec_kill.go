package service

import (
	"log"
	"summerCourse/model"
	"sync"
	"time"
)

type User struct {
	UserId  string
	GoodsId uint
}

var OrderChan = make(chan User, 1024)

var ItemMap = make(map[uint]*Item)

type Item struct {
	ID        uint   // 商品id
	Name      string // 名字
	Total     int    // 商品总量
	Left      int    // 商品剩余数量
	IsSoldOut bool   // 是否售罄
	leftCh    chan int
	sellCh    chan int
	done      chan struct{}
	Lock      sync.Mutex
}

// TODO 写一个定时任务，每天定时从数据库加载数据到Map

//每天获取商品
func initMap() {
	goods, err := model.SelectGoods()
	if err != nil {
		log.Println("get goods error")
	}
	for _, v := range goods {
		item := newItem(v.ID, v.Name, v.Num, v.Num)
		//覆盖
		ItemMap[item.ID] = item
	}
}

//数据库加载到Map
func set_inMap() {

	//假设每天凌晨12点加载
	t := time.Now()
	//获取明天的时间
	next := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
	//设置timer
	timer := time.NewTimer(next.Sub(t))
	//如果啥子商品都没用，就读取
	<-timer.C
	//添加
	initMap()
	//自己调用自己
	set_inMap()

}

//构造函数
func newItem(id uint, name string, total int, left int) *Item {
	return &Item{
		ID:        id,
		Name:      name,
		Total:     total,
		Left:      left,
		IsSoldOut: false,
		leftCh:    make(chan int),
		sellCh:    make(chan int),
	}
}

func getItem(itemId uint) *Item {
	return ItemMap[itemId]
}

func order() {
	for {
		//传一个用户
		user := <-OrderChan
		//获取到对应的商品
		item := getItem(user.GoodsId)
		//购买
		item.SecKilling(user.UserId)
	}
}

func (item *Item) SecKilling(userId string) {

	item.Lock.Lock()
	defer item.Lock.Unlock()
	// 等价
	// var lock = make(chan struct{}, 1}
	// lock <- struct{}{}
	// defer func() {
	// 		<- lock
	// }

	//如果买完了
	if item.IsSoldOut {
		return
	}
	//买一个
	item.BuyGoods(1)
	//执行
	MakeOrder(userId, item.ID, 1)

}

// 定时下架
//这个下架。。。
func (item *Item) OffShelve() {
	beginTime := time.Now()
	// 获取第二天时间
	//nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品下架的时间
	//offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
	offShelveTime := beginTime.Add(time.Second * 5)
	timer := time.NewTimer(offShelveTime.Sub(beginTime))

	<-timer.C
	delete(ItemMap, item.ID)
	close(item.done)

}

// 出售商品
func (item *Item) SalesGoods() {
	for {
		select {
		//出售
		case num := <-item.sellCh:
			if item.Left -= num; item.Left <= 0 {
				item.IsSoldOut = true
			}
		//剩下的
		case item.leftCh <- item.Left:
		//下架了
		case <-item.Done():
			log.Println("我自闭了")

			return
		}
	}
}

func (item *Item) Done() <-chan struct{} {
	if item.done == nil {
		item.done = make(chan struct{})
	}
	d := item.done
	return d
}

//监视出售商品
func (item *Item) Monitor() {
	go item.SalesGoods()
}

// 获取剩余库存 -.-这个目前未用
func (item *Item) GetLeft() int {
	var left int
	left = <-item.leftCh
	return left
}

// 购买商品
func (item *Item) BuyGoods(num int) {
	item.sellCh <- num
}

func InitService() {
	//初始化ItemMap
	initMap()
	//开启定时器，每天加载到map
	go set_inMap()
	//如果在凌晨12点开启服务，，这里会gg ItemMap是普通的map，在多个协程中会发生冲突
	for _, item := range ItemMap {
		//监听
		item.Monitor()
		//协程进行计时，并定时刷新
		go item.OffShelve()
	}
	for i := 0; i < 10; i++ {
		//用户命名，开10个协程
		go order()
	}

}
