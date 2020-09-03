package main

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/elgs/gosqljson"
	_ "github.com/mattn/go-adodb"
)

const INSERT_LOG = `insert into tx_pay_list(store,pos,orderid,dtm,amt,articlename,aid,num,price,status)
	values(@store,@pos,@order,@dtm,@amt,@name,@aid,@num,@price,@status)`

const UPDATE_LOG = `update tx_pay_list set status = @status,uptdtm = @dtm where orderid = @order`

const SELECCT_LIST = `select store,pos,orderid,dtm,amt,articlename,aid,num,price,isnull(num,0) * isnull(price,0) as amt1, status from tx_pay_list where status='1' and dtm >=@dt1 and dtm <=@dt2 order by dtm`

const SELECCT_LIST_SUM = `select store,pos,orderid,min(dtm) dtm,amt,sum(num) as num from tx_pay_list 
	where status='1' and dtm >=@dt1 and dtm <=@dt2 
	group by store,pos,orderid,amt
	order by orderid`

const SELECT_ESCALE = `select a1,article_name,p,a2,q,falg,supplier_id,counter_id,segregate_id 
	from v_esacle_article_list
	`

const SELECT_GROUPSALE_D1 = `select 
groupsale_bi as group_sale_bi,groupsale_d.aid,unit,article_name,barcode,q as group_q,favour_a,tax_gross,memo,groupsale_d.retail_price as sett_price,groupsale_d.retail_a as sett_a,specification,duty_paragraph
from groupsale_d,article
where groupsale_d.aid = article.aid 
and groupsale_bi in(select groupsale_bi from groupsale_m where groupsale_bi in (`

const SELECT_GROUPSALE_M1 = `select 
groupsale_bi as group_sale_bi,folio_ref,customer_id,input_dt as draw_dt,audite_dt,auditer,
memo,customer_id as member_id ,(select SUM(retail_a) from GROUPSALE_D where GROUPSALE_D.GROUPSALE_BI = GROUPSALE_M.GROUPSALE_BI) as sett_a from groupsale_m `

const SELECT_GROUPSALE_D = `select 
group_sale_bi,group_sale_d.aid,unit,article_name,barcode,group_q,bale_q,entry_price,tax_entry_price,entry_price_a,tax_a,
tax_retail_a,favour_a,tax_gross,memo,sett_price,sett_a,specification,duty_paragraph
from group_sale_d,article
where group_sale_d.aid = article.aid 
and group_sale_bi in(select group_sale_bi from group_sale_m where group_sale_bi in (`

const SELECT_GROUPSALE_M = `select 
group_sale_bi,folio_ref,customer_id,draw_dt,drawer,audite_dt,auditer,sys_dt,
settlement_status,flag,memo,addr_id,member_id,(select sum(sett_a) from group_sale_d where group_sale_d.group_sale_bi = group_sale_m.group_sale_bi) as sett_a from group_sale_m `

func getDB1() (db *sql.DB) {

	//dbconn := os.Getenv("RESTDB")
	//dbconn := "server=192.168.0.72;port=1433;user id=sa;password=soft;database=th_ol025"
	dbconn := "Provider=SQLOLEDB;Initial Catalog=dhy_zb;Data Source=126.10.9.242,1433;user id=sa;password=soft;database=dhy_zb"

	envdbconn := os.Getenv("RESTDB1")
	if envdbconn != "" {
		dbconn = envdbconn
	}

	db, err := sql.Open("adodb", dbconn)

	if err != nil {
		ErrorLog("Open database error: %s\n", err)
	}
	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		ErrorLog("db ping: ", err)
		//log.Fatal(err)
		return nil
	} else {
		//originalDb = db
		return db
	}
}

func getDB() (db *sql.DB) {

	//dbconn := os.Getenv("RESTDB")
	//dbconn := "server=192.168.0.72;port=1433;user id=sa;password=soft;database=th_ol025"
	dbconn := "Provider=SQLOLEDB;Initial Catalog=dhy;Data Source=126.10.9.99,1433;user id=sa;password=newer0710;database=dhy"
	//dbconn := "Provider=SQLOLEDB;Initial Catalog=th_ol025;Data Source=192.168.0.72,1433;user id=sa;password=soft;database=th_ol025"
	//dbconn := "server=10.120.99.2;port=1433;user id=sa;password=soft;database=th_hq_new"
	//dbconn := "server=dd.eastime.top;port=11433;user id=sa;password=cepDemo123;database=cepone"
	//db, err := sql.Open("sqlserver", dbconn)

	envdbconn := os.Getenv("RESTDB")
	if envdbconn != "" {
		dbconn = envdbconn
	}

	db, err := sql.Open("adodb", dbconn)
	//db, err := sql.Open("mssql", dbconn)

	if err != nil {
		ErrorLog("Open database error: %s\n", err)
	}
	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		ErrorLog("db ping: ", err)
		//log.Fatal(err)
		return nil
	} else {
		//originalDb = db
		return db
	}
}

func QueryList(dtm1 string, dtm2 string) ([]map[string]string, error) {

	db := getDB()
	defer db.Close()
	params := make([]interface{}, 0)
	params = append(params, sql.Named("@dt1", dtm1))
	params = append(params, sql.Named("@dt2", dtm2))
	return QueryDbMap(db, SELECCT_LIST, params)
}

func QueryListSum(dtm1 string, dtm2 string) ([]map[string]string, error) {

	db := getDB()
	defer db.Close()
	params := make([]interface{}, 0)
	params = append(params, sql.Named("@dt1", dtm1))
	params = append(params, sql.Named("@dt2", dtm2))
	return QueryDbMap(db, SELECCT_LIST_SUM, params)
}

