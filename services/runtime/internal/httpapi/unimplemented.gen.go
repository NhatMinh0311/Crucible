// Code generated from api.StrictServerInterface; DO NOT EDIT. Tạo lại bằng go generate ./internal/httpapi

package httpapi

import (
	"context"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/api"
	"github.com/NhatMinh0311/crucible/services/runtime/internal/problem"
)

// Unimplemented trả 501 NOT_IMPLEMENTED cho mọi operation. Server nhúng kiểu này và chỉ ghi đè operation có ở chế độ đang chạy.
type Unimplemented struct{}

func (Unimplemented) GetCapacity(context.Context, api.GetCapacityRequestObject) (api.GetCapacityResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) FindCapacityWindows(context.Context, api.FindCapacityWindowsRequestObject) (api.FindCapacityWindowsResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) ListQuotaAccounts(context.Context, api.ListQuotaAccountsRequestObject) (api.ListQuotaAccountsResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) GetQuotaAccount(context.Context, api.GetQuotaAccountRequestObject) (api.GetQuotaAccountResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) PutQuotaAccount(context.Context, api.PutQuotaAccountRequestObject) (api.PutQuotaAccountResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) CreateReservation(context.Context, api.CreateReservationRequestObject) (api.CreateReservationResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) CancelReservation(context.Context, api.CancelReservationRequestObject) (api.CancelReservationResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) GetReservation(context.Context, api.GetReservationRequestObject) (api.GetReservationResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) ResizeReservation(context.Context, api.ResizeReservationRequestObject) (api.ResizeReservationResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) ListResourceViolations(context.Context, api.ListResourceViolationsRequestObject) (api.ListResourceViolationsResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) CreateSession(context.Context, api.CreateSessionRequestObject) (api.CreateSessionResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) GetSession(context.Context, api.GetSessionRequestObject) (api.GetSessionResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) ExecInSession(context.Context, api.ExecInSessionRequestObject) (api.ExecInSessionResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) ExtendSession(context.Context, api.ExtendSessionRequestObject) (api.ExtendSessionResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) GetRecording(context.Context, api.GetRecordingRequestObject) (api.GetRecordingResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) CreateTerminalTicket(context.Context, api.CreateTerminalTicketRequestObject) (api.CreateTerminalTicketResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) TerminateSession(context.Context, api.TerminateSessionRequestObject) (api.TerminateSessionResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) CreateViewTicket(context.Context, api.CreateViewTicketRequestObject) (api.CreateViewTicketResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) ListSettings(context.Context, api.ListSettingsRequestObject) (api.ListSettingsResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) UpdateSetting(context.Context, api.UpdateSettingRequestObject) (api.UpdateSettingResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

func (Unimplemented) ListSettingChanges(context.Context, api.ListSettingChangesRequestObject) (api.ListSettingChangesResponseObject, error) {
	return nil, problem.ErrNotImplemented
}

var _ api.StrictServerInterface = Unimplemented{}
