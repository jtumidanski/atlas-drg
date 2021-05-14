package drop

type DropOperator func(*Drop)
type DropsOperator func([]*Drop)

func ForEachDrop(f DropOperator) {
	ForAllDrops(ExecuteForEachDrop(f))
}

func ForAllDrops(f DropsOperator) {
	drops := GetRegistry().GetAllDrops()
	f(drops)
}

func ExecuteForEachDrop(f DropOperator) DropsOperator {
	return func(drops []*Drop) {
		for _, drop := range drops {
			f(drop)
		}
	}
}
