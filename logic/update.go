package logic

import "github.com/qianguozheng/goadmin/model"

type UpgradeLogic struct{}

var DefaultUpgrade = UpgradeLogic{}

//Return Unique version of models
func (upgrade *UpgradeLogic) GetAllVersions() []string {
	var upg []model.Upgrade
	model.DB.Find(&upg)

	var a []string
	for _, v := range upg {
		a = append(a, v.Version)
	}
	return RemoveDuplicateAndEmpty(a)
}

func RemoveDuplicateAndEmpty(a []string) (ret []string) {
	alen := len(a)
	for i := 0; i < alen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
