package tron_lib

type AccountInfoData struct {
	Address               string                  `json:"address"`
	Balance               int64                  `json:"balance"`
	CreateTime            int64                   `json:"create_time"`
	LatestOprationTime    int64                   `json:"latest_opration_time"`
	LatestConsumeFreeTime int64                   `json:"latest_consume_free_time"`
	AccountResource       *AccountResourceData    `json:"account_resource"`
	OwnerPermission       *OwnerPermissionData    `json:"owner_permission"`
	ActivePermission      []*ActivePermissionData `json:"active_permission"`
}

type AccountResourceData struct {
	LatestConsumeTimeForEnergy int64 `json:"latest_consume_time_for_energy"`
}

type OwnerPermissionData struct {
	PermissionName string `json:"permission_name"`
	Threshold      int64  `json:"threshold"`
	Keys           []*struct {
		Address string `json:"address"`
		Weight  int64  `json:"weight"`
	} `json:"keys"`
}

type ActivePermissionData struct {
	Type           string `json:"type"`
	Id             int64  `json:"id"`
	PermissionName string `json:"permission_name"`
	Threshold      int64  `json:"threshold"`
	Operations     string `json:"operations"`
	Keys           []*struct {
		Address string `json:"address"`
		Weight  int64  `json:"weight"`
	} `json:"keys"`
}
