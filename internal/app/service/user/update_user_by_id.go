package user

import "context"

func (s *service) UpdateUserByID(ctx context.Context, uid, displayName, email string) error {
	currentUser, err := s.GetSelfInfo(ctx, uid)
	if err != nil {
		return err
	}

	updatedUser := currentUser
	updatedUser.DisplayName = displayName
	updatedUser.Email = email

	return s.repo.UpdateUserByID(ctx, uid, updatedUser)
}
