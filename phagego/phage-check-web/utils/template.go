package utils

func GetDialStatusNameMap() map[int64]string {
	m := map[int64]string{
		0:  "未拨打",
		1:  "有效接听",
		2:  "不是本人",
		3:  "接通挂断",
		4:  "无法接通",
		5:  "无人接听",
		6:  "拒接",
		7:  "关机",
		8:  "停机",
		9:  "空号",
		10: "其他",
	}
	return m
}
