package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yilin/ai-for-backend/internal/application/dto"
	progressapp "github.com/yilin/ai-for-backend/internal/application/interview_progress"
	questionapp "github.com/yilin/ai-for-backend/internal/application/interview_question"
	applicationapp "github.com/yilin/ai-for-backend/internal/application/resume_application"
	domainquestion "github.com/yilin/ai-for-backend/internal/domain/interview_question"
)

type resumeInterviewsInput struct {
	ResumeVersionID string `json:"resumeVersionId" jsonschema:"resume version ID"`
	Page            int    `json:"page,omitempty"`
	PageSize        int    `json:"pageSize,omitempty"`
}
type workflowsInput struct {
	CurrentStatus string `json:"currentStatus,omitempty"`
	Page          int    `json:"page,omitempty"`
	PageSize      int    `json:"pageSize,omitempty"`
}
type questionListInput struct {
	Keyword       string `json:"keyword,omitempty"`
	Tags          string `json:"tags,omitempty"`
	Difficulty    string `json:"difficulty,omitempty"`
	QuestionType  string `json:"questionType,omitempty"`
	MasteryStatus string `json:"masteryStatus,omitempty"`
	IsFavorite    *bool  `json:"isFavorite,omitempty"`
	Page          int    `json:"page,omitempty"`
	PageSize      int    `json:"pageSize,omitempty"`
}
type questionCreateInput struct {
	Questions []dto.CreateInterviewQuestionRequest `json:"questions" jsonschema:"questions to create"`
}

func NewToolServer(apps *applicationapp.ResumeApplicationService, progresses *progressapp.InterviewProgressService, questions *questionapp.InterviewQuestionService) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "ai4job", Version: "1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "get_resume_interviews", Description: "查看当前用户指定简历版本关联的投递记录和面试进展"}, func(ctx context.Context, _ *mcp.CallToolRequest, in resumeInterviewsInput) (*mcp.CallToolResult, any, error) {
		page, size := normalizePage(in.Page, in.PageSize)
		rows, total, err := apps.List(ctx, UserID(ctx), in.ResumeVersionID, "", "", "", page, size)
		if err != nil {
			return nil, nil, err
		}
		items := make([]map[string]any, 0, len(rows))
		for _, row := range rows {
			progress, err := progresses.List(ctx, UserID(ctx), row.ID)
			if err != nil {
				return nil, nil, err
			}
			items = append(items, map[string]any{"application": row, "progresses": progress})
		}
		return JSONResult(map[string]any{"list": items, "total": total, "page": page, "pageSize": size})
	})
	mcp.AddTool(server, &mcp.Tool{Name: "list_interview_workflows", Description: "查看当前用户全部投递和面试流程状态"}, func(ctx context.Context, _ *mcp.CallToolRequest, in workflowsInput) (*mcp.CallToolResult, any, error) {
		page, size := normalizePage(in.Page, in.PageSize)
		rows, total, err := apps.List(ctx, UserID(ctx), "", "", "", in.CurrentStatus, page, size)
		if err != nil {
			return nil, nil, err
		}
		items := make([]map[string]any, 0, len(rows))
		for _, row := range rows {
			progress, err := progresses.List(ctx, UserID(ctx), row.ID)
			if err != nil {
				return nil, nil, err
			}
			items = append(items, map[string]any{"application": row, "progresses": progress})
		}
		return JSONResult(map[string]any{"list": items, "total": total, "page": page, "pageSize": size})
	})
	mcp.AddTool(server, &mcp.Tool{Name: "list_interview_questions", Description: "查询当前用户题库。difficulty 仅可为 easy、medium、hard；questionType 仅可为 knowledge、scenario、coding、algorithm、project、behavior、discussion；masteryStatus 仅可为 unreviewed、weak、normal、mastered。", InputSchema: questionListInputSchema()}, func(ctx context.Context, _ *mcp.CallToolRequest, in questionListInput) (*mcp.CallToolResult, any, error) {
		page, size := normalizePage(in.Page, in.PageSize)
		rows, total, err := questions.ListQuestions(ctx, UserID(ctx), dto.InterviewQuestionQuery{Keyword: in.Keyword, Tags: in.Tags, Difficulty: domainquestion.Difficulty(in.Difficulty), QuestionType: domainquestion.QuestionType(in.QuestionType), MasteryStatus: domainquestion.MasteryStatus(in.MasteryStatus), IsFavorite: in.IsFavorite, Page: page, PageSize: size})
		if err != nil {
			return nil, nil, err
		}
		return JSONResult(map[string]any{"list": rows, "total": total, "page": page, "pageSize": size})
	})
	mcp.AddTool(server, &mcp.Tool{Name: "create_interview_questions", Description: "向当前用户题库新增 1 至 50 条题目。difficulty 仅可为 easy、medium、hard；questionType 仅可为 knowledge、scenario、coding、algorithm、project、behavior、discussion；masteryStatus 仅可为 unreviewed、weak、normal、mastered。题目将固定归属当前 API Key 用户，source 由服务端写为 mcp。", InputSchema: questionCreateInputSchema()}, func(ctx context.Context, _ *mcp.CallToolRequest, in questionCreateInput) (*mcp.CallToolResult, any, error) {
		for i := range in.Questions {
			in.Questions[i].Source = "mcp"
		}
		rows, err := questions.CreateQuestions(ctx, UserID(ctx), in.Questions)
		if err != nil {
			return nil, nil, err
		}
		return JSONResult(map[string]any{"list": rows})
	})
	return server
}

func questionListInputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type":"object",
		"properties":{
			"keyword":{"type":"string","description":"按题目标题、内容或解析检索"},
			"tags":{"type":"string","description":"按标签筛选"},
			"difficulty":{"type":"string","enum":["easy","medium","hard"],"description":"题目难度：easy（简单）、medium（中等）、hard（困难）"},
			"questionType":{"type":"string","enum":["knowledge","scenario","coding","algorithm","project","behavior","discussion"],"description":"题型：knowledge（知识点）、scenario（场景）、coding（编程）、algorithm（算法）、project（项目）、behavior（行为）、discussion（开放讨论）"},
			"masteryStatus":{"type":"string","enum":["unreviewed","weak","normal","mastered"],"description":"掌握状态：unreviewed（未复习）、weak（薄弱）、normal（一般）、mastered（已掌握）"},
			"isFavorite":{"type":"boolean","description":"是否仅查询收藏题目"},
			"page":{"type":"integer","minimum":1,"description":"页码，默认 1"},
			"pageSize":{"type":"integer","minimum":1,"maximum":100,"description":"每页数量，默认 20，最大 100"}
		},
		"additionalProperties":false
	}`)
}

func questionCreateInputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type":"object",
		"properties":{
			"questions":{
				"type":"array",
				"minItems":1,
				"maxItems":50,
				"description":"要创建的题目列表，单次 1 至 50 条",
				"items":{
					"type":"object",
					"properties":{
						"title":{"type":"string","maxLength":200,"description":"题目标题"},
						"content":{"type":"string","description":"题目内容"},
						"answerAnalysis":{"type":"string","description":"答案或解析"},
						"tags":{"type":"array","items":{"type":"string"},"description":"标签列表，可省略"},
						"difficulty":{"type":"string","enum":["easy","medium","hard"],"description":"题目难度：easy（简单）、medium（中等）、hard（困难）"},
						"questionType":{"type":"string","enum":["knowledge","scenario","coding","algorithm","project","behavior","discussion"],"description":"题型：knowledge（知识点）、scenario（场景）、coding（编程）、algorithm（算法）、project（项目）、behavior（行为）、discussion（开放讨论）"},
						"masteryStatus":{"type":"string","enum":["unreviewed","weak","normal","mastered"],"description":"掌握状态：unreviewed（未复习）、weak（薄弱）、normal（一般）、mastered（已掌握）"},
						"isFavorite":{"type":"boolean","description":"是否收藏，默认 false"},
						"remark":{"type":"string","description":"备注，可省略"}
					},
					"required":["title","content","answerAnalysis","difficulty","questionType","masteryStatus"],
					"additionalProperties":false
				}
			}
		},
		"required":["questions"],
		"additionalProperties":false
	}`)
}

func normalizePage(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
