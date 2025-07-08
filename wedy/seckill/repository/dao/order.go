package dao

type Order struct {
	OrderID        uint64  `gorm:"primaryKey;autoIncrement;comment:订单ID"`
	UserID         uint64  `gorm:"not null;comment:用户ID"`
	OriginalAmount float64 `gorm:"type:decimal(10,2);not null;comment:订单原价"`
	DiscountAmount float64 `gorm:"type:decimal(10,2);default:0;comment:优惠金额"`
	ShippingAmount float64 `gorm:"type:decimal(10,2);default:0;comment:运费"`
	PayableAmount  float64 `gorm:"type:decimal(10,2);not null;comment:应付金额"`
	OrderStatus    string  `gorm:"type:varchar(20);default:'pending';comment:订单状态(pending/paid/shipped/delivered/closed/canceled)"`
}

type Promote struct {
	OrderID        uint64  `gorm:"primaryKey;autoIncrement;comment:订单ID"`
	CouponID       uint64  `gorm:"comment:优惠券ID"`
	CouponCode     string  `gorm:"type:varchar(32);comment:优惠券码"`
	CouponDiscount float64 `gorm:"type:decimal(10,2);default:0;comment:优惠券优惠金额"`
}
