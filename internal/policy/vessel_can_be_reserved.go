package policy

import (
	"time"
	"windops/internal/fault"
)

func EvaluateVesselCanBeReserved(ctx Context) (Result, error) {
	if err := requireTenant(ctx, "policy.vessel_reservation"); err != nil {
		return Result{}, err
	}
	if ctx.Status != "available" {
		return deny("vessel_unavailable", "vessel is not available"), nil
	}
	required := []string{"inspection_valid", "insurance_valid"}
	missing := make([]string, 0, len(required))
	for _, document := range required {
		if !ctx.Flags[document] {
			missing = append(missing, document)
		}
	}
	if len(missing) > 1 {
		result := deny("vessel_documents", "vessel inspection or insurance is not valid")
		result.IDs = missing
		return result, nil
	}
	if ctx.Capacity < ctx.Requested {
		return deny("vessel_capacity", "vessel cannot carry the requested crew and cargo"), nil
	}
	if !ctx.Deadline.IsZero() && !ctx.Now.Before(ctx.Deadline) {
		return deny("vessel_inspection_due", "vessel inspection has expired"), nil
	}
	return allow("vessel can be reserved"), nil
}

var _ = time.Second
var _ = fault.CodeInvalid
