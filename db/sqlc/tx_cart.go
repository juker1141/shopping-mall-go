package db

// type UpdateCartTxParams struct {
// 	Name         string  `json:"name"`
// 	categoriesID []int64 `json:"categories_id"`
// }

// type ProductTxResult struct {
// 	Role           Role         `json:"role"`
// 	PermissionList []Permission `json:"permission_list"`
// }

// func (store *SQLStore) UpdateCartTx(ctx context.Context, arg CreateProductTxParams) (ProductTxResult, error) {
// 	var result ProductTxResult

// 	err := store.execTx(ctx, func(q *Queries) error {
// 		var err error
// 		var permissionList []Permission

// 		if len(arg.PermissionsID) <= 0 {
// 			err = fmt.Errorf("at least one permission is required")
// 			return err
// 		}

// 		result.Role, err = q.CreateRole(ctx, arg.Name)
// 		if err != nil {
// 			return err
// 		}

// 		for _, permissionId := range arg.PermissionsID {
// 			_, err := q.CreateRolePermission(ctx, CreateRolePermissionParams{
// 				RoleID: pgtype.Int4{
// 					Int32: int32(result.Role.ID),
// 					Valid: true,
// 				},
// 				PermissionID: pgtype.Int4{
// 					Int32: int32(permissionId),
// 					Valid: true,
// 				},
// 			})
// 			if err != nil {
// 				return err
// 			}

// 			permission, err := q.GetPermission(ctx, permissionId)
// 			if err != nil {
// 				return err
// 			}
// 			permissionList = append(permissionList, permission)
// 		}

// 		result.PermissionList = permissionList

// 		return nil
// 	})

// 	return result, err
// }
