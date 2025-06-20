package service

import (
	"context"
	"fmt"

	config "github.com/minilikmila/standard-auth-go/configs"
)

type SMSServiceImpl struct {
	config *config.Config
}

func NewSMSService(config *config.Config) SMSService {
	return &SMSServiceImpl{
		config: config,
	}
}

func (s *SMSServiceImpl) SendVerificationSMS(ctx context.Context, phone, code string) error {
	// TODO: Implement actual SMS sending logic
	fmt.Printf("Sending verification SMS to %s with code %s\n", phone, code)
	return nil
}

func (s *SMSServiceImpl) SendPasswordResetSMS(ctx context.Context, phone string, code string) error {
	// TODO: Implement actual SMS sending logic
	fmt.Printf("Sending password reset SMS to %s with code %s\n", phone, code)
	return nil
}
