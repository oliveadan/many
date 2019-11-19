package utils

func GetAnimal(num int) string{
	var name string
	switch num {
	case 1:
		name = "鼠"
		break
	case 2:
		name = "牛"
		break
	case 3:
		name = "虎"
		break
	case 4:
		name = "兔"
		break
	case 5:
		name = "龙"
		break
	case 6:
		name = "蛇"
		break
	case 7:
		name = "马"
		break
	case 8:
		name = "羊"
		break
	case 9:
		name = "猴"
		break
	case 10:
		name = "鸡"
		break
	case 11:
		name = "狗"
		break
	case 12:
		name = "猪"
		break
	}
	return name
}
