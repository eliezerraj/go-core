package policy

import (
	"os"
	"fmt"
	"errors"
	"strings"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v3"

	"github.com/cedar-policy/cedar-go"
	"github.com/eliezerraj/go-core/v3/auth"

	"github.com/eliezerraj/go-core/v3/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var (
	ErrAuthorizationDenied = errors.New("authorization denied")
	ErrAuthorizerNotReady  = errors.New("cedar authorizer is not initialized")
)

type CedarClient struct {
	policies	*cedar.PolicySet
	dryRun bool
}

type AuthorizationResult struct {
	Decision   cedar.Decision
	Diagnostic cedar.Diagnostic
}

func NewCedarClient(policyPath string, dryRun bool) (*CedarClient, error) {
	logger.DebugOutCtx("Creating new Cedar client SUCCESSFULLY")

	document, err := os.ReadFile(policyPath)
	if err != nil {
		logger.DebugOutCtx("Failed to read Cedar policy", zap.Error(err))
		return nil, fmt.Errorf("read Cedar policy %q: %w", policyPath, err)
	}

	policies, err := cedar.NewPolicySetFromBytes(policyPath,document)
	if err != nil {
		logger.DebugOutCtx("Failed to parse Cedar policy", zap.Error(err))
		return nil, fmt.Errorf("parse Cedar policy %q: %w", policyPath, err)
	}

	return &CedarClient{
		policies: policies,
		dryRun:   dryRun,
	}, nil
}

func (p *CedarClient) EvaluatePolicy(req cedar.Request, entities cedar.EntityGetter) (AuthorizationResult, error) {
	logger.DebugOutCtx("Evaluating policy")

	if p.policies == nil {
		logger.DebugOutCtx("Cedar authorizer is not initialized")
		return AuthorizationResult{}, ErrAuthorizerNotReady
	}

	decision, diagnostic := cedar.Authorize(p.policies, entities, req)

	return AuthorizationResult{
		Decision:   decision,
		Diagnostic: diagnostic,
	}, nil
}

func (p *CedarClient) Authorize(entities cedar.EntityGetter, request cedar.Request) error {
	logger.DebugOutCtx("Authorizing request")

	result, err := p.EvaluatePolicy(request, entities)
	if err != nil {
		logger.DebugOutCtx("Failed to evaluate policy", zap.Error(err))
		return err
	}

	if result.Decision != cedar.Allow {
		logger.DebugOutCtx("Authorization denied")
		return ErrAuthorizationDenied
	}

	return nil
}

func (p *CedarClient) FiberPolicyEvaluationMiddleware(action string) fiber.Handler {
	return func(c fiber.Ctx) error {

		if p.dryRun {
			return c.Next()
		}

		ctx, span := otel.Tracer("go-core").Start(c.Context(), "auth.policy_evaluation")
		c.SetContext(ctx)

		if p.policies == nil {
			logger.Debug(ctx, "Cedar authorizer is not initialized")
			return c.Status(fiber.StatusInternalServerError).SendString("Cedar authorizer is not initialized")
		}

		accessTokenClaims, ok := c.Context().Value(auth.AccessTokenClaimsKey).(*auth.AccessTokenClaims)
		if !ok {
			logger.Error(ctx, "Access token claims not found in context")
			span.RecordError(errors.New("Unauthorized by token claims malformed"))
			span.SetStatus(codes.Error, "Unauthorized by token claims malformed")
			span.End()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status_code": fiber.StatusUnauthorized,
				"reason": "Unauthorized by token claims malformed",
				"error":"Unauthorized by token claims malformed",
			})
		} else {
			logger.Info(ctx, "Access token claims found", zap.Any("claims", accessTokenClaims))
		}

		principalUID := cedar.NewEntityUID(
			cedar.EntityType("Client"),
			cedar.String(accessTokenClaims.ClientID),
		)

		// Parse scopes from accessTokenClaims.Scope string
		scopeList := strings.Fields(accessTokenClaims.Scope)
		scopeSet := make([]cedar.Value, len(scopeList))
		for i, scope := range scopeList {
			scopeSet[i] = cedar.String(scope)
		}

		roleList := strings.Fields(accessTokenClaims.Role)
		roleSet := make([]cedar.Value, len(roleList))
		for i, role := range roleList {
			roleSet[i] = cedar.String(role)
		}

		principal := cedar.Entity{
			UID: principalUID,
			Attributes: cedar.NewRecord(
				cedar.RecordMap{
					"client_id": cedar.String(accessTokenClaims.ClientID),
					"iss":       cedar.String(accessTokenClaims.Issuer),
					"sub":       cedar.String(accessTokenClaims.Subject),
					"scopes":    cedar.NewSet(scopeSet...),
					"roles":     cedar.NewSet(roleSet...),
				},
			),
			Parents: cedar.NewEntityUIDSet(),
			Tags:    cedar.NewRecord(cedar.RecordMap{}),
		}

		// NOT IMPLEMENTED: Resource identification is not yet implemented
		resourceUID := cedar.NewEntityUID(
			cedar.EntityType("Resource"),
			cedar.String(""),
		)

		resource := cedar.Entity{
			UID: resourceUID,
			Attributes: cedar.NewRecord(
				cedar.RecordMap{},
			),
			Parents: cedar.NewEntityUIDSet(),
			Tags:    cedar.NewRecord(cedar.RecordMap{}),
		}

		entities := cedar.EntityMap{
			principalUID: principal,
			resourceUID:  resource,
		}

		request := cedar.Request{
			Principal: principalUID,
			Action: cedar.NewEntityUID(
				cedar.EntityType("Action"),
				cedar.String(action),
			),
			Resource: resourceUID,
			Context: cedar.NewRecord(
				cedar.RecordMap{},
			),
		}

		result, err := p.EvaluatePolicy(
			request,
			entities,
		)

		if err != nil {
			logger.Debug(ctx, "Failed to evaluate policy", zap.Error(err))
			return err
		}

		if result.Decision != cedar.Allow {
			logger.Warn(ctx, "Authorization denied")
			span.RecordError(errors.New("Unauthorized by policy"))
			span.SetStatus(codes.Error, "Unauthorized by policy")
			span.End()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status_code": fiber.StatusUnauthorized,
				"reason": "Unauthorized by CEDAR policy: " + action,
				"error":"Unauthorized by CEDAR policy",
			})
		} 

		logger.Info(ctx, "Authorization granted")
		
		span.End()
		
        return c.Next()
	}
}
