package ORM

import (
	"GolangWebFramwork/Helper"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
)

// ORMCache保存ORM
var ORMCache map[int]*ORM

// 属性信息
type Field struct {
	FiledName string
	FiledType reflect.Type
}

// model注册信息
type ModelInfo struct {
	ModelType  reflect.Type
	Filed      []Field
	FiledNames []string
	NumField   int
}

// ORM 结构
type ORM struct {
	db               *sql.DB
	modelRespository map[string]ModelInfo
}

// GetOrNewORM获取ORM
func GetOrNewORM() *ORM {
	if len(ORMCache) == 0 {
		ORMCache = make(map[int]*ORM)
		ORMCache[0] = &ORM{nil, make(map[string]ModelInfo)}
		return ORMCache[0]
	}
	return ORMCache[0]
}

// ConnectMtSQLDB 连接mysql数据库
func (orm *ORM) ConnectMySQLDB(dataSource string) {
	var err error
	orm.db, err = sql.Open("mysql", dataSource)
	if err != nil {
		fmt.Println("数据库连接失败")
		return
	}
}

// RegistModel注册模型
func (orm *ORM) RegistModel(model interface{}) {
	modelType := reflect.TypeOf(model)
	modelInfo := ModelInfo{}
	modelInfo.ModelType = modelType
	modelField := Field{}
	for i := 0; i < modelType.NumField(); i++ {
		modelField.FiledName = modelType.Field(i).Name
		modelInfo.FiledNames = append(modelInfo.FiledNames, modelField.FiledName)
		modelField.FiledType = modelType.Field(i).Type
		modelInfo.Filed = append(modelInfo.Filed, modelField)
		modelInfo.NumField++
	}
	if _, ok := orm.modelRespository[modelType.Name()]; !ok {
		orm.modelRespository[strings.ToLower(modelType.Name())] = modelInfo
	}
}

// 数据库同步，根据force来判断是否重置数据库
func (orm *ORM) ResetAndSyncDB(force bool) {
	createTableSQL := "create table if not exists %s(%s)"
	dropTableSQL := "drop table if exists %s"
	showTablesSQL := "show tables"
	descTableSQL := "select column_name, data_type from information_schema.columns where table_name = '%s'"
	dropColumnSQL := "alter table %s drop column %s"
	alterColumnTypeSQL := "alter table %s modify %s %s"
	addColumnSQL := "alter table %s add column %s %s"
	//fmt.Println(tableNames)
	//fmt.Println(tableColsMap)
	// 自动删表建表
	for k, v := range orm.modelRespository {
		if force {
			dropSql := fmt.Sprintf(dropTableSQL, k)
			orm.db.Exec(dropSql)
			//fmt.Println(dropSql)
		}
		var cols []string
		for _, fi := range v.Filed {
			cols = append(cols, fi.FiledName+" "+ Helper.GoTypeToSQLTypeString(fi.FiledType))
		}
		createSql := fmt.Sprintf(createTableSQL, k, strings.Join(cols, ","))
		//fmt.Println(createSql)
		orm.db.Exec(createSql)
	}
	// 获取表结构信息
	rows, err := orm.db.Query(showTablesSQL)
	if err != nil {
		fmt.Println("数据库同步失败")
		return
	}
	var tableNames []string
	var tableName string
	for rows.Next() {
		if err = rows.Scan(&tableName); err != nil {
			fmt.Println("数据库同步失败")
			return
		}
		tableNames = append(tableNames, tableName)
	}
	tableColsMap := make(map[string][]Field)
	for _, name := range tableNames {
		rows, err := orm.db.Query(fmt.Sprintf(descTableSQL, name))
		if err != nil {
			fmt.Println("数据库同步失败")
			return
		}
		var colName string
		var dataType string
		for rows.Next() {
			if err = rows.Scan(&colName, &dataType); err != nil {
				fmt.Println("数据库同步失败")
				return
			}
			dt := Helper.SQLTypeStringToGoType(dataType)
			tableColsMap[name] = append(tableColsMap[name], Field{
				FiledName: colName,
				FiledType: dt,
			})
		}
	}
	// 对不存在的列进行删减（若属性名被改变，则这一列会被删除）
	for k1, v1 := range tableColsMap {
		fieldNames := orm.modelRespository[k1].FiledNames
		for _, col := range v1 {
			flag := false
			for _, fieledName := range fieldNames {
				if fieledName == col.FiledName {
					flag = true
				}
			}
			if !flag {
				orm.db.Exec(fmt.Sprintf(dropColumnSQL, k1, col.FiledName))
			}
		}
	}
	// 对缺少的列进行增加（对被改名的列进行重建）
	for k, v := range orm.modelRespository {
		fields := tableColsMap[k]
		for _, newField := range v.Filed {
			flag := false
			for _, field := range fields {
				if field.FiledName == newField.FiledName {
					flag = true
				}
				if !flag {
					orm.db.Exec(fmt.Sprintf(addColumnSQL, k, newField.FiledName, Helper.GoTypeToSQLTypeString(newField.FiledType)))
				}
			}
		}
	}
	// 对属性变化的列进行修改
	for k, v := range orm.modelRespository {
		fields := tableColsMap[k]
		for _, newField := range v.Filed {
			for _, field := range fields {
				if field.FiledName == newField.FiledName && field.FiledType != newField.FiledType {
					orm.db.Exec(fmt.Sprintf(alterColumnTypeSQL, k, newField.FiledName, Helper.GoTypeToSQLTypeString(newField.FiledType)))
				}
			}
		}
	}
}

// GetModelInfo获取已注册模型的信息
func (orm ORM) GetModelInfo(model interface{}) (ModelInfo, error) {
	modelType := reflect.TypeOf(model)
	if _, ok := orm.modelRespository[strings.ToLower(modelType.Name())]; ok {
		modelInfo := orm.modelRespository[strings.ToLower(modelType.Name())]
		return modelInfo, nil
	}
	return ModelInfo{}, errors.New("此model未经注册")
}



// Query查询抽象
func (orm *ORM) Query(model interface{}) []reflect.Value {
	modelInfo, err := orm.GetModelInfo(model)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	modelType := modelInfo.ModelType
	modelValue := reflect.ValueOf(model)

	cols := modelInfo.FiledNames
	var conditions []string
	for i := 0; i < modelInfo.NumField; i++ {
		if modelValue.Field(i).Interface() != "" && modelValue.Field(i).Interface() != 0 {
			conditions = append(conditions, cols[i]+" = "+fmt.Sprintf("%v", modelValue.Field(i).Interface()))
		}
	}
	var selectSQL string
	if len(conditions) != 0 {
		selectSQL = fmt.Sprintf("select %s from %s where %s", strings.Join(cols, ","), modelType.Name(),
			strings.Join(conditions, " and "))
	} else {
		selectSQL = fmt.Sprintf("select %s from %s", strings.Join(cols, ","), modelType.Name())
	}
	//fmt.Println(selectSQL)
	var rows *sql.Rows

	rows, err = orm.db.Query(selectSQL)
	if err != nil {
		fmt.Println("查询失败")
		return nil
	}
	values := make([]sql.RawBytes, len(cols))
	refs := make([]interface{}, len(cols))
	for i := range refs {
		refs[i] = &values[i]
	}
	var rst []reflect.Value
	var elem reflect.Value
	for rows.Next() {
		if err = rows.Scan(refs...); err != nil {
			return nil
		}
		elem = reflect.New(modelType)
		elem = reflect.Indirect(elem)
		Helper.SetElemValues(&elem, values)
		rst = append(rst, elem)
	}
	return rst
}
