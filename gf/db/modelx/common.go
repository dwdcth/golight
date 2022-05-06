package modelx

import (
	"context"
	"math"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const Success = 0
const Failed = -1
const SuccessInfo = "操作成功"
const FailedInfo = "操作失败"
const FailedNotFound = "数据不存在"
const SystemError = "系统错误"

type Pagination struct {
	PageIndex   int         `json:"page_index" example:"1"`        //第几页
	PageSize    int         `json:"page_size" example:"10"`        // 页面大小
	IsFirstPage bool        `json:"is_first_page" example:"false"` // 是否是第一页
	IsLastPage  bool        `json:"is_last_page" example:"true"`   // 是否是最后一页
	PageTotal   int         `json:"page_total" example:"10"`       // 总页数
	HasRecords  bool        `json:"has_records" example:"true" `   // 是否有记录
	DataList    interface{} `json:"data_list"`                     // 记录列表
	Total       int         `json:"total" example:"100"`           // 总条数
	// Title       []map[string]string `json:"Title"`                       // 标题列表
}

func (p *Pagination) Offset() int {
	if p.PageIndex <= 0 {
		p.PageIndex = 1
	}
	return (p.PageIndex - 1) * p.PageSize
}
func (p *Pagination) SetData(data interface{}) {
	p.DataList = data
}

//传入页码是否超过总页码
func (p *Pagination) IsOverflow() bool {
	if p.PageIndex > p.PageTotal {
		return true
	}
	return false
}

// NewPagination 新建分页
func NewPagination(total, pageSize, page int, dataList interface{}) *Pagination {
	p := &Pagination{Total: total, PageSize: pageSize, PageIndex: page, DataList: dataList}

	p.PageTotal = int(math.Ceil(float64(total) / float64(pageSize)))
	p.HasRecords = p.PageTotal > 0
	p.IsFirstPage = page == 1
	p.IsLastPage = page == p.PageTotal
	if dataList == nil {
		p.DataList = make([]interface{}, 0)
	}
	return p
}

/*
返回分页结构体
querys  查询 用的query, 如果只有一个query, 则只需要传入一个query
如果有需要则传入两个query,避免count的时候失败, 第一个query是查询的query, 第二个query是count的query
example:

where := g.Map{
	"table.column":      param,
}

queryCount := dao.table.Ctx(ctx).
	LeftJoin("table1", "table1.column1 = table2.column2")
query := queryCount.Fields("table1.*,table2.column3")

return PageResult([]*gdb.Model{query, queryCount}, req.PageIndex, req.PageSize, where, dao.table.Columns.CreatedTime)

*/

func PageResult(ctx context.Context, querys []*gdb.Model, pageIndex int, pageSize int, where interface{}, orderBy string) (*Pagination, error) {
	var query, countQuery *gdb.Model
	if len(querys) == 0 {
		return nil, gerror.New("查询模型不能为空")
	}

	if pageIndex <= 0 || pageSize <= 0 {
		return nil, gerror.New("页码和页面大小必须大于0")
	}

	query = querys[0]
	countQuery = querys[0]
	if len(querys) == 2 {
		countQuery = querys[1]
	}

	count, err := countQuery.Where(where).Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(SystemError)
	}

	p := &Pagination{Total: count, PageSize: pageSize, PageIndex: pageIndex}

	p.PageTotal = int(math.Ceil(float64(count) / float64(pageSize)))
	p.HasRecords = p.PageTotal > 0
	p.IsFirstPage = pageIndex == 1
	p.IsLastPage = pageIndex == p.PageTotal

	if count == 0 {
		return p, nil
	}

	if p.IsOverflow() {
		return p, nil
	}

	dataList, err := query.Where(where).
		Order(orderBy).
		Limit(p.PageSize, p.Offset()).
		All()
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(SystemError)
	}
	p.DataList = dataList
	return p, nil
}

// 简单分页，返回 model供后续处理 和总页数
// ! 注意筛选列 Fields 函数不能在 query 里 ，需要在返回后里筛选
/*
example:
	many := make([]model.PointInfo, 0)
	query := dao.Trend.Ctx(ctx).Distinct().Fields(
		dao.Trend.Columns().PointUuid,
		dao.Trend.Columns().PointId,
		dao.Trend.Columns().EquipmentId,
	).
		WhereLT(dao.Trend.Columns().Trendtime, req.EndTime).
		WhereGTE(dao.Trend.Columns().Trendtime, req.StartTime).
		Where(dao.Trend.Columns().EquipmentId, req.EquipmentId)
	query, sum, err = model.PageModel(query, req.PageSize, req.PageIndex, dao.Trend.Columns().PointUuid+" asc")
	err = query.Scan(&many)
*/
func PageModel(query *gdb.Model, pageSize int, pageIndex int, orderBy string) (*gdb.Model, int, error) {
	if pageSize <= 0 || pageIndex <= 0 {
		return query, 0, gerror.New("页码错误")
	}
	sum, err := query.Count()
	if err != nil {
		return query, -1, err
	}
	if sum == 0 {
		return query, 0, nil
	}
	pageTotal := int(math.Ceil(float64(sum) / float64(pageSize)))
	if pageIndex > pageTotal {
		return query, pageTotal, gerror.New("页码超出范围")
	}
	query = query.Order(orderBy).Limit(pageSize, (pageIndex-1)*pageSize)
	return query, pageTotal, nil
}
