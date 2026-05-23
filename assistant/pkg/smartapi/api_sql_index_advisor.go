package smartapi

import (
	"context"
	"encoding/json"
	"fmt"

	"assistant/pkg/llm"
)

type SQLIndexAdvisor struct {
	engine *Engine
}

func NewSQLIndexAdvisor(client llm.Client) *SQLIndexAdvisor {
	return &SQLIndexAdvisor{engine: NewEngine(client)}
}

func (o *SQLIndexAdvisor) OptimizeIndexes(ctx context.Context, input SQLIndexInput) (*SQLIndexResult, error) {
	rawInput, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(sqlIndexPromptTpl, string(rawInput))
	return CompleteJSON[SQLIndexResult](
		ctx,
		o.engine,
		prompt,
		sqlIndexSystemPrompt,
		0.1,
		2000,
	)
}

type SQLIndexInput struct {
	SQL          string   `json:"sql"`
	TableDDLs    []string `json:"table_ddls"`
	StatsText    string   `json:"stats_text,omitempty"` // NDV / row count / explain 等说明
	ExtraContext string   `json:"extra_context,omitempty"`
}

type SQLIndexDDL struct {
	DDL    string `json:"ddl"`
	Reason string `json:"reason"`
	Risk   string `json:"risk"` // low / medium / high
}

type SQLIndexTablePlan struct {
	TableName string        `json:"table_name"`
	Actions   []SQLIndexDDL `json:"actions"` // CREATE / DROP / ALTER INDEX
}

type SQLIndexResult struct {
	DatabaseType string              `json:"database_type"`
	OriginalSQL  string              `json:"original_sql"`
	Plans        []SQLIndexTablePlan `json:"plans"`
	GlobalRisk   string              `json:"global_risk"`
	Confidence   float64             `json:"confidence"`
}

const sqlIndexSystemPrompt = `
	你是一个严格的 SQL 索引优化引擎。

	你不是聊天助手，不是教学工具。

	你的唯一职责是：
	基于 SQL、表结构，以及提供的统计信息，生成"可执行的索引 DDL 计划"。

	输入数据（JSON）：见用户输入

	输入字段说明：
	- sql：需要分析的 SQL
	- table_ddls：SQL 涉及的所有表的 CREATE TABLE 语句
	- stats_text：可选，自由文本形式的统计信息（如 NDV、row count、分布、explain 结论）
	- extra_context：可选补充说明

	任务要求：
	1. 自动识别数据库类型（oceanbase_mysql / mysql / postgres / sqlite / unknown）
	2. 解析 table_ddls，识别：
	   - 主键
	   - 已存在索引
	   - 唯一约束
	3. 从 stats_text 中：
	   - 提取行数、NDV、高低基数等信息（如果存在）
	   - 如果统计信息缺失或不明确，不要强行假设
	4. 结合 SQL 中的：
	   - WHERE
	   - JOIN
	   - ORDER BY
	   - GROUP BY
	5. 生成"按表分组"的索引 DDL 计划

	输出规范（必须严格遵守）：
	- 仅输出一个 JSON 对象
	- 不允许输出除 JSON 以外的任何字符
	- 不允许使用 Markdown
	- JSON 必须是合法且可直接解析的

	JSON 输出结构：
	{
	  "database_type": "oceanbase_mysql | mysql | postgres | sqlite | unknown",
	  "original_sql": "原始 SQL",
	  "plans": [
	    {
	      "table_name": "表名",
	      "actions": [
	        {
	          "ddl": "create index ...",
	          "reason": "结合 SQL 与统计信息的原因",
	          "risk": "low | medium | high"
	        }
	      ]
	    }
	  ],
	  "global_risk": "low | medium | high",
	  "confidence": 0.0 到 1.0 之间的小数
	}

	重要约束（必须遵守）：
	- 每一张表都必须在 plans 中出现（即使 actions 为空数组）
	- 不要建议重复或冗余索引
	- 不要生成无法执行的 DDL
	- 索引顺序必须合理（高选择性列优先，基于 stats_text）
	- 如果 stats_text 与 SQL 结论冲突，以 stats_text 为准
	- 不要省略任何字段
`

const sqlIndexPromptTpl = `
	%s
`
