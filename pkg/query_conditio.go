package pkg

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type SearchCondition struct {
	Comparator string      `json:"comparator"`
	Value      interface{} `json:"value"`
	Value2     interface{} `json:"value2"`
}

type QueryCondition struct {
	Page          int       `json:"page"`
	Size          int       `json:"size"`
	Sort          string    `json:"sort"`
	Order         string    `json:"order"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	modelDBFields map[string]map[string]string
}

func NewQueryCondition(modelDBFields map[string]map[string]string) *QueryCondition {
	return &QueryCondition{
		modelDBFields: modelDBFields,
	}
}

// 生成高级过滤条件，支持字段搜索条件定制（Comparator：>、==、<、!=、like等）
func (m *QueryCondition) GetQuery(tableName string, condition map[string]SearchCondition, query *gorm.DB, countQuery *gorm.DB) (err error) {
	m.Page = 0
	m.Size = 0
	m.Sort = ""
	m.Order = ""
	m.StartDate = time.Time{}
	m.EndDate = time.Time{}
	//清理前端传入的没有数据的空条件
	for field, item := range condition {
		val := strings.TrimSpace(cast.ToString(item.Value))
		if val == "" && (strings.ToLower(item.Comparator) == "==" || strings.ToLower(item.Comparator) == "") {
			delete(condition, field)
		}
	}
	//处理分页参数
	if _, hasValue := condition["pageNumber"]; hasValue {
		m.Page = cast.ToInt(condition["pageNumber"].Value)
		delete(condition, "pageNumber")
	}
	if _, hasValue := condition["pageSize"]; hasValue {
		m.Size = cast.ToInt(condition["pageSize"].Value)
		delete(condition, "pageSize")
	}
	if _, hasValue := condition["sort"]; hasValue {
		m.Sort = cast.ToString(condition["sort"].Value)
		delete(condition, "sort")
	}
	if _, hasValue := condition["order"]; hasValue {
		m.Order = cast.ToString(condition["order"].Value)
		delete(condition, "order")
	}
	if _, hasValue := condition["startDate"]; hasValue {
		m.StartDate, _ = time.ParseInLocation("2006-01-02", cast.ToString(condition["startDate"].Value), time.Local)
		delete(condition, "startDate")
	}
	if _, hasValue := condition["endDate"]; hasValue {
		m.EndDate, _ = time.ParseInLocation("2006-01-02", cast.ToString(condition["endDate"].Value), time.Local)
		delete(condition, "endDate")
	}
	if !m.StartDate.IsZero() && !m.EndDate.IsZero() {
		query = query.Where("`create_time` >= ?", m.StartDate).Where("`create_time` <= ?", m.EndDate)
		countQuery = countQuery.Where("`create_time` >= ?", m.StartDate).Where("`create_time` <= ?", m.EndDate)
	}
	DBFieldNames, isExists := m.modelDBFields[tableName]
	if isExists {
		// 常规参数
		for key, value := range condition {
			comparator := strings.TrimSpace(value.Comparator)
			val := strings.TrimSpace(cast.ToString(value.Value))
			fieldName := DBFieldNames[key]
			if fieldName != "" {
				switch comparator {
				case "==":
					query = query.Where(fmt.Sprintf("`%s` = ?", fieldName), val)
					countQuery = countQuery.Where(fmt.Sprintf("`%s` = ?", fieldName), val)
				case "!=":
					query = query.Where(fmt.Sprintf("`%s` <> ?", fieldName), val)
					countQuery = countQuery.Where(fmt.Sprintf("`%s` <> ?", fieldName), val)
				case ">":
					query = query.Where(fmt.Sprintf("`%s` > ?", fieldName), val)
					countQuery = countQuery.Where(fmt.Sprintf("`%s` > ?", fieldName), val)
				case ">=":
					query = query.Where(fmt.Sprintf("`%s` >= ?", fieldName), val)
					countQuery = countQuery.Where(fmt.Sprintf("`%s` >= ?", fieldName), val)
				case "<":
					query = query.Where(fmt.Sprintf("`%s` < ?", fieldName), val)
					countQuery = countQuery.Where(fmt.Sprintf("`%s` < ?", fieldName), val)
				case "<=":
					query = query.Where(fmt.Sprintf("`%s` <= ?", fieldName), val)
					countQuery = countQuery.Where(fmt.Sprintf("`%s` <= ?", fieldName), val)
				case "like":
					query = query.Where(fmt.Sprintf("`%s` like ?", fieldName), fmt.Sprintf("%%%s%%", val))
					countQuery = countQuery.Where(fmt.Sprintf("`%s` like ?", fieldName), fmt.Sprintf("%%%s%%", val))
				case "not like":
					query = query.Where(fmt.Sprintf("`%s` not like ?", fieldName), fmt.Sprintf("%%%s%%", val))
					countQuery = countQuery.Where(fmt.Sprintf("`%s` not like ?", fieldName), fmt.Sprintf("%%%s%%", val))
				case "between":
					val2 := strings.TrimSpace(cast.ToString(value.Value2))
					query = query.Where(fmt.Sprintf("(`%s` >= ? and `%s` <= ?)", fieldName, fieldName), val, val2)
					countQuery = countQuery.Where(fmt.Sprintf("(`%s` >= ? and `%s` <= ?)", fieldName, fieldName), val, val2)
				case "in":
					tmpArray := strings.Split(val, ",")
					query = query.Where(fmt.Sprintf("(`%s` in ?)", fieldName), tmpArray)
					countQuery = countQuery.Where(fmt.Sprintf("(`%s` in ?)", fieldName), tmpArray)
				// case "is null":
				// 	query = query.Where(fmt.Sprintf("%s is null", fieldName))
				// 	countQuery = countQuery.Where(fmt.Sprintf("%s is null", fieldName))
				// case "is not null":
				// 	query = query.Where(fmt.Sprintf("%s is not null", fieldName))
				// 	countQuery = countQuery.Where(fmt.Sprintf("%s is not null", fieldName))
				default:
					query = query.Where(fmt.Sprintf("`%s` = ?", fieldName), val)
					countQuery = countQuery.Where(fmt.Sprintf("`%s` = ?", fieldName), val)
				}
			}
		}
		// 处理排序
		if m.Sort != "" && m.Order != "" {
			fieldName := ""
			sortField := make([]string, 0)
			if strings.Contains(m.Sort, ",") {
				for _, field := range strings.Split(m.Sort, ",") {
					tmpFieldName := fmt.Sprintf("`%s`", DBFieldNames[field])
					if tmpFieldName == "" {
						err = fmt.Errorf("%s不是有效的排序字段", field)
						return
					} else {
						sortField = append(sortField, tmpFieldName)
					}
				}
				fieldName = strings.Join(sortField, ",")
			} else {
				fieldName = fmt.Sprintf("`%s`", DBFieldNames[m.Sort])
				if fieldName == "" {
					err = fmt.Errorf("%s不是有效的排序字段", m.Sort)
				}
			}
			if err == nil {
				if m.Page > 0 && m.Size > 0 {
					query.Order(fmt.Sprintf("%s %s", fieldName, m.Order)).Offset((m.Page - 1) * m.Size).Limit(m.Size)
				} else {
					query.Order(fmt.Sprintf("%s %s", fieldName, m.Order))
				}
			}
		} else {
			if m.Page > 0 && m.Size > 0 {
				query.Offset((m.Page - 1) * m.Size).Limit(m.Size)
			}
		}
	} else {
		query = query.Where(fmt.Sprintf("`%s` = ?", "true"), false)
		countQuery = countQuery.Where(fmt.Sprintf("`%s` = ?", "true"), false)
	}
	return
}
