package handler

import (
	"context"
	"strconv"

	"ecommerce/gateway/internal/client"
	"ecommerce/gateway/pkg/response"
	"ecommerce/product-service/kitex_gen/api"

	"github.com/cloudwego/hertz/pkg/app"
)

// 添加缺失的辅助函数
func safeInt64(i int64) int64 {
	return i
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// CreateProduct 创建商品（管理员）
func CreateProduct(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		var req api.CreateProductReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		resp, err := clientManager.ProductClient.CreateProduct(c, &req)
		if err != nil {
			response.Error(ctx, 500, "创建商品失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeString(resp.Message))
			return
		}
		var productID int64
		if resp.Product != nil {
			productID = resp.Product.Id
		}

		response.Success(ctx, map[string]interface{}{
			"product_id": productID,
			"message":    "商品创建成功",
		})
	}
}

// GetProduct 获取商品详情
func GetProduct(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "商品ID格式错误")
			return
		}

		req := &api.GetProductReq{
			Id: productID,
		}

		resp, err := clientManager.ProductClient.GetProduct(c, req)
		if err != nil {
			response.Error(ctx, 500, "获取商品信息失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeString(resp.Message))
			return
		}

		response.Success(ctx, resp.Product)
	}
}

// UpdateProduct 更新商品（管理员）
func UpdateProduct(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "商品ID格式错误")
			return
		}

		var req api.UpdateProductReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		req.Id = productID

		resp, err := clientManager.ProductClient.UpdateProduct(c, &req)
		if err != nil {
			response.Error(ctx, 500, "更新商品失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeString(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "商品更新成功",
		})
	}
}

// ListProducts 获取商品列表（用户搜索）
func ListProducts(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		page, _ := strconv.Atoi(ctx.Query("page"))
		pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
		keyword := ctx.Query("keyword")
		category := ctx.Query("category")

		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		req := &api.UserSearchProductsReq{
			Page:     int32(page),
			PageSize: int32(pageSize),
		}

		if keyword != "" {
			req.Keyword = &keyword
		}
		if category != "" {
			req.Category = &category
		}

		resp, err := clientManager.ProductClient.UserSearchProducts(c, req)
		if err != nil {
			response.Error(ctx, 500, "获取商品列表失败: "+err.Error())
			return
		}

		response.SuccessWithPagination(ctx, resp.Products, int64(resp.Total), page, pageSize)
	}
}

// GetProductsByCategory 按分类获取商品
func GetProductsByCategory(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		category := ctx.Param("category")
		if category == "" {
			response.Error(ctx, 400, "分类不能为空")
			return
		}

		page, _ := strconv.Atoi(ctx.Query("page"))
		pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}

		req := &api.UserSearchProductsReq{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Category: &category,
		}

		resp, err := clientManager.ProductClient.UserSearchProducts(c, req)
		if err != nil {
			response.Error(ctx, 500, "获取商品列表失败: "+err.Error())
			return
		}

		response.SuccessWithPagination(ctx, resp.Products, int64(resp.Total), page, pageSize)
	}
}

// SearchProducts 搜索商品（用户）
func SearchProducts(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		var req api.UserSearchProductsReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		resp, err := clientManager.ProductClient.UserSearchProducts(c, &req)
		if err != nil {
			response.Error(ctx, 500, "搜索商品失败: "+err.Error())
			return
		}

		response.SuccessWithPagination(ctx, resp.Products, int64(resp.Total), int(req.Page), int(req.PageSize))
	}
}

// AdminSearchProducts 管理员搜索商品
func AdminSearchProducts(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		var req api.AdminSearchProductsReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		resp, err := clientManager.ProductClient.AdminSearchProducts(c, &req)
		if err != nil {
			response.Error(ctx, 500, "搜索商品失败: "+err.Error())
			return
		}

		response.SuccessWithPagination(ctx, resp.Products, int64(resp.Total), int(req.Page), int(req.PageSize))
	}
}

// OnlineProduct 上架商品（管理员）
func OnlineProduct(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "商品ID格式错误")
			return
		}

		req := &api.OnlineProductReq{
			Id: productID,
		}

		resp, err := clientManager.ProductClient.OnlineProduct(c, req)
		if err != nil {
			response.Error(ctx, 500, "上架商品失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeString(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "商品上架成功",
		})
	}
}

// OfflineProduct 下架商品（管理员）
func OfflineProduct(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "商品ID格式错误")
			return
		}

		req := &api.OfflineProductReq{
			Id: productID,
		}

		resp, err := clientManager.ProductClient.OfflineProduct(c, req)
		if err != nil {
			response.Error(ctx, 500, "下架商品失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeString(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "商品下架成功",
		})
	}
}

// DeleteProduct 删除商品（管理员）
func DeleteProduct(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		productID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "商品ID格式错误")
			return
		}

		req := &api.DeleteProductReq{
			Id: productID,
		}

		resp, err := clientManager.ProductClient.DeleteProduct(c, req)
		if err != nil {
			response.Error(ctx, 500, "删除商品失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeString(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "商品删除成功",
		})
	}
}
