package handler

import (
	"context"
	"strconv"

	"ecommerce/product-service/internal/service"
	api "ecommerce/product-service/kitex_gen/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type ProductHTTPHandler struct {
	productService service.ProductService
}

func NewProductHTTPHandler(productService service.ProductService) *ProductHTTPHandler {
	return &ProductHTTPHandler{
		productService: productService,
	}
}

// HTTP响应包装器
type HTTPResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
}

func stringPtr(s string) *string {
	return &s
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// 创建产品
func (h *ProductHTTPHandler) CreateProduct(ctx context.Context, c *app.RequestContext) {
	var req api.CreateProductReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.productService.CreateProduct(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "创建产品失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "创建产品失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data:    resp.Product,
	})
}

// 获取产品详情
func (h *ProductHTTPHandler) GetProduct(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.productService.GetProduct(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取产品失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "获取产品失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		status := consts.StatusNotFound
		if resp.Code == 500 {
			status = consts.StatusInternalServerError
		}
		c.JSON(status, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data:    resp.Product,
	})
}

// 更新产品
func (h *ProductHTTPHandler) UpdateProduct(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}

	var req api.UpdateProductReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	req.Id = id
	resp, err := h.productService.UpdateProduct(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "更新产品失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "更新产品失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data:    resp.Product,
	})
}

// 删除产品
func (h *ProductHTTPHandler) DeleteProduct(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.productService.DeleteProduct(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "删除产品失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "删除产品失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
	})
}

// 搜索产品
func (h *ProductHTTPHandler) SearchProducts(ctx context.Context, c *app.RequestContext) {
	var req api.UserSearchProductsReq
	if category := c.Query("category"); category != "" {
		req.Category = &category
	}
	if keyword := c.Query("keyword"); keyword != "" {
		req.Keyword = &keyword
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			req.MinPrice = &min
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			req.MaxPrice = &max
		}
	}
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	req.Page = int32(page)

	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	req.PageSize = int32(pageSize)
	resp, err := h.productService.UserSearchProducts(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "搜索产品失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "搜索产品失败",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: utils.H{
			"total":     resp.Total,
			"page":      resp.Page,
			"page_size": resp.PageSize,
			"products":  resp.Products,
		},
	})
}

// 管理员搜索产品
func (h *ProductHTTPHandler) AdminSearchProducts(ctx context.Context, c *app.RequestContext) {
	var req api.AdminSearchProductsReq
	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			req.Id = &id
		}
	}
	if category := c.Query("category"); category != "" {
		req.Category = &category
	}
	if keyword := c.Query("keyword"); keyword != "" {
		req.Keyword = &keyword
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			req.MinPrice = &min
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			req.MaxPrice = &max
		}
	}
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	req.Page = int32(page)
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	req.PageSize = int32(pageSize)
	resp, err := h.productService.AdminSearchProducts(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "管理员搜索产品失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "搜索产品失败",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: utils.H{
			"total":     resp.Total,
			"page":      resp.Page,
			"page_size": resp.PageSize,
			"products":  resp.Products,
		},
	})
}

// 上架产品
func (h *ProductHTTPHandler) OnlineProduct(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}
	resp, err := h.productService.OnlineProduct(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "上架产品失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "上架产品失败",
			Error:   err.Error(),
		})
		return
	}
	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}
	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: utils.H{
			"old_status":  resp.OldStatus,
			"new_status":  resp.NewStatus_,
			"operated_at": resp.OperatedAt,
		},
	})
}

// 下架产品
func (h *ProductHTTPHandler) OfflineProduct(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}
	resp, err := h.productService.OfflineProduct(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "下架产品失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "下架产品失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}
	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: utils.H{
			"old_status":  resp.OldStatus,
			"new_status":  resp.NewStatus_,
			"operated_at": resp.OperatedAt,
		},
	})
}
