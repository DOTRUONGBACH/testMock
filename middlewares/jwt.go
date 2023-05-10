package middlewares

import (
	"context"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func JWTUnaryInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	// Lấy thông tin về token từ metadata của yêu cầu
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "Metadata not found")
	}
	tokenString := strings.Join(md["authorization"], "")
	if tokenString == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "Authorization token not found")
	}

	// Xác thực token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("mysecretkey"), nil
	})
	if err != nil || !token.Valid {
		return nil, grpc.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	// Lưu thông tin về token vào context của yêu cầu
	ctx = context.WithValue(ctx, "token", token)

	// Xử lý yêu cầu bằng handler tiếp theo trong chuỗi middleware
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func JWTStreamInterceptor(
	srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) error {
	// Lấy thông tin về token từ metadata của yêu cầu
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "Metadata not found")
	}
	tokenString := strings.Join(md["authorization"], "")
	if tokenString == "" {
		return grpc.Errorf(codes.Unauthenticated, "Authorization token not found")
	}

	// Xác thực token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("mysecretkey"), nil
	})
	if err != nil || !token.Valid {
		return grpc.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	// Lưu thông tin về token vào context của yêu cầu
	ctx = context.WithValue(ctx, "token", token)

	// Xử lý yêu cầu bằng handler tiếp theo trong chuỗi middleware
	err = handler(srv, &jwtServerStream{ServerStream: ss})
	if err != nil {
		return err
	}
	return nil
}

// Định nghĩa đối tượng ServerStream cho middleware JWT
type jwtServerStream struct {
	grpc.ServerStream
}

func (j *jwtServerStream) Context() context.Context {
	ctx := j.ServerStream.Context()
	token := ctx.Value("token").(*jwt.Token)
	ctx = context.WithValue(ctx, "token", token)
	return ctx
}
