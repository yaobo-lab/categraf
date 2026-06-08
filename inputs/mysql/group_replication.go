package mysql

import (
	"database/sql"
	"log"
	"strings"

	"flashcat.cloud/categraf/types"
)

func (ins *Instance) gatherGroupRelication(slist *types.SampleList, db *sql.DB, globalTags map[string]string) {
	rows, err := db.Query(SQL_GROUP_REPLICATION_MEMBER)
	if err != nil {
		log.Println("E! failed to query group relication innodb status:", err)
		return
	}

	defer rows.Close()

	for rows.Next() {

		var model struct {
			channel_name string
			member_role  string
			member_state string
		}

		err = rows.Scan(&model.channel_name, &model.member_state, &model.member_role)

		if err != nil {
			continue
		}

		model.member_role = strings.ToLower(model.member_role)
		model.member_state = strings.ToLower(model.member_state)

		//PRIMARY
		is_primary := 0
		if model.member_role == "primary" {
			is_primary = 1
		}
		slist.PushFront(types.NewSample(inputName, "replication_group_members_is_primary", is_primary, globalTags))

		status := 0
		switch model.member_state {
		//正常工作状态
		case "online":
			status = 1
			//表示节点正在加入组中，这个状态有可能是正在同步数据，也有可能是正在和主节点发 生通信
			//如果长期处于这个状态，往往是 host 没配，需要检查下 host 配置
		case "recovering":
			status = 2
		case "offline":
			//表示这个节点的组复制插件已经加载
			status = 3
		case "unreachable":
			//表示经过仲裁，某个节点已经崩溃或者不可访问
			status = 4
		}
		slist.PushFront(types.NewSample(inputName, "replication_group_members_state", status, globalTags))
		return
	}
}
