package ticker

import (
	"testing"
	"time"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
)

func TestTicker(t *testing.T) {
	ticker := NewManager()
	Convey("TestTickerManager", t, func() {
		var result int
		key := "test1"
		Convey("AddTicker", func() {
			// 测试没有ticker 去判断是否存在
			So(ticker.HasTicker(key), ShouldEqual, false)
			// 取消一个不存在的定时器
			ticker.CancelTicker(key)
			// 添加ticker
			ticker.AddTicker(key, 1 * time.Second, func() {
				t.Logf("ticker: %s", key)
				result++
			})
			So(ticker.HasTicker(key), ShouldEqual, true)
			time.Sleep(2 * time.Second)
		})
		Convey("CancelTimer", func() {
			// 取消ticker
			ticker.CancelTicker(key)
			result = 0
			So(ticker.HasTicker(key), ShouldEqual, false)
			time.Sleep(2 * time.Second)
			So(result == 0, ShouldEqual, true)
		})
		Convey("ResetTimer", func() {
			// 添加多个ticker，然后重置
			for i := 1; i < 5; i++ {
				ticker.AddTicker(key + strconv.Itoa(i), 1 * time.Second, func(){
					t.Errorf("should not be execute")
				})
			}
			ticker.Reset()
			for i := 1; i < 5; i++ {
				So(ticker.HasTicker(key + strconv.Itoa(i)), ShouldEqual, false)
			}
		})
	})
}