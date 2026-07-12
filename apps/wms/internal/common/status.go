// Package common provides shared helpers for admin route modules.
package common

// ParcelStatusCN returns a user-friendly Chinese display name for a parcel status string.
func ParcelStatusCN(status string) string {
	switch status {
	case "pre_declared":
		return "预报"
	case "received":
		return "已入仓"
	case "weighed":
		return "已称重"
	case "stored":
		return "已上架"
	case "picked":
		return "已拣货"
	case "packed":
		return "已打包"
	case "outbound":
		return "已出货"
	case "container_area":
		return "待装柜"
	case "loaded":
		return "已装柜"
	case "shipped":
		return "运输中"
	case "customs":
		return "清关中"
	case "delivering":
		return "配送中"
	case "delivered":
		return "已签收"
	case "abnormal":
		return "异常"
	case "returned":
		return "已退货"
	default:
		return status
	}
}

// OrderStatusCN returns a user-friendly Chinese display name for an order status string.
func OrderStatusCN(status string) string {
	switch status {
	case "pending_picking":
		return "待拣货"
	case "picking":
		return "拣货中"
	case "pending_packing":
		return "待打包"
	case "pending_loading":
		return "待装柜"
	case "loaded":
		return "已装柜"
	case "in_transit":
		return "运输中"
	case "customs_clearance":
		return "清关中"
	case "out_for_delivery":
		return "派送中"
	case "completed":
		return "已完成"
	case "cancelled":
		return "已取消"
	case "shipped":
		return "已发货"
	default:
		return status
	}
}

// CargoTypeCN returns a user-friendly Chinese display name for a cargo type string.
func CargoTypeCN(cargoType string) string {
	switch cargoType {
	case "general":
		return "普货"
	case "sensitive":
		return "特货"
	case "dangerous":
		return "危险品"
	default:
		return cargoType
	}
}

// TransportTypeCN returns a user-friendly Chinese display name for a transport type string.
func TransportTypeCN(transportType string) string {
	switch transportType {
	case "air":
		return "空运"
	case "sea":
		return "海运"
	case "sea_express":
		return "海快"
	default:
		return transportType
	}
}

// ActionBadgeClass returns a badge CSS class for an audit action type.
func ActionBadgeClass(action string) string {
	switch action {
	case "CREATE":
		return "success"
	case "UPDATE":
		return "brand"
	case "DELETE":
		return "danger"
	case "LOGIN":
		return "info"
	case "EXPORT":
		return "warning"
	default:
		return "secondary"
	}
}
