package oidc4vc

import "context"

func (s *Service) ValidatePreAuthorizeCode(
	ctx context.Context,
	preAuthorizeCode string,
	pin string,
) (*Transaction, error) {
	panic("todo")
}
