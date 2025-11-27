package authentication

const USER_MANAGEMENT_APP string = "user-management"

/////// Currently OUT OF USE
// // encore:middleware target=tag:user-management
// func ValidationMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
// 	params := req.Data().PathParams
// 	piidStr := params.Get("piid")
// 	url_piid, err := uuid.FromString(piidStr)
// 	if err != nil {
// 		return middleware.Response{Err: errors.ErrorUnauthenticated}
// 	}

// 	data, ok := auth.Data().(entity.AuthData)
// 	if !ok {
// 		return next(req)
// 	}
// 	for _, i := range data.Instances {
// 		if i.PIID == url_piid {
// 			if !i.AppMapping[USER_MANAGEMENT_APP] {
// 				return middleware.Response{Err: errors.ErrorUnauthenticated}
// 			}
// 			break
// 		}
// 	}
// 	return next(req)
// }
