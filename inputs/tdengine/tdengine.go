package tdengine

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"database/sql"

	"flashcat.cloud/categraf/config"
	"flashcat.cloud/categraf/inputs"
	"flashcat.cloud/categraf/pkg/conv"
	"flashcat.cloud/categraf/types"

	_ "github.com/taosdata/driver-go/v3/taosRestful"
)

const inputPrefixName = "tdengine"

type QueryConfig struct {
	Mesurement    string          `toml:"mesurement"`
	LabelFields   []string        `toml:"label_fields"`
	MetricFields  []string        `toml:"metric_fields"`
	FieldToAppend string          `toml:"field_to_append"`
	Timeout       config.Duration `toml:"timeout"`
	Request       string          `toml:"request"`
}

type TDengine struct {
	config.PluginConfig
	Instances []*Instance `toml:"instances"`
}

func init() {
	inputs.Add("tdengine", func() inputs.Input {
		return &TDengine{}
	})
}

func (a *TDengine) Clone() inputs.Input {
	return &TDengine{}
}

func (a *TDengine) Name() string {
	return inputPrefixName
}

type Instance struct {
	config.InstanceConfig
	Queries  []QueryConfig `toml:"queries"`
	Address  string        `toml:"address"`
	User     string        `toml:"user"`
	Password string        `toml:"password"`
}

// 多实例
func (m *TDengine) GetInstances() []inputs.Instance {
	ret := make([]inputs.Instance, len(m.Instances))
	for i := 0; i < len(m.Instances); i++ {
		ret[i] = m.Instances[i]
	}
	return ret
}

func (ins *Instance) Gather(slist *types.SampleList) {

	var taosurl = fmt.Sprintf("%s:%s@http(%s)/", ins.User, ins.Password, ins.Address)
	taos, err := sql.Open("taosRestful", taosurl)

	if err != nil {
		slist.PushFront(types.NewSample(inputPrefixName, "health", 0))
		log.Println("failed to connect TDengine, err:", err)
		return
	}

	defer taos.Close()

	err = ins.CheckHealth(slist, taos)
	if err != nil {
		return
	}

	ins.getApps(taos, slist)
	ins.getDNODES(taos, slist)
	ins.getMNODES(taos, slist)
	ins.getInsDatabases(taos, slist)
	ins.gatherCustomQueries(slist, taos)
}

// /SHOW CLUSTER;
func (ins *Instance) CheckHealth(slist *types.SampleList, taos *sql.DB) error {

	//正常返回1
	sqlstr := "SHOW CLUSTER;"
	rows, err := taos.Query(sqlstr)

	if err != nil {
		slist.PushFront(types.NewSample(inputPrefixName, "health", 0))
		return err
	}

	defer rows.Close()

	for rows.Next() {

		var r struct {
			id          string
			name        string
			uptime      int
			create_time string
			version     string
			expire_time sql.NullString
		}

		err := rows.Scan(&r.id, &r.name, &r.uptime, &r.create_time, &r.version, &r.expire_time)

		if err != nil {
			log.Println("getMNODES scan error:\n", err)
			continue
		}
		slist.PushFront(types.NewSample(inputPrefixName, "uptime", r.uptime))
	}

	slist.PushFront(types.NewSample(inputPrefixName, "health", 1))

	return err
}

// 当前连接客户端
// SHOW CONNECTIONS;
// SHOW APPS;(更详细)
func (ins *Instance) getApps(taos *sql.DB, slist *types.SampleList) {
	//正常返回1
	sql := "SHOW APPS;"
	rows, err := taos.Query(sql)

	defer rows.Close()

	if err != nil {
		log.Println("getApps error:\n", err)
		return
	}

	for rows.Next() {

		var r struct {
			app_id       string
			ip           string
			pid          string
			name         string
			start_time   string
			insert_req   int
			insert_row   string
			insert_time  int
			insert_bytes string
			fetch_bytes  string
			query_time   int
			slow_query   int
			total_req    int
			current_req  int
			last_access  string
		}

		err := rows.Scan(&r.app_id, &r.ip, &r.pid, &r.name, &r.start_time, &r.insert_req, &r.insert_row,
			&r.insert_time,
			&r.insert_bytes,
			&r.fetch_bytes,
			&r.query_time,
			&r.slow_query,
			&r.total_req,
			&r.current_req,
			&r.last_access)

		if err != nil {
			log.Println("getApps error:\n", err)
			return
		}

		labels := map[string]string{}
		labels["ip"] = r.ip
		labels["name"] = r.name
		labels["insert_row"] = r.insert_row
		labels["start_time"] = r.start_time
		slist.PushFront(types.NewSample(inputPrefixName, "conn_clients", 1, labels))
	}
}

// 获取 MNODES 信息
func (ins *Instance) getMNODES(taos *sql.DB, slist *types.SampleList) {

	sql := "SHOW MNODES;"
	rows, err := taos.Query(sql)
	if err != nil {
		log.Println("getMNODES err:", err)
		return
	}

	defer rows.Close()

	metric_str := "mnodes_role_status"

	for rows.Next() {

		var r struct {
			id          int
			endpoint    string
			role        string
			status      string
			create_time string
		}

		err := rows.Scan(&r.id, &r.endpoint, &r.role, &r.status, &r.create_time)

		if err != nil {
			fmt.Println("getMNODES scan error:\n", err)
			return
		}

		labels := map[string]string{}
		labels["role"] = r.role
		labels["endpoint"] = r.endpoint
		labels["status"] = r.status
		status := 1
		if r.role == "offline" {
			status = 0
		}
		slist.PushFront(types.NewSample(inputPrefixName, metric_str, status, labels))
	}
}

