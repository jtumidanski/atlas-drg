package inventory

import "strings"

const (
	TypeValueEquip int8 = 1
	TypeValueUse   int8 = 2
	TypeValueSetup int8 = 3
	TypeValueETC   int8 = 4
	TypeValueCash  int8 = 5
	TypeEquip           = "EQUIP"
	TypeUse             = "USE"
	TypeSetup           = "SETUP"
	TypeETC             = "ETC"
	TypeCash            = "CASH"

	ItemTypeEquip = "EQUIPMENT"
	ItemTypeItem  = "ITEM"
)

func GetInventoryType(itemId uint32) (int8, bool) {
	t := int8(itemId / 1000000)
	if t >= 1 && t <= 5 {
		return t, true
	}
	return 0, false
}

func GetTypeFromByte(opt int8) (string, bool) {
	if opt >= 1 && opt <= 5 {
		switch opt {
		case TypeValueEquip:
			return TypeEquip, true
		case TypeValueUse:
			return TypeUse, true
		case TypeValueSetup:
			return TypeSetup, true
		case TypeValueETC:
			return TypeETC, true
		case TypeValueCash:
			return TypeCash, true
		}
	}
	return "", false
}

func GetByteFromName(name string) (int8, bool) {
	if strings.EqualFold(name, TypeEquip) {
		return TypeValueEquip, true
	} else if strings.EqualFold(name, TypeUse) {
		return TypeValueUse, true
	} else if strings.EqualFold(name, TypeSetup) {
		return TypeValueSetup, true
	} else if strings.EqualFold(name, TypeETC) {
		return TypeValueETC, true
	} else if strings.EqualFold(name, TypeCash) {
		return TypeValueCash, true
	}
	return 0, false
}
