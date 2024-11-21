package statistics

import (
	"maps"
	"slices"
	"tool/app/internal/helpers"
	"tool/app/internal/models"
)

func permissionsArrayCount(permissions []models.PermissionsStatistics) models.PermissionsStatistics {
	var valsRead []int
	var valsWrite []int
	var valsNone []int
	var valsReadAll []int
	var valsWriteAll []int
	var valsNoneAll []int
	var valsDefault []int
	var counts []int

	for _, permission := range permissions {
		valsRead = append(valsRead, permission.FineGrained.Read.Total)
		valsWrite = append(valsWrite, permission.FineGrained.Write.Total)
		valsNone = append(valsNone, permission.FineGrained.None.Total)
		valsReadAll = append(valsReadAll, permission.CoarseGrained.ReadAll.Total)
		valsWriteAll = append(valsWriteAll, permission.CoarseGrained.WriteAll.Total)
		valsNoneAll = append(valsNoneAll, permission.CoarseGrained.NoneAll.Total)
		valsDefault = append(valsDefault, permission.CoarseGrained.Default.Total)
		counts = append(counts, helpers.ComputeSum([]int{
			permission.FineGrained.Read.Total, permission.FineGrained.Write.Total, permission.FineGrained.None.Total,
			permission.CoarseGrained.ReadAll.Total, permission.CoarseGrained.WriteAll.Total,
			permission.CoarseGrained.NoneAll.Total, permission.CoarseGrained.Default.Total,
		}))
	}

	return models.PermissionsStatistics{
		FineGrained: models.FineGrainedStatistics{
			Read:  BuildIntStatistics(valsRead),
			Write: BuildIntStatistics(valsWrite),
			None:  BuildIntStatistics(valsNone),
		},
		CoarseGrained: models.CoarseGrainedStatistics{
			ReadAll:  BuildIntStatistics(valsReadAll),
			WriteAll: BuildIntStatistics(valsWriteAll),
			NoneAll:  BuildIntStatistics(valsNoneAll),
			Default:  BuildIntStatistics(valsDefault),
		},
		Count: BuildIntStatistics(counts),
	}
}

func permissionsCount(permissions interface{}) models.PermissionsStatistics {
	permissionsFine := map[string]int{
		"read":  0,
		"write": 0,
		"none":  0,
	}

	permissionsCoarse := map[string]int{
		"read-all":  0,
		"write-all": 0,
		"none-all":  0,
		"default":   0,
	}

	switch permissions := permissions.(type) {
	case string:
		if slices.Contains(slices.Collect(maps.Keys(permissionsCoarse)), permissions) {
			permissionsCoarse[permissions] += 1
		} else {
			permissionsCoarse["default"] += 1
		}
	case map[string]interface{}:
		if len(permissions) == 0 {
			permissionsCoarse["none-all"] += 1
		} else {

			for _, access := range permissions {
				permissionsFine[access.(string)] += 1
			}
		}
	default:
		permissionsCoarse["default"] += 1
	}

	count := 0
	for key := range permissionsFine {
		count += permissionsFine[key]
	}

	for key := range permissionsCoarse {
		count += permissionsCoarse[key]
	}

	return models.PermissionsStatistics{
		FineGrained: models.FineGrainedStatistics{
			Read:  models.IntStatistics{Total: permissionsFine["read"]},
			Write: models.IntStatistics{Total: permissionsFine["write"]},
			None:  models.IntStatistics{Total: permissionsFine["none"]},
		},
		CoarseGrained: models.CoarseGrainedStatistics{
			ReadAll:  models.IntStatistics{Total: permissionsCoarse["read-all"]},
			WriteAll: models.IntStatistics{Total: permissionsCoarse["write-all"]},
			NoneAll:  models.IntStatistics{Total: permissionsCoarse["none-all"]},
			Default:  models.IntStatistics{Total: permissionsCoarse["default"]},
		},
		Count: models.IntStatistics{Total: count},
	}
}
