package statistics

import (
	"log"
	"maps"
	"reflect"
	"slices"
	"tool/app/internal/models"
)

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

	log.Print(reflect.TypeOf(permissions).Kind())
	switch permissions := permissions.(type) {
	case string:
		if slices.Contains(slices.Collect(maps.Keys(permissionsCoarse)), permissions) {
			permissionsCoarse[permissions] += 1
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
