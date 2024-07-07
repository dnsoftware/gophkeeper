// Серверные перехватчики
package handlers

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	"github.com/dnsoftware/gophkeeper/internal/utils"
)

// checkUserInterceptor проверка авторизованности пользователя
func checkUserInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if info.FullMethod == constants.ExcludeMethodRegistration || info.FullMethod == constants.ExcludeMethodPing || info.FullMethod == constants.ExcludeMethodLogin {
		return handler(ctx, req)
	}

	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		tok := headers.Get(constants.TokenKey)
		if len(tok) <= 0 || tok[0] == "" || tok == nil {
			return nil, status.Errorf(codes.PermissionDenied, `Unauthorized`)
		}

		userID := utils.GetUserID(tok[0])
		if userID <= 0 {
			return nil, status.Errorf(codes.PermissionDenied, `Unauthorized`)
		}

	}

	return handler(ctx, req)
}
