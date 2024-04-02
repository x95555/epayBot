package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Msg struct {
	Key  string  `json:"key"`
	Data []Order `json:"data"`
}

type Order struct {
	Uid      int       `json:"uid"`
	Name     string    `json:"name"`
	Amount   float64   `json:"amount"`
	OrderNum string    `json:"order_num"`
	Time     time.Time `json:"time"`

	Bill    bool   `json:"bill"`    //是否为结算
	Address string `json:"address"` //结算地址
}

func payOrder(t time.Time, db *sql.DB) []Order {

	//`trade_no` varchar(64) NOT NULL,
	//`uid` int(11) unsigned NOT NULL,
	//`name` varchar(64) NOT NULL,
	//`money` decimal(10,2) NOT NULL,
	//`endtime` datetime DEFAULT NULL,
	//`status` tinyint(1) NOT NULL DEFAULT '0',

	query := `SELECT 
	             uid,
	             name,
	             money AS "Amount",
	             trade_no AS "OrderNum",
	             endtime AS "Time"
	         FROM pay_order 
	         WHERE endtime >= ? 
	         ORDER BY trade_no ASC `

	rows, err := db.Query(query, t)
	if err != nil {
		return nil
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	var orders []Order
	for rows.Next() {
		var o Order
		err = rows.Scan(&o.Uid, &o.Name, &o.Amount, &o.OrderNum, &o.Time)
		if err != nil {
			continue
		}
		orders = append(orders, o)
	}

	return orders
}

// 如果使用易支付pay_settle表没有account字段，注销即可
func getBill(t time.Time, db *sql.DB) []Order {

	//  ALTER TABLE `pay_settle`
	//	ADD COLUMN `realmoney` decimal(10,2) NOT NULL,
	//	CHANGE COLUMN `time` `addtime` datetime DEFAULT NULL,
	//	ADD COLUMN `endtime` datetime DEFAULT NULL,
	//	ADD COLUMN `result` varchar(64) DEFAULT NULL;

	query := `SELECT 
	             uid,
	             realmoney AS "Amount",
	             account AS "Address",
	             endtime AS "Time"
	         FROM pay_settle
	         WHERE endtime >= ? 
	         ORDER BY trade_no ASC `

	rows, err := db.Query(query, t)
	if err != nil {
		return nil
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	var orders []Order
	for rows.Next() {
		var o Order
		err = rows.Scan(&o.Uid, &o.Address, &o.Time)
		if err != nil {
			continue
		}
		o.Bill = true
		orders = append(orders, o)
	}

	return orders

}