// 获取所有表元数据
func (ins *Instance) getInsDatabases(taos *sql.DB, slist *types.SampleList) {
	sqlstr := "select name,ntables,`vgroups`,`replica`,`keep`,status from information_schema.ins_databases;"
	rows, err := taos.Query(sqlstr)
	if err != nil {
		log.Println("getInsDatabases err:", err)
		return
	}

	defer rows.Close()

	metric_str := "ins_databases"

	for rows.Next() {

		var r struct {
			name    string
			ntables string
			vgroups sql.NullString
			replica sql.NullString
			keep    sql.NullString
			status  string
		}

		err := rows.Scan(&r.name, &r.ntables, &r.vgroups, &r.replica, &r.keep, &r.status)

		if err != nil {
			fmt.Println("getInsDatabases scan error:\n", err)
			continue
		}

		if r.name == "information_schema" || r.name == "performance_schema" {
			continue
		}

		labels := map[string]string{}
		labels["name"] = r.name
		labels["ntables"] = r.ntables
		labels["vgroups"] = r.vgroups.String
		labels["replica"] = r.replica.String
		labels["keep"] = r.keep.String
		labels["status"] = r.status
		status := 0
		if r.status == "ready" {
			status = 1
		}
		slist.PushFront(types.NewSample(inputPrefixName, metric_str, status, labels))
	}
}

// 获取 DNODES 信息
func (ins *Instance) getDNODES(taos *sql.DB, slist *types.SampleList) {

	sql := "SHOW DNODES;"
	rows, err := taos.Query(sql)
	if err != nil {
		fmt.Println("getDNODES err:", err)
	}

	defer rows.Close()

	metric_str := "dnodes_status"

	for rows.Next() {

		var r struct {
			id             int
			endpoint       string
			vnodes         string
			support_vnodes string
			status         string
			create_time    string
			note           string
		}

		err := rows.Scan(&r.id, &r.endpoint, &r.vnodes, &r.support_vnodes, &r.status, &r.create_time, &r.note)
		if err != nil {
			fmt.Println("getDNODES scan error:\n", err)
			return
		}

		labels := map[string]string{}
		labels["vnodes"] = r.vnodes
		labels["support_vnodes"] = r.support_vnodes
		labels["endpoint"] = r.endpoint
		labels["status"] = r.status
		status := 0
		if r.status == "ready" {
			status = 1
		}
		slist.PushFront(types.NewSample(inputPrefixName, metric_str, status, labels))
	}
}

func (ins *Instance) gatherCustomQueries(slist *types.SampleList, db *sql.DB) {
	wg := new(sync.WaitGroup)
	defer wg.Wait()

	for i := 0; i < len(ins.Queries); i++ {
		wg.Add(1)
		go ins.gatherOneQuery(slist, db, wg, ins.Queries[i])
	}
}

func (ins *Instance) gatherOneQuery(slist *types.SampleList, db *sql.DB, wg *sync.WaitGroup, query QueryConfig) {
	defer wg.Done()

	timeout := time.Duration(query.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query.Request)

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("E! query timeout, request:", query.Request)
		return
	}

	if err != nil {
		log.Println("E! failed to query:", err)
		return
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println("E! failed to get columns:", err)
		return
	}

	for rows.Next() {

		columns := make([]sql.RawBytes, len(cols))

		columnPointers := make([]interface{}, len(cols))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			log.Println("E! failed to scan:", err)
			return
		}

		row := make(map[string]string)
		for i, colName := range cols {
			val := columnPointers[i].(*sql.RawBytes)
			row[strings.ToLower(colName)] = string(*val)
		}

		if err = ins.parseRow(row, query, slist); err != nil {
			log.Println("E! failed to parse row:", err, "sql:", query.Request)
		}
	}
}

func (ins *Instance) parseRow(row map[string]string, query QueryConfig, slist *types.SampleList) error {
	labels := map[string]string{}

	for _, label := range query.LabelFields {
		labelValue, has := row[label]
		if has {
			labels[label] = strings.Replace(labelValue, " ", "_", -1)
		}
	}

	for _, column := range query.MetricFields {
		value, err := conv.ToFloat64(row[column])
		if err != nil {
			log.Println("E! failed to convert field:", column, "value:", value, "error:", err)
			return err
		}

		if query.FieldToAppend == "" {
			slist.PushFront(types.NewSample(inputPrefixName, column, value, labels))
		} else {
			suffix := cleanName(row[query.FieldToAppend])
			slist.PushFront(types.NewSample(inputPrefixName, query.Mesurement+"_"+suffix+"_"+column, value, labels))
		}
	}

	return nil
}

func cleanName(s string) string {
	s = strings.Replace(s, " ", "_", -1) // Remove spaces
	s = strings.Replace(s, "(", "", -1)  // Remove open parenthesis
	s = strings.Replace(s, ")", "", -1)  // Remove close parenthesis
	s = strings.Replace(s, "/", "", -1)  // Remove forward slashes
	s = strings.Replace(s, "*", "", -1)  // Remove asterisks
	s = strings.Replace(s, "%", "percent", -1)
	s = strings.ToLower(s)
	return s
}
