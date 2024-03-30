// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"tontelebot/model"
)

func newJetton(db *gorm.DB, opts ...gen.DOOption) jetton {
	_jetton := jetton{}

	_jetton.jettonDo.UseDB(db, opts...)
	_jetton.jettonDo.UseModel(&model.Jetton{})

	tableName := _jetton.jettonDo.TableName()
	_jetton.ALL = field.NewAsterisk(tableName)
	_jetton.ID = field.NewInt32(tableName, "id")
	_jetton.Address = field.NewString(tableName, "address")
	_jetton.CreateTime = field.NewTime(tableName, "create_time")
	_jetton.UpdateTime = field.NewTime(tableName, "update_time")

	_jetton.fillFieldMap()

	return _jetton
}

// jetton jetton表
type jetton struct {
	jettonDo jettonDo

	ALL        field.Asterisk
	ID         field.Int32  // id
	Address    field.String // jetton地址
	CreateTime field.Time   // 创建时间
	UpdateTime field.Time   // 更新时间

	fieldMap map[string]field.Expr
}

func (j jetton) Table(newTableName string) *jetton {
	j.jettonDo.UseTable(newTableName)
	return j.updateTableName(newTableName)
}

func (j jetton) As(alias string) *jetton {
	j.jettonDo.DO = *(j.jettonDo.As(alias).(*gen.DO))
	return j.updateTableName(alias)
}

func (j *jetton) updateTableName(table string) *jetton {
	j.ALL = field.NewAsterisk(table)
	j.ID = field.NewInt32(table, "id")
	j.Address = field.NewString(table, "address")
	j.CreateTime = field.NewTime(table, "create_time")
	j.UpdateTime = field.NewTime(table, "update_time")

	j.fillFieldMap()

	return j
}

func (j *jetton) WithContext(ctx context.Context) IJettonDo { return j.jettonDo.WithContext(ctx) }

func (j jetton) TableName() string { return j.jettonDo.TableName() }

func (j jetton) Alias() string { return j.jettonDo.Alias() }

func (j jetton) Columns(cols ...field.Expr) gen.Columns { return j.jettonDo.Columns(cols...) }

func (j *jetton) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := j.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (j *jetton) fillFieldMap() {
	j.fieldMap = make(map[string]field.Expr, 4)
	j.fieldMap["id"] = j.ID
	j.fieldMap["address"] = j.Address
	j.fieldMap["create_time"] = j.CreateTime
	j.fieldMap["update_time"] = j.UpdateTime
}

func (j jetton) clone(db *gorm.DB) jetton {
	j.jettonDo.ReplaceConnPool(db.Statement.ConnPool)
	return j
}

func (j jetton) replaceDB(db *gorm.DB) jetton {
	j.jettonDo.ReplaceDB(db)
	return j
}

type jettonDo struct{ gen.DO }

