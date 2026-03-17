package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	ds "team-task-manager/internal/datastruct"
	"team-task-manager/internal/supports"
)

func getStatusCode(s string) int {
	switch s {
	case ds.StatusUserNotFound:
		return http.StatusNotFound
	case ds.StatusResurceNotFound:
		return http.StatusNotFound
	case ds.StatusServiceError:
		return http.StatusInternalServerError
	case ds.StatusUserAlreadyExists:
		return http.StatusConflict
	case ds.StatusResourceAlreadyExists:
		return http.StatusConflict
	case ds.StatusWrongLoginOrPassword:
		return http.StatusUnauthorized
	case ds.StatusInvalidToken:
		return http.StatusUnauthorized
	case ds.StatusSessionReset:
		return http.StatusUnauthorized
	case ds.StatusForbidden:
		return http.StatusForbidden
	case ds.StatusNotOwner:
		return http.StatusForbidden
	case ds.StatusNotMember:
		return http.StatusForbidden
	case ds.StatusConflict:
		return http.StatusConflict
	case ds.StatusIvalidVersion:
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func extractJsonBody[ReqT any](r *http.Request, v ReqT) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil && err != io.EOF {
		return err
	}

	return setJWTUserCredsIfRequire(r, v)
}

func writeJsonResponse[RespT IWithStatus](w *http.ResponseWriter, resp RespT) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&resp); err != nil {
		return err
	}

	(*w).Header().Set(contentLenKey, strconv.Itoa(len(buf.Bytes())))
	(*w).Header().Set(contentTypeKey, appJSONValue)

	code := getStatusCode(resp.GetStatus())
	(*w).WriteHeader(code)
	_, err := (*w).Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func extractSchemaQuery[ReqT any](r *http.Request, v ReqT) error {
	if err := schemaDecoder.Decode(v, r.URL.Query()); err != nil && err != io.EOF {
		return err
	}

	return setJWTUserCredsIfRequire(r, v)
}

func extractJWTCredsOnly[ReqT IWithJWTUserCreds](r *http.Request, v ReqT) error {
	ctxIdVal := r.Context().Value(ds.UserIDKey)
	ctxRoleVal := r.Context().Value(ds.UserRoleKey)
	if ctxIdVal == nil || ctxRoleVal == nil {
		return fmt.Errorf("struct expected to get jwt creds but was not privided")
	}

	id, okId := ctxIdVal.(int64)
	if !okId {
		return fmt.Errorf("key '%s' expected to be 'int64' type but was not", ds.UserIDKey)
	}
	role, okRole := ctxRoleVal.(string)
	if !okRole {
		return fmt.Errorf("key '%s' expected to be 'string' type but was not", ds.UserRoleKey)
	}

	v.SetUserId(id)
	v.SetUserRole(role)

	return nil
}

func setJWTUserCredsIfRequire(r *http.Request, v any) error {
	vv, okV := v.(IWithJWTUserCreds)
	if !okV {
		return nil
	}
	return extractJWTCredsOnly(r, vv)
}

func structValidator[ReqT any](s ReqT) error {
	return supports.StructValidator().Struct(s)
}

func Exec[ReqT any, RespT IWithStatus](a ExecArgs[ReqT, RespT]) {
	var req ReqT

	if err := a.requestExtractor(a.httpRequest, &req); err != nil {
		msg := "failed extracting request"
		a.api.logger.ErrorKV(msg, "error", err.Error())

		resp := ds.Status{Message: supports.Concat(msg, ": ", err.Error())}
		err = writeJsonResponse(a.httpResponse, resp)
		if err != nil {
			a.api.logger.ErrorKV("failed write response",
				"error", err.Error(), "response", resp)
		}

		return
	}

	if a.validator == nil {
		a.validator = structValidator
	}

	if err := a.validator(&req); err != nil {
		msg := "failed validating request"
		a.api.logger.ErrorKV(msg, "error", err.Error(), "request", req)

		resp := ds.Status{Message: supports.Concat(msg, ": ", err.Error())}
		err = writeJsonResponse(a.httpResponse, resp)
		if err != nil {
			a.api.logger.ErrorKV("failed write response",
				"error", err.Error(), "response", resp)
		}
		return
	}

	resp := a.serviceFunc(&req)
	if resp == nil {
		resp := ds.Status{Message: ds.StatusServiceError}
		msg := "failed execute request on service"
		a.api.logger.ErrorKV(msg, "error", "service return no response", "request", req)
		err := writeJsonResponse(a.httpResponse, resp)
		if err != nil {

		}
		return
	}

	if err := a.responseWriter(a.httpResponse, resp); err != nil {
		msg := "failed writing response"
		http.Error(*a.httpResponse, msg, http.StatusInternalServerError)
		a.api.logger.ErrorKV(msg, "error", err.Error(), "request", req)
	}
}
