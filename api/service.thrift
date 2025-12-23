namespace go api

include "product.thrift"

enum OrderStatus{
    UNPAID = 0      //未支付
    PAID = 1        //已支付
    SHIPPED = 2     //已发货
    COMPLETED = 3   //已完成
    CANCELLED = 4   //已取消
}
struct OrderItem{
    1:i64 productId
    2:string productName
    3:i32 quantity
    4:double price
}

struct Order{
    1:i64 orderId
    2:string orderNo
    3:i64 userId
    4:double totalAmount
    5:OrderStatus status
    6:list<OrderItem> items
    7:i64 createdAt
    8:string receiver
    9:string phone
    10:string address
}

struct CreateOrderReq{
    1:i64 userId
    2:list<OrderItem> items
    3:string receiver
    4:string phone
    5:string address
    6:string paymentMethod
}

struct CreateOrderResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:string orderNo
    5:double totalAmount
    6:string paymentUrl
}

struct PayOrderReq{
    1:string orderNo
    2:i64 userId
}

struct PayOrderResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:OrderStatus newStatus
}

struct CancelOrderReq{
    1:string orderNo
    2:i64 userId
    3:optional string reason
}

struct CancelOrderResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:OrderStatus newStatus
}

struct GetOrderDetailReq{
    1:string orderNo
    2:i64 userId
}

struct GetOrderDetailResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:Order order
}

struct QueryOrderReq{
    1:optional string orderNo
    2:optional i64 userId
    3:optional OrderStatus status
    4:i32 page = 1
    5:i32 pageSize = 10
}

struct QueryOrderResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:i32 total
    5:i32 page
    6:i32 pageSize
    7:list<Order> orders
}

struct UpdateOrderStatusReq{
    1:string orderNo
    2:OrderStatus status
    3:i64 userId
}

struct UpdateOrderStatusResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:OrderStatus oldStatus
    5:OrderStatus newStatus
}

service OrderService{
    CreateOrderResp CreateOrder(1:CreateOrderReq req)
    PayOrderResp PayOrder(1:PayOrderReq req)
    CancelOrderResp CancelOrder(1:CancelOrderReq req)
    GetOrderDetailResp GetOrderDetail(1:GetOrderDetailReq req)
    QueryOrderResp QueryOrders(1:QueryOrderReq req)
    UpdateOrderStatusResp UpdateOrderStatus(1:UpdateOrderStatusReq req)
}