type IJettonDo interface {
	gen.SubQuery
	Debug() IJettonDo
	WithContext(ctx context.Context) IJettonDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IJettonDo
	WriteDB() IJettonDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IJettonDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IJettonDo
	Not(conds ...gen.Condition) IJettonDo
	Or(conds ...gen.Condition) IJettonDo
	Select(conds ...field.Expr) IJettonDo
	Where(conds ...gen.Condition) IJettonDo
	Order(conds ...field.Expr) IJettonDo
	Distinct(cols ...field.Expr) IJettonDo
	Omit(cols ...field.Expr) IJettonDo
	Join(table schema.Tabler, on ...field.Expr) IJettonDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IJettonDo
	RightJoin(table schema.Tabler, on ...field.Expr) IJettonDo
	Group(cols ...field.Expr) IJettonDo
	Having(conds ...gen.Condition) IJettonDo
	Limit(limit int) IJettonDo
	Offset(offset int) IJettonDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IJettonDo
	Unscoped() IJettonDo
	Create(values ...*model.Jetton) error
	CreateInBatches(values []*model.Jetton, batchSize int) error
	Save(values ...*model.Jetton) error
	First() (*model.Jetton, error)
	Take() (*model.Jetton, error)
	Last() (*model.Jetton, error)
	Find() ([]*model.Jetton, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Jetton, err error)
	FindInBatches(result *[]*model.Jetton, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Jetton) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IJettonDo
	Assign(attrs ...field.AssignExpr) IJettonDo
	Joins(fields ...field.RelationField) IJettonDo
	Preload(fields ...field.RelationField) IJettonDo
	FirstOrInit() (*model.Jetton, error)
	FirstOrCreate() (*model.Jetton, error)
	FindByPage(offset int, limit int) (result []*model.Jetton, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IJettonDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (j jettonDo) Debug() IJettonDo {
	return j.withDO(j.DO.Debug())
}

func (j jettonDo) WithContext(ctx context.Context) IJettonDo {
	return j.withDO(j.DO.WithContext(ctx))
}

func (j jettonDo) ReadDB() IJettonDo {
	return j.Clauses(dbresolver.Read)
}

func (j jettonDo) WriteDB() IJettonDo {
	return j.Clauses(dbresolver.Write)
}

func (j jettonDo) Session(config *gorm.Session) IJettonDo {
	return j.withDO(j.DO.Session(config))
}

func (j jettonDo) Clauses(conds ...clause.Expression) IJettonDo {
	return j.withDO(j.DO.Clauses(conds...))
}

func (j jettonDo) Returning(value interface{}, columns ...string) IJettonDo {
	return j.withDO(j.DO.Returning(value, columns...))
}

func (j jettonDo) Not(conds ...gen.Condition) IJettonDo {
	return j.withDO(j.DO.Not(conds...))
}

func (j jettonDo) Or(conds ...gen.Condition) IJettonDo {
	return j.withDO(j.DO.Or(conds...))
}

func (j jettonDo) Select(conds ...field.Expr) IJettonDo {
	return j.withDO(j.DO.Select(conds...))
}

func (j jettonDo) Where(conds ...gen.Condition) IJettonDo {
	return j.withDO(j.DO.Where(conds...))
}

func (j jettonDo) Order(conds ...field.Expr) IJettonDo {
	return j.withDO(j.DO.Order(conds...))
}

func (j jettonDo) Distinct(cols ...field.Expr) IJettonDo {
	return j.withDO(j.DO.Distinct(cols...))
}

func (j jettonDo) Omit(cols ...field.Expr) IJettonDo {
	return j.withDO(j.DO.Omit(cols...))
}

func (j jettonDo) Join(table schema.Tabler, on ...field.Expr) IJettonDo {
	return j.withDO(j.DO.Join(table, on...))
}

func (j jettonDo) LeftJoin(table schema.Tabler, on ...field.Expr) IJettonDo {
	return j.withDO(j.DO.LeftJoin(table, on...))
}

func (j jettonDo) RightJoin(table schema.Tabler, on ...field.Expr) IJettonDo {
	return j.withDO(j.DO.RightJoin(table, on...))
}

func (j jettonDo) Group(cols ...field.Expr) IJettonDo {
	return j.withDO(j.DO.Group(cols...))
}

func (j jettonDo) Having(conds ...gen.Condition) IJettonDo {
	return j.withDO(j.DO.Having(conds...))
}

func (j jettonDo) Limit(limit int) IJettonDo {
	return j.withDO(j.DO.Limit(limit))
}

func (j jettonDo) Offset(offset int) IJettonDo {
	return j.withDO(j.DO.Offset(offset))
}

func (j jettonDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IJettonDo {
	return j.withDO(j.DO.Scopes(funcs...))
}

func (j jettonDo) Unscoped() IJettonDo {
	return j.withDO(j.DO.Unscoped())
}

func (j jettonDo) Create(values ...*model.Jetton) error {
	if len(values) == 0 {
		return nil
	}
	return j.DO.Create(values)
}

func (j jettonDo) CreateInBatches(values []*model.Jetton, batchSize int) error {
	return j.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (j jettonDo) Save(values ...*model.Jetton) error {
	if len(values) == 0 {
		return nil
	}
	return j.DO.Save(values)
}

func (j jettonDo) First() (*model.Jetton, error) {
	if result, err := j.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Jetton), nil
	}
}

func (j jettonDo) Take() (*model.Jetton, error) {
	if result, err := j.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Jetton), nil
	}
}

func (j jettonDo) Last() (*model.Jetton, error) {
	if result, err := j.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Jetton), nil
	}
}

func (j jettonDo) Find() ([]*model.Jetton, error) {
	result, err := j.DO.Find()
	return result.([]*model.Jetton), err
}

func (j jettonDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Jetton, err error) {
	buf := make([]*model.Jetton, 0, batchSize)
	err = j.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (j jettonDo) FindInBatches(result *[]*model.Jetton, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return j.DO.FindInBatches(result, batchSize, fc)
}

func (j jettonDo) Attrs(attrs ...field.AssignExpr) IJettonDo {
	return j.withDO(j.DO.Attrs(attrs...))
}

func (j jettonDo) Assign(attrs ...field.AssignExpr) IJettonDo {
	return j.withDO(j.DO.Assign(attrs...))
}

func (j jettonDo) Joins(fields ...field.RelationField) IJettonDo {
	for _, _f := range fields {
		j = *j.withDO(j.DO.Joins(_f))
	}
	return &j
}

func (j jettonDo) Preload(fields ...field.RelationField) IJettonDo {
	for _, _f := range fields {
		j = *j.withDO(j.DO.Preload(_f))
	}
	return &j
}

func (j jettonDo) FirstOrInit() (*model.Jetton, error) {
	if result, err := j.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Jetton), nil
	}
}

func (j jettonDo) FirstOrCreate() (*model.Jetton, error) {
	if result, err := j.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Jetton), nil
	}
}

func (j jettonDo) FindByPage(offset int, limit int) (result []*model.Jetton, count int64, err error) {
	result, err = j.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = j.Offset(-1).Limit(-1).Count()
	return
}

func (j jettonDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = j.Count()
	if err != nil {
		return
	}

	err = j.Offset(offset).Limit(limit).Scan(result)
	return
}

func (j jettonDo) Scan(result interface{}) (err error) {
	return j.DO.Scan(result)
}

func (j jettonDo) Delete(models ...*model.Jetton) (result gen.ResultInfo, err error) {
	return j.DO.Delete(models)
}

func (j *jettonDo) withDO(do gen.Dao) *jettonDo {
	j.DO = *do.(*gen.DO)
	return j
}