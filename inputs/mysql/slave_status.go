package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"flashcat.cloud/categraf/pkg/tagx"
	"flashcat.cloud/categraf/types"
)

var slaveStatusQuerySuffixes = [3]string{" NONBLOCKING", " NOLOCK", ""}
var replicaStatusQuery = [2]string{"SHOW ALL REPLICAS STATUS", "SHOW REPLICA STATUS"}
var binaryLogsQuery = `SHOW BINARY LOGS`

func columnIndex(slaveCols []string, colName string) int {
	for idx := range slaveCols {
		if slaveCols[idx] == colName {
			return idx
		}
	}
	return -1
}

func columnValue(scanArgs []interface{}, slaveCols []string, colName string) string {
	var columnIndex = columnIndex(slaveCols, colName)
	if columnIndex == -1 {
		return ""
	}
	return string(*scanArgs[columnIndex].(*sql.RawBytes))
}

func (ins *Instance) gatherBinaryLogs(slist *types.SampleList, db *sql.DB, tags map[string]string) error {
	if !ins.GatherBinaryLogs {
		return nil
	}
	// run query
	rows, err := db.Query(binaryLogsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var (
		size      uint64
		count     uint64
		fileSize  uint64
		fileName  string
		encrypted string
	)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	numColumns := len(columns)

	// iterate over rows and count the size and count of files
	for rows.Next() {
		if numColumns == 3 {
			if err := rows.Scan(&fileName, &fileSize, &encrypted); err != nil {
				return err
			}
		} else {
			if err := rows.Scan(&fileName, &fileSize); err != nil {
				return err
			}
		}

		size += fileSize
		count++
	}
	fields := map[string]interface{}{
		"binary_size_bytes":  size,
		"binary_files_count": count,
	}

	slist.PushSamples(inputName, fields, tags)
	return nil
}

func (ins *Instance) gatherReplicaStatus(slist *types.SampleList, db *sql.DB, globalTags map[string]string) error {
	if !ins.GatherReplicaStatus {
		return nil
	}
	var err error
	for _, query := range replicaStatusQuery {
		err = ins.gatherReplicaStatusOnce(slist, db, globalTags, query)
		if err == nil {
			return nil
		}
	}

	log.Println("E! failed to gather replica status:", err)
	return err
}

func (ins *Instance) gatherReplicaStatusOnce(slist *types.SampleList, db *sql.DB, globalTags map[string]string, query string) error {
	// run query
	var rows *sql.Rows
	var err error
	tags := tagx.Copy(globalTags)

	rows, err = db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	fields := make(map[string]interface{})

	// for each channel record
	for rows.Next() {
		// to save the column names as a field key
		// scanning keys and values separately

		// get columns names, and create an array with its length
		cols, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		vals := make([]sql.RawBytes, len(cols))
		valPtrs := make([]interface{}, len(cols))
		// fill the array with sql.Rawbytes
		for i := range vals {
			vals[i] = sql.RawBytes{}
			valPtrs[i] = &vals[i]
		}
		if err := rows.Scan(valPtrs...); err != nil {
			return err
		}

		// range over columns, and try to parse values
		for i, col := range cols {
			colName := col.Name()

			colName = strings.ToLower(colName)

			colValue := vals[i]

			if ins.GatherAllSlaveChannels &&
				(strings.EqualFold(colName, "channel_name") || strings.EqualFold(colName, "connection_name")) {
				// Since the default channel name is empty, we need this block
				channelName := "default"
				if len(colValue) > 0 {
					channelName = string(colValue)
				}
				tags["channel"] = channelName
				continue
			}

			if len(colValue) == 0 {
				continue
			}

			value, err := ins.parseValueByDatabaseTypeName(colValue, col.DatabaseTypeName())
			if err != nil {
				errString := fmt.Errorf("error parsing mysql slave status %q=%q: %w", colName, string(colValue), err)
				log.Println(errString)
				continue
			}

			fields["slave_"+colName] = value
		}
		slist.PushSamples(inputName, fields, tags)

		// Only the first row is relevant if not all slave-channels should be gathered,
		// so break here and skip the remaining rows
		if !ins.GatherAllSlaveChannels {
			break
		}
	}

	return nil
}

func (ins *Instance) parseValueByDatabaseTypeName(value sql.RawBytes, databaseTypeName string) (interface{}, error) {
	if databaseTypeName == "VARCHAR" {
		return string(value), nil
	}
	if bytes.EqualFold(value, []byte("YES")) || bytes.Equal(value, []byte("ON")) {
		return int64(1), nil
	}

	if bytes.EqualFold(value, []byte("NO")) || bytes.Equal(value, []byte("OFF")) {
		return int64(0), nil
	}
	if val, err := strconv.ParseInt(string(value), 10, 64); err == nil {
		return val, nil
	}
	if val, err := strconv.ParseUint(string(value), 10, 64); err == nil {
		return val, nil
	}
	if val, err := strconv.ParseFloat(string(value), 64); err == nil {
		return val, nil
	}
	return nil, fmt.Errorf("unconvertible value: %v", string(value))
}
