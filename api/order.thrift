namespace go api

include "product.thrift"
include "user.thrift"

enum OrderStatus {
    PENDING = 0     // 待支付
    PAID = 1        // 已支付  
    SHIPPED = 2     // 已发货
    COMPLETED = 3   // 已完成
    CANCELLED = 4   // 已取消
    REFUNDED = 5    // 已退款
}

enum RefundStatus {
    PENDING = 0     // 待处理
    APPROVED = 1    // 已同意
    REJECTED = 2    // 已拒绝
    PROCESSING = 3  // 处理中
    COMPLETED = 4   // 已完成
}

enum TimeoutType {
    ORDER_UNPAID = 0      // 订单未支付
    STOCK_RESERVATION = 1 // 库存预占
}

struct OrderItem {
    1:i64 productId
    2:string productName
    3:i32 quantity
    4:double price
    5:optional string productImage  // 商品图片
}

struct Order {
    1:i64 id
    2:string orderNo
    3:i64 userId
    4:double totalAmount
    5:OrderStatus status
    6:list<OrderItem> items
    7:string address
    8:string phone
    9:i64 createdAt
    10:i64 updatedAt
    11:optional string receiver     // 收货人姓名
    12:optional string paymentNo    // 支付单号
    13:optional string shippingNo   // 物流单号
}

struct StockReservation {
    1:string reserveId
    2:string orderNo
    3:i64 productId
    4:i32 quantity
    5:string status
    6:i64 expireTime
    7:i64 createdAt
    8:i64 updatedAt
}

struct RefundOrder {
    1:string refundNo
    2:string orderNo
    3:i64 userId
    4:double amount
    5:RefundStatus status
    6:string reason
    7:i64 createdAt
    8:i64 updatedAt
    9:optional string processor    // 处理人
    10:optional i64 processedAt    // 处理时间
}

struct TimeoutTask {
    1:string taskId
    2:string orderNo
    3:TimeoutType type
    4:string status
    5:i64 expireTime
    6:i32 retryCount
    7:i64 createdAt
    8:i64 updatedAt
}

// 创建订单
struct CreateOrderReq {
    1:i64 userId
    2:list<OrderItem> items
    3:string address
    4:string phone
    5:optional string receiver      // 收货人
    6:optional string paymentMethod // 支付方式
}

struct CreateOrderResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:string orderNo
    5:double totalAmount
    6:optional string paymentUrl    // 支付链接
}

// 获取订单详情
struct GetOrderReq {
    1:string orderNo
    2:optional i64 userId          // 用于验证订单所属用户
}

struct GetOrderResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:Order order
}

// 查询订单列表
struct ListOrdersReq {
    1:i64 userId
    2:optional OrderStatus status
    3:i32 page = 1
    4:i32 pageSize = 10
    5:optional i64 startTime       // 开始时间
    6:optional i64 endTime         // 结束时间
}

struct ListOrdersResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:i32 total
    5:i32 page
    6:i32 pageSize
    7:list<Order> orders
}

// 支付订单
struct PayOrderReq {
    1:string orderNo
    2:i64 userId
    3:optional string paymentNo    // 支付单号
}

struct PayOrderResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:OrderStatus newStatus
    5:optional i64 paidAt          // 支付时间
}

// 取消订单
struct CancelOrderReq {
    1:string orderNo
    2:i64 userId
    3:string reason
}

struct CancelOrderResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:OrderStatus newStatus
    5:optional i64 cancelledAt     // 取消时间
}

// 退款申请
struct ApplyRefundReq {
    1:string orderNo
    2:i64 userId
    3:string reason
    4:optional double amount       // 退款金额（部分退款）
}

struct ApplyRefundResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:string refundNo
    5:RefundStatus status
}

// 处理退款
struct ProcessRefundReq {
    1:string refundNo
    2:i64 processorId              // 处理人ID
    3:RefundStatus action          // 处理动作：APPROVED/REJECTED
    4:optional string remark       // 处理备注
}

struct ProcessRefundResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:RefundStatus newStatus
}

// 库存预占
struct ReserveStockReq {
    1:string orderNo
    2:i64 productId
    3:i32 quantity
    4:i64 expireSeconds = 900      // 默认15分钟
}

struct ReserveStockResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:string reserveId
}

// 释放库存
struct ReleaseStockReq {
    1:string reserveId
    2:string reason
}

struct ReleaseStockResp {
    1:bool success
    2:i32 code = 0
    3:string message
}

// 确认库存（扣减）
struct ConfirmStockReq {
    1:string reserveId
    2:string orderNo
}

struct ConfirmStockResp {
    1:bool success
    2:i32 code = 0
    3:string message
}

// 超时处理
struct ProcessTimeoutReq {
    1:string taskId
    2:TimeoutType type
}

struct ProcessTimeoutResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:map<string, string> results
}

// 订单统计
struct OrderStatsReq {
    1:i64 userId
    2:optional i64 startTime
    3:optional i64 endTime
}

struct OrderStatsResp {
    1:bool success
    2:i32 code = 0
    3:string message
    4:i32 totalOrders
    5:double totalAmount
    6:map<string, i32> statusCounts  // 各状态订单数
}

service OrderService {
    // 订单生命周期
    CreateOrderResp CreateOrder(1:CreateOrderReq req)
    PayOrderResp PayOrder(1:PayOrderReq req)
    CancelOrderResp CancelOrder(1:CancelOrderReq req)
    
    // 订单查询
    GetOrderResp GetOrder(1:GetOrderReq req)
    ListOrdersResp ListOrders(1:ListOrdersReq req)
    
    // 退款管理
    ApplyRefundResp ApplyRefund(1:ApplyRefundReq req)
    ProcessRefundResp ProcessRefund(1:ProcessRefundReq req)
    
    // 库存管理（与库存服务交互）
    ReserveStockResp ReserveStock(1:ReserveStockReq req)
    ReleaseStockResp ReleaseStock(1:ReleaseStockReq req)
    ConfirmStockResp ConfirmStock(1:ConfirmStockReq req)
    
    // 超时处理（定时任务）
    ProcessTimeoutResp ProcessTimeout(1:ProcessTimeoutReq req)
    
    // 统计
    OrderStatsResp GetOrderStats(1:OrderStatsReq req)
    
    // 订单状态更新（内部/管理用）
    CancelOrderResp UpdateOrderStatus(1:CancelOrderReq req)
    
    // 发货
    PayOrderResp ShipOrder(1:PayOrderReq req)  // 复用PayOrderReq结构
    
    // 确认收货
    PayOrderResp ConfirmReceipt(1:PayOrderReq req)
}