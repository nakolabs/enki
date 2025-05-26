package handler

import (
	"enuma-elish/internal/exam/service"
	"enuma-elish/internal/exam/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	service   service.Service
	validator *validator.Validate
}

func New(service service.Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		validator: validate,
	}
}

func (h *Handler) CreateExam(c *gin.Context) {
	data := request.CreateExamRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	// Additional validation for question types
	if err := data.Validate(); err != nil {
		c.Error(err)
		return
	}

	err = h.service.CreateExam(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusCreated).
		SetMessage("exam created successfully").
		SetData(data)

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetDetailExam(c *gin.Context) {
	examIDStr := c.Param("exam_id")
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, err := h.service.GetDetailExam(c.Request.Context(), examID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get detail exam error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get detail exam success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetListExams(c *gin.Context) {
	httpQuery := request.GetListExamQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()
	err := c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, meta, err := h.service.GetListExams(c.Request.Context(), httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get list exams error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get list exams success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateExam(c *gin.Context) {
	examIDStr := c.Param("exam_id")
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data := request.CreateExamRequest{}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	// Additional validation for question types
	if err := data.Validate(); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdateExam(c.Request.Context(), examID, data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("exam updated successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteExam(c *gin.Context) {
	examIDStr := c.Param("exam_id")
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	err = h.service.DeleteExam(c.Request.Context(), examID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("delete exam error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("exam deleted successfully")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) AssignExamToClass(c *gin.Context) {
	data := request.AssignExamToClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.AssignExamToClass(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("exam assigned to class successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GradeExam(c *gin.Context) {
	data := request.GradeExamRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.GradeExam(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("exam graded successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetExamStudents(c *gin.Context) {
	examIDStr := c.Param("exam_id")
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	httpQuery := request.GetExamStudentsQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()

	err = c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(httpQuery); err != nil {
		c.Error(err)
		return
	}

	data, meta, err := h.service.GetExamStudents(c.Request.Context(), examID, httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get exam students error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get exam students success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) SubmitExamAnswers(c *gin.Context) {
	studentIDStr := c.GetString("user_id") // Assuming user_id is set by auth middleware
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data := request.SubmitExamAnswersRequest{}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.SubmitExamAnswers(c.Request.Context(), studentID, data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("exam answers submitted successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetStudentExams(c *gin.Context) {
	studentIDStr := c.GetString("user_id") // Assuming user_id is set by auth middleware
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	httpQuery := request.GetStudentExamsQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()
	err = c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, meta, err := h.service.GetStudentExams(c.Request.Context(), studentID, httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get student exams error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get student exams success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetStudentExamDetail(c *gin.Context) {
	studentIDStr := c.GetString("user_id") // Assuming user_id is set by auth middleware
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	examIDStr := c.Param("exam_id")
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, err := h.service.GetStudentExamDetail(c.Request.Context(), examID, studentID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get student exam detail error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get student exam detail success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}