//db1 groupsale
func QueryListGroupM1(dtm1 string, dtm2 string, bill, cus string) ([]map[string]string, error) {

	db := getDB1()
	defer db.Close()
	//params := make([]interface{}, 0)
	//params = append(params, sql.Named("dt1", dtm1))
	//params = append(params, sql.Named("dt2", dtm2))
	str := SELECT_GROUPSALE_M1 + "where input_dt >='" + dtm1 + "' and input_dt <='" + dtm2 + "'"
	str = str + " and  isnull(customer_id,'') like '" + cus + "' and groupsale_bi like '" + bill + "'"
	return QueryDbMap(db, str, nil)
}

func QueryListGroupD1(bill string) ([]map[string]string, error) {

	db := getDB1()
	defer db.Close()
	//params := make([]interface{}, 0)
	//params = append(params, sql.Named("bill", bill))
	str := SELECT_GROUPSALE_D1 + bill + ")) order by groupsale_bi"
	return QueryDbMap(db, str, nil)
}

//db group_sale
func QueryListGroupM(dtm1 string, dtm2 string, bill, cus string) ([]map[string]string, error) {

	db := getDB()
	defer db.Close()
	//params := make([]interface{}, 0)
	//params = append(params, sql.Named("dt1", dtm1))
	//params = append(params, sql.Named("dt2", dtm2))
	str := SELECT_GROUPSALE_M + "where draw_dt >='" + dtm1 + "' and draw_dt <='" + dtm2 + "'"
	str = str + " and  isnull(customer_id,'') like '" + cus + "' and group_sale_bi like '" + bill + "'"
	return QueryDbMap(db, str, nil)
}

func QueryListGroupD(bill string) ([]map[string]string, error) {

	db := getDB()
	defer db.Close()
	//params := make([]interface{}, 0)
	//params = append(params, sql.Named("bill", bill))
	str := SELECT_GROUPSALE_D + bill + ")) order by group_sale_bi"
	return QueryDbMap(db, str, nil)
}

func QueryListScale(seg string, count string, supplier string) ([]map[string]string, error) {

	db := getDB()
	defer db.Close()
	//params := make([]interface{}, 0)
	//
	//params = append(params, supplier)
	//params = append(params, count)
	//params = append(params, seg)

	//params = append(params, sql.Named("segregate", seg))
	//params = append(params, sql.Named("count", count))
	//params = append(params, sql.Named("supplier", supplier))
	str := SELECT_ESCALE + " where supplier_id like '" + supplier + "' and counter_id like '" + count + "' and segregate_id like '" + seg + "' "
	return QueryDbMap(db, str, nil)
}

func QueryDbMap(db *sql.DB, strsql string, args []interface{}) ([]map[string]string, error) {
	if db == nil {
		return nil, errors.New("db error")
	}
	DebugLog("query sql:", strsql, args)
	data, err := gosqljson.QueryDbToMap(db, "default", strsql, args...)
	if err != nil {
		ErrorLog(err)
	}
	DebugLog("query result:", len(data))
	return data, err
}

func UpdateStatus(order string, store string, dtm string, status string) error {
	db := getDB()

	if db == nil {
		return errors.New("服务端系统异常，数据库连接失败")
	}
	defer db.Close()

	stmt, err := db.Prepare(UPDATE_LOG)
	defer stmt.Close()

	if err != nil {
		ErrorLog("Sql error,", UPDATE_LOG)
		return err
	}

	result, err := stmt.Exec(
		sql.Named("store", store),
		sql.Named("order", order),
		sql.Named("dtm", dtm),
		sql.Named("status", status))
	if err != nil {
		ErrorLog("Update failed:", err.Error())
		return err
	}
	rtn, _ := result.RowsAffected()
	DebugLog("Update result:", rtn)

	return nil
}

func SaveList(items []OrderItem, order string, store string, pos string, amt int, dtm string) error {

	db := getDB()
	if db == nil {
		return errors.New("服务端系统异常，数据库连接失败")
	}
	ctx, errRtn := db.Begin()

	if errRtn != nil {
		ErrorLog("begin transaction error")
		return errRtn
	}

	defer func() {
		if errRtn != nil {
			ErrorLog(errRtn.Error())
			errRtn = ctx.Rollback()
		}
		db.Close()
	}()

	stmt, err := db.Prepare(INSERT_LOG)
	if err != nil {
		ErrorLog("Sql error,", INSERT_LOG)
		errRtn = err
		return errRtn
	}
	for _, item := range items {
		result, err := stmt.Exec(
			sql.Named("store", store),
			sql.Named("pos", pos),
			sql.Named("order", order),
			sql.Named("dtm", dtm),
			sql.Named("amt", amt/100.0),
			sql.Named("name", item.Goodsname),
			sql.Named("num", item.Goodsnum),
			sql.Named("price", item.Goodsprice),
			sql.Named("aid", item.goodsID),
			sql.Named("status", "0"))
		//	store,pos,order,dtm,amt,item.Goodsname,item.Goodsnum,item.Goodsprice,"0")
		if err != nil {
			ErrorLog("Insert failed:", err.Error())
			errRtn = err
			return errRtn
		}
		rtn, _ := result.LastInsertId()
		DebugLog("exec insert, newID:", rtn)
	}

	//fmt.Println(res.RowsAffected())
	defer stmt.Close()

	errRtn = ctx.Commit()

	return errRtn
}
