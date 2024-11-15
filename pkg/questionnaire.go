package otlh

import "fmt"

type Questionnaire struct {
	ID        int    `json:"id"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Draft     bool   `json:"draft,omitempty"`
	Note      string `json:"note,omitempty"`
	CreatedBy struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
		Type  string `json:"type,omitempty"`
	} `json:"created_by,omitempty"`
	Questions []Questionn `json:"questions,omitempty"`
}

type Questionn struct {
	ID                  int      `json:"id"`
	QuestionnaireID     int      `json:"questionnaire_id"`
	RecommendedQuestion bool     `json:"recommended_question,omitempty"`
	DisplayOrder        int      `json:"display_order,omitempty"`
	Type                string   `json:"type,omitempty"`
	Required            bool     `json:"required,omitempty"`
	Text                string   `json:"text,omitempty"`
	CreatedAt           string   `json:"created_at,omitempty"`
	UpdatedAt           string   `json:"updated_at,omitempty"`
	Answers             []Answer `json:"answers,omitempty"`
}

type Answer struct {
	ID           int    `json:"id"`
	QuestionID   int    `json:"question_id"`
	Text         string `json:"text,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	DisplayOrder int    `json:"display_order,omitempty"`
}

type AnswerAttributes struct {
	ID           int    `json:"id"`
	Text         string `json:"text,omitempty"`
	DisplayOrder int    `json:"display_order,omitempty"`
	Destroy      string `json:"_destroy,omitempty"`
}

type QuestionsAttributes struct {
	ID                  int                `json:"id"`
	Text                string             `json:"text,omitempty"`
	Type                string             `json:"type,omitempty"`
	DisplayOrder        int                `json:"display_order,omitempty"`
	Destroy             string             `json:"_destroy,omitempty"`
	Required            bool               `json:"required,omitempty"`
	RecommendedQuestion bool               `json:"recommended_question,omitempty"`
	AnswersAttributes   []AnswerAttributes `json:"answers_attributes"`
}

type QuestionnaireAttributes struct {
	Name                string                `json:"name"`
	Note                string                `json:"note,omitempty"`
	QuestionsAttributes []QuestionsAttributes `json:"questions_attributes"`
}

type QuestionnaireRequest struct {
	id int
	Request
}

type Questionnaires []Questionnaire

type QuestionnairesResponse struct {
	DefaultEntityListInfo
	Embedded struct {
		Questionnaires Questionnaires `json:"questionnaires"`
	} `json:"_embedded"`
}

type QuestionnaireRequestBuilder struct {
	*QuestionnaireRequest
}

func (b *QuestionnaireRequestBuilder) WithID(id int) *QuestionnaireRequestBuilder {
	b.id = id
	return b
}

func (b *QuestionnaireRequestBuilder) Build() (*QuestionnaireRequest, error) {
	return b.QuestionnaireRequest, nil
}

func (req *QuestionnaireRequest) Endpoint() string {
	if req.id == 0 {
		return fmt.Sprintf("/t/%s/api/%s/questionnaires", req.tenant, APIVERSION)
	}
	return fmt.Sprintf("/t/%s/api/%s/questionnaires/%d", req.tenant, APIVERSION, req.id)
}
