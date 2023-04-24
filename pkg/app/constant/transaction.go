package appconstant

import (
	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/external/utils/file"
)

// TransactionProcessIcon ...
type TransactionProcessIcon struct {
	Icon  *file.FilePhoto `json:"icon"`
	Color string          `json:"color"`
}

// transactionIconPending ...
var transactionIconPending = file.FilePhoto{
	ID:   "6260d931b3bd4aaf29073b8a",
	Name: "icon_affiliate_transaction_pending.png",
	Dimensions: &file.FileDimensions{
		Small: &file.FileSize{
			Width:  24,
			Height: 24,
		},
		Medium: &file.FileSize{
			Width:  48,
			Height: 48,
		},
	},
}

// transactionIconApproved ...
var transactionIconApproved = file.FilePhoto{
	ID:   "6260d931b3bd4aaf29073b8a",
	Name: "icon_affiliate_transaction_approved.png",
	Dimensions: &file.FileDimensions{
		Small: &file.FileSize{
			Width:  22,
			Height: 14,
		},
		Medium: &file.FileSize{
			Width:  44,
			Height: 28,
		},
	},
}

// transactionIconCashback ...
var transactionIconCashback = file.FilePhoto{
	ID:   "6260d931b3bd4aaf29073b8a",
	Name: "icon_affiliate_transaction_cashback.png",
	Dimensions: &file.FileDimensions{
		Small: &file.FileSize{
			Width:  22,
			Height: 14,
		},
		Medium: &file.FileSize{
			Width:  44,
			Height: 28,
		},
	},
}

// transactionIconRejected ...
var transactionIconRejected = file.FilePhoto{
	ID:   "6260d931b3bd4aaf29073b8a",
	Name: "icon_affiliate_transaction_rejected.png",
	Dimensions: &file.FileDimensions{
		Small: &file.FileSize{
			Width:  22,
			Height: 14,
		},
		Medium: &file.FileSize{
			Width:  44,
			Height: 28,
		},
	},
}

// TransactionProcessIcons ...
var TransactionProcessIcons = struct {
	Pending  []TransactionProcessIcon
	Approved []TransactionProcessIcon
	Cashback []TransactionProcessIcon
	Rejected []TransactionProcessIcon
}{
	Pending: []TransactionProcessIcon{
		{
			Icon:  transactionIconPending.GetResponseData(),
			Color: constant.ColorYellow,
		},
		{
			Icon:  transactionIconApproved.GetResponseData(),
			Color: constant.ColorGray,
		},
		{
			Icon:  transactionIconCashback.GetResponseData(),
			Color: constant.ColorGray,
		},
	},
	Approved: []TransactionProcessIcon{
		{
			Icon:  transactionIconPending.GetResponseData(),
			Color: constant.ColorBlue,
		},
		{
			Icon:  transactionIconApproved.GetResponseData(),
			Color: constant.ColorBlue,
		},
		{
			Icon:  transactionIconCashback.GetResponseData(),
			Color: constant.ColorGray,
		},
	},
	Cashback: []TransactionProcessIcon{
		{
			Icon:  transactionIconPending.GetResponseData(),
			Color: constant.ColorCyan,
		},
		{
			Icon:  transactionIconApproved.GetResponseData(),
			Color: constant.ColorCyan,
		},
		{
			Icon:  transactionIconCashback.GetResponseData(),
			Color: constant.ColorCyan,
		},
	},
	Rejected: []TransactionProcessIcon{
		{
			Icon:  transactionIconPending.GetResponseData(),
			Color: constant.ColorRedLight,
		},
		{
			Icon:  transactionIconRejected.GetResponseData(),
			Color: constant.ColorRedLight,
		},
	},
}
