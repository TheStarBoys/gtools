package timer

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"time"
	"strconv"
)

func TestTimerManager(t *testing.T) {
	manager := NewManager()
	Convey("TestTimerManager", t, func() {
		Convey("AddTimer", func() {
			// 正常的添加一个定时器
			key := "AddTimer"
			So(manager.HasTimer(key), ShouldEqual, false)
			manager.AddTimer(key, 2 * time.Second, func() {
				t.Logf("key: %s, time out", key)
			})
			So(manager.HasTimer(key), ShouldEqual, true)
			time.Sleep(3 * time.Second)

			// 正常的添加一个定时器，然后再重复添加一次
			manager.AddTimer(key, 1 * time.Second, func() {
				t.Errorf("key: %s, should not be execute", key)
			})
			manager.AddTimer(key, 1 * time.Second, func() {
				t.Logf("key: %s, time out", key)
			})
			time.Sleep(2 * time.Second)
			So(manager.HasTimer(key), ShouldEqual, false)
		})
		Convey("CancelTimer", func() {
			key := "CancelTimer"
			// 正常的取消一个存在的定时器
			manager.AddTimer(key, 2 * time.Second, func() {
				t.Errorf("key: %s, should not be execute", key)
			})
			So(manager.CancelTimer(key), ShouldEqual, true)
			time.Sleep(3 * time.Second)

			// 取消一个不存在的定时器
			So(manager.CancelTimer(key), ShouldEqual, false)
		})
		Convey("OnDoTimer", func() {
			key := "OnDoTimer"
			// 执行一个不存在的定时器，什么也不会发生
			manager.OnDoTimer(key)

			// 正常执行一个定时器
			manager.AddTimer(key, 10 * time.Second, func() {
				t.Logf("key: %s, time out", key)
			})
			manager.OnDoTimer(key)
			// 执行后的定时器应该被清空
			So(manager.HasTimer(key), ShouldEqual, false)
		})
		Convey("IncreaseLeftTime", func() {
			key := "IncreaseLeftTime"
			// 增加一个不存在的定时器时间，什么也不会发生
			manager.IncreaseLeftTime(key, 1 * time.Second)

			// 正常添加一个定时器，然后增加其剩余时间
			manager.AddTimer(key, 1 * time.Second, func() {
				t.Logf("key: %s, time out", key)
			})
			manager.IncreaseLeftTime(key, 1 * time.Second)
			time.Sleep(2 * time.Second)
		})
		Convey("Reset", func() {
			key := "Reset"
			// 添加多个定时器
			for i := 1; i <= 4; i++ {
				manager.AddTimer(key + strconv.Itoa(i), 1 * time.Second, func() {
					t.Errorf("key: %s, should not be execute", key + strconv.Itoa(i))
				})
			}

			manager.Reset()
			for i := 1; i <= 4; i++ {
				So(manager.HasTimer(key + strconv.Itoa(i)), ShouldEqual, false)
			}
			time.Sleep(1 * time.Second)
		})
	})
}
