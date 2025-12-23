namespace go api

enum ProductStatus{
    DRAFT = 0   //预上架
    ONLINE = 1  //上架
    OFFLINE = 2 //下架
    DELETED = 3 //删除
}

struct Product{
    1:i64 id
    2:string name
    3:string avatar
    4:string category
    5:double price
    6:i32 stock
    7:ProductStatus status = ProductStatus.DRAFT
    8:i64 createdAt
    9:i64 updatedAt
    10:optional string brand
}

struct SimpleProduct{
    1:i64 id
    2:string category
    3:double price
    4:i32 stock
    5:ProductStatus status = ProductStatus.ONLINE
    6:optional string brand
    7:string name
    8:string avatar
}

struct CreateProductReq{
    1:string name
    2:string avatar
    3:string category
    4:double price
    5:i32 stock
    6:optional string brand
    7:optional ProductStatus status = ProductStatus.DRAFT
}

struct CreateProductResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:Product product
}

struct OnlineProductReq{
    1:i64 id
}

struct OnlineProductResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:ProductStatus oldStatus
    5:ProductStatus newStatus
    6:i64 operatedAt
}

struct OfflineProductReq{
    1:i64 id
}

struct OfflineProductResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:ProductStatus oldStatus
    5:ProductStatus newStatus
    6:i64 operatedAt
}

struct UserSearchProductsReq{
    1:optional string category
    2:optional double minPrice
    3:optional double maxPrice
    4:optional string keyword
    5:i32 page = 1
    6:i32 pageSize = 20
}

struct UserSearchProductsResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:i32 total
    5:i32 page
    6:i32 pageSize
    7:list<SimpleProduct> products
}

struct AdminSearchProductsReq{
    1:optional i64 id
    2:optional string category
    3:optional double minPrice
    4:optional double maxPrice
    5:optional string keyword
    6:i32 page = 1
    7:i32 pageSize = 20
}

struct AdminSearchProductsResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:i32 total
    5:i32 page
    6:i32 pageSize
    7:list<Product> products
}

service ProductService{
    CreateProductResp CreateProduct(1:CreateProductReq req)
    GetProductResp GetProduct(1:GetProductReq req)
    UpdateProductResp UpdateProduct(1:UpdateProductReq req)
    OnlineProductResp OnlineProduct(1:OnlineProductReq req)
    OfflineProductResp OfflineProduct(1:OfflineProductReq req)
    DeleteProductResp DeleteProduct(1:DeleteProductReq req)
    UserSearchProductsResp UserSearchProducts(1:UserSearchProductsReq req)
    AdminSearchProductsResp AdminSearchProducts(1:AdminSearchProductsReq req)
